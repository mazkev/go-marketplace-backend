package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-market/internal/domain"
)

var (
	ErrCartIsEmpty             = errors.New("cart is empty")
	ErrAddressNotFound         = errors.New("delivery address not found")
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderItemNotFound       = errors.New("order item not found")
	ErrNotAuthorizedForStore   = errors.New("you are not authorized to manage this store's orders")
	ErrCannotShipStatus        = errors.New("order item cannot be shipped from current status")
	ErrCannotCompleteStatus    = errors.New("only shipped or delivered order items can be marked as completed")
	ErrNotOrderOwner           = errors.New("you do not own this order")
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	cartRepo    domain.CartRepository
	userRepo    domain.UserRepository
	storeRepo   domain.StoreRepository
	productRepo domain.ProductRepository
	voucherRepo domain.VoucherRepository
	walletRepo  domain.WalletRepository
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	userRepo domain.UserRepository,
	storeRepo domain.StoreRepository,
	productRepo domain.ProductRepository,
	voucherRepo domain.VoucherRepository,
	walletRepo domain.WalletRepository,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		userRepo:    userRepo,
		storeRepo:   storeRepo,
		productRepo: productRepo,
		voucherRepo: voucherRepo,
		walletRepo:  walletRepo,
	}
}

func (u *orderUsecase) Checkout(ctx context.Context, userID uint, req *domain.CheckoutRequest) (*domain.Order, error) {
	// 1. Validate cart
	carts, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil || len(carts) == 0 {
		return nil, ErrCartIsEmpty
	}

	// 2. Validate Address
	addresses, err := u.userRepo.GetAddressesByUserID(ctx, userID)
	if err != nil || len(addresses) == 0 {
		return nil, ErrAddressNotFound
	}
	var addressFound bool
	for _, addr := range addresses {
		if addr.ID == req.AddressID {
			addressFound = true
			break
		}
	}
	if !addressFound {
		return nil, ErrAddressNotFound
	}

	// 3. Map courier options per store
	courierMap := make(map[uint]string)
	for _, opt := range req.Stores {
		courierMap[opt.StoreID] = opt.CourierName
	}

	// 4. Build OrderItems and calculate totals
	var orderItems []domain.OrderItem
	var totalItemsAmount float64
	var totalShippingCost float64

	storeShippingApplied := make(map[uint]bool)

	for _, c := range carts {
		if c.Product == nil {
			return nil, ErrProductNotFound
		}

		price := c.Product.Price
		if c.Variant != nil && c.Variant.PriceOverride != nil {
			price = *c.Variant.PriceOverride
		}

		storeID := c.Product.StoreID
		courier := courierMap[storeID]
		if courier == "" {
			courier = "JNE Regular"
		}

		var shippingCost float64
		if !storeShippingApplied[storeID] {
			shippingCost = 10000.0 // Default mock shipping cost per store
			storeShippingApplied[storeID] = true
			totalShippingCost += shippingCost
		}

		itemTotal := price * float64(c.Quantity)
		totalItemsAmount += itemTotal

		orderItems = append(orderItems, domain.OrderItem{
			StoreID:      storeID,
			ProductID:    c.ProductID,
			VariantID:    c.VariantID,
			Quantity:     c.Quantity,
			Price:        price,
			ShippingCost: shippingCost,
			CourierName:  courier,
			Status:       domain.OrderItemStatusProcessing,
		})
	}

	totalOrderAmount := totalItemsAmount + totalShippingCost
	discountAmount := 0.0

	// 5. Apply Voucher if provided
	if req.VoucherCode != "" {
		voucher, err := u.voucherRepo.GetByCode(ctx, req.VoucherCode)
		if err == nil && voucher != nil {
			now := time.Now()
			if now.After(voucher.StartDate) && now.Before(voucher.EndDate) && voucher.UsedCount < voucher.Quota && totalOrderAmount >= voucher.MinSpend {
				if voucher.VoucherType == domain.VoucherTypePercentage && voucher.DiscountPercent != nil {
					discountAmount = (totalOrderAmount * (*voucher.DiscountPercent)) / 100.0
					if voucher.MaxDiscount != nil && discountAmount > *voucher.MaxDiscount {
						discountAmount = *voucher.MaxDiscount
					}
				} else if voucher.VoucherType == domain.VoucherTypeFixed && voucher.DiscountAmount != nil {
					discountAmount = *voucher.DiscountAmount
				}
				_ = u.voucherRepo.IncrementUsage(ctx, voucher.ID)
			}
		}
	}

	finalAmount := totalOrderAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	// 6. Generate unique invoice number & VA
	invoiceNo := fmt.Sprintf("INV/%s/%04d%04d",
		time.Now().Format("20060102"),
		rand.Intn(10000),
		userID,
	)

	payMethod := req.PaymentMethod
	if payMethod == "" {
		payMethod = "BCA_VA"
	}
	vaNumber := fmt.Sprintf("8808%04d%04d", rand.Intn(10000), userID)
	expiredAt := time.Now().Add(24 * time.Hour)

	order := &domain.Order{
		UserID:           userID,
		InvoiceNumber:    invoiceNo,
		TotalAmount:      totalOrderAmount,
		DiscountAmount:   discountAmount,
		FinalAmount:      finalAmount,
		VoucherCode:      req.VoucherCode,
		PaymentMethod:    payMethod,
		VANumber:         vaNumber,
		PaymentStatus:    domain.PaymentStatusPaid, // In MVP simulation, payment is paid immediately (Escrow held by platform)
		PaymentExpiredAt: &expiredAt,
	}

	// 7. Execute atomic transaction
	if err := u.orderRepo.CreateOrderWithItems(ctx, order, orderItems); err != nil {
		return nil, err
	}

	return u.orderRepo.GetByID(ctx, order.ID)
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, userID, orderID uint) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID != userID {
		return nil, ErrNotOrderOwner
	}

	return order, nil
}

