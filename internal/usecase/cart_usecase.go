package usecase

import (
	"context"
	"errors"

	"go-market/internal/domain"
)

var (
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrInvalidVariant   = errors.New("invalid product variant")
	ErrExceedStock      = errors.New("requested quantity exceeds available stock")
)

type cartUsecase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
}

func NewCartUsecase(cartRepo domain.CartRepository, productRepo domain.ProductRepository) domain.CartUsecase {
	return &cartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *cartUsecase) AddToCart(ctx context.Context, userID uint, req *domain.AddToCartRequest) (*domain.Cart, error) {
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Validate stock for product or variant
	availableStock := product.Stock
	if req.VariantID != nil && *req.VariantID > 0 {
		var found bool
		for _, v := range product.Variants {
			if v.ID == *req.VariantID {
				availableStock = v.Stock
				found = true
				break
			}
		}
		if !found {
			return nil, ErrInvalidVariant
		}
	}

	existingItem, _ := u.cartRepo.GetByUserAndProduct(ctx, userID, req.ProductID, req.VariantID)
	if existingItem != nil {
		newQty := existingItem.Quantity + req.Quantity
		if newQty > availableStock {
			return nil, ErrExceedStock
		}
		existingItem.Quantity = newQty
		if err := u.cartRepo.Update(ctx, existingItem); err != nil {
			return nil, err
		}
		return existingItem, nil
	}

	if req.Quantity > availableStock {
		return nil, ErrExceedStock
	}

	cart := &domain.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
	}

	if err := u.cartRepo.Create(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecase) GetCartGroupedByStore(ctx context.Context, userID uint) (*domain.CartSummaryResponse, error) {
	carts, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	storeMap := make(map[uint]*domain.StoreCartGroup)
	var storeOrder []uint
	var grandTotal float64
	var totalItems int

	for _, c := range carts {
		if c.Product == nil || c.Product.Store == nil {
			continue
		}

		storeID := c.Product.StoreID
		group, exists := storeMap[storeID]
		if !exists {
			group = &domain.StoreCartGroup{
				StoreID:   storeID,
				StoreName: c.Product.Store.StoreName,
				CityID:    c.Product.Store.CityID,
				Items:     []domain.CartItem{},
				Subtotal:  0,
			}
			storeMap[storeID] = group
			storeOrder = append(storeOrder, storeID)
		}

		price := c.Product.Price
		variantName := ""
		if c.Variant != nil {
			variantName = c.Variant.VariantName
			if c.Variant.PriceOverride != nil {
				price = *c.Variant.PriceOverride
			}
		}

		itemTotal := price * float64(c.Quantity)
		group.Subtotal += itemTotal
		grandTotal += itemTotal
		totalItems += c.Quantity

		group.Items = append(group.Items, domain.CartItem{
			ID:          c.ID,
			ProductID:   c.ProductID,
			ProductName: c.Product.Name,
			VariantID:   c.VariantID,
			VariantName: variantName,
			Price:       price,
			Quantity:    c.Quantity,
			Weight:      c.Product.Weight,
			ItemTotal:   itemTotal,
		})
	}

	var storeGroups []domain.StoreCartGroup
	for _, storeID := range storeOrder {
		storeGroups = append(storeGroups, *storeMap[storeID])
	}

	if storeGroups == nil {
		storeGroups = []domain.StoreCartGroup{}
	}

	return &domain.CartSummaryResponse{
		Stores:     storeGroups,
		TotalItems: totalItems,
		TotalPrice: grandTotal,
	}, nil
}

func (u *cartUsecase) UpdateCartItem(ctx context.Context, userID, cartID uint, req *domain.UpdateCartItemRequest) error {
	cart, err := u.cartRepo.GetByID(ctx, cartID)
	if err != nil || cart.UserID != userID {
		return ErrCartItemNotFound
	}

	availableStock := cart.Product.Stock
	if cart.Variant != nil {
		availableStock = cart.Variant.Stock
	}

	if req.Quantity > availableStock {
		return ErrExceedStock
	}

	cart.Quantity = req.Quantity
	return u.cartRepo.Update(ctx, cart)
}

func (u *cartUsecase) DeleteCartItem(ctx context.Context, userID, cartID uint) error {
	cart, err := u.cartRepo.GetByID(ctx, cartID)
	if err != nil || cart.UserID != userID {
		return ErrCartItemNotFound
	}

	return u.cartRepo.Delete(ctx, cartID)
}