func (u *orderUsecase) GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error) {
	return u.orderRepo.GetByUserID(ctx, userID)
}

func (u *orderUsecase) GetStoreOrders(ctx context.Context, userID uint) ([]domain.OrderItem, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrNotAuthorizedForStore
	}

	return u.orderRepo.GetOrderItemsByStoreID(ctx, store.ID)
}

func (u *orderUsecase) ShipOrderItem(ctx context.Context, sellerUserID, orderItemID uint, req *domain.ShipOrderRequest) (*domain.OrderItem, error) {
	item, err := u.orderRepo.GetOrderItemByID(ctx, orderItemID)
	if err != nil {
		return nil, ErrOrderItemNotFound
	}

	store, err := u.storeRepo.GetByUserID(ctx, sellerUserID)
	if err != nil || store.ID != item.StoreID {
		return nil, ErrNotAuthorizedForStore
	}

	if item.Status != domain.OrderItemStatusProcessing && item.Status != domain.OrderItemStatusPending {
		return nil, ErrCannotShipStatus
	}

	item.TrackingNumber = req.TrackingNumber
	if req.CourierName != "" {
		item.CourierName = req.CourierName
	}
	item.Status = domain.OrderItemStatusShipped

	if err := u.orderRepo.UpdateOrderItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (u *orderUsecase) CompleteOrderItem(ctx context.Context, buyerUserID, orderItemID uint) (*domain.OrderItem, error) {
	item, err := u.orderRepo.GetOrderItemByID(ctx, orderItemID)
	if err != nil {
		return nil, ErrOrderItemNotFound
	}

	order, err := u.orderRepo.GetByID(ctx, item.OrderID)
	if err != nil || order.UserID != buyerUserID {
		return nil, ErrNotOrderOwner
	}

	if item.Status != domain.OrderItemStatusShipped && item.Status != domain.OrderItemStatusDelivered {
		return nil, ErrCannotCompleteStatus
	}

	item.Status = domain.OrderItemStatusCompleted

	if err := u.orderRepo.UpdateOrderItem(ctx, item); err != nil {
		return nil, err
	}

	// Escrow Release: Add item total revenue (price * qty + shipping_cost) to seller's store balance
	payoutAmount := (item.Price * float64(item.Quantity)) + item.ShippingCost
	_ = u.storeRepo.AddBalance(ctx, item.StoreID, payoutAmount)

	// Record Balance Mutation Ledger
	store, _ := u.storeRepo.GetByID(ctx, item.StoreID)
	newBalance := payoutAmount
	if store != nil {
		newBalance = store.Balance
	}
	_ = u.walletRepo.RecordMutation(ctx, &domain.BalanceMutation{
		StoreID:      item.StoreID,
		Amount:       payoutAmount,
		Type:         domain.MutationTypeCredit,
		Description:  fmt.Sprintf("Escrow released for completed order item #%d (Order %s)", item.ID, order.InvoiceNumber),
		BalanceAfter: newBalance,
	})

	return item, nil
}
