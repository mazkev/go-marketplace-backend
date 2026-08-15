package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithItems(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create main order record
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Process each item: deduct stock and attach order_id
		for i := range items {
			items[i].OrderID = order.ID

			// Deduct stock safely with transaction
			if items[i].VariantID != nil && *items[i].VariantID > 0 {
				res := tx.Model(&domain.ProductVariant{}).
					Where("id = ? AND product_id = ? AND stock >= ?", *items[i].VariantID, items[i].ProductID, items[i].Quantity).
					Update("stock", gorm.Expr("stock - ?", items[i].Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return ErrInsufficientStock
				}
			} else {
				res := tx.Model(&domain.Product{}).
					Where("id = ? AND stock >= ?", items[i].ProductID, items[i].Quantity).
					Update("stock", gorm.Expr("stock - ?", items[i].Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return ErrInsufficientStock
				}
			}

			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		// 3. Clear user cart
		if err := tx.Where("user_id = ?", order.UserID).Delete(&domain.Cart{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) GetByID(ctx context.Context, orderID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("OrderItems.Store").
		Preload("OrderItems.Product").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Review").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByInvoiceNumber(ctx context.Context, invoiceNo string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("OrderItems.Store").
		Preload("OrderItems.Product").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Review").
		Where("invoice_number = ?", invoiceNo).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Store").
		Preload("OrderItems.Product").
		Preload("OrderItems.Variant").
		Preload("OrderItems.Review").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetOrderItemByID(ctx context.Context, orderItemID uint) (*domain.OrderItem, error) {
	var item domain.OrderItem
	err := r.db.WithContext(ctx).
		Preload("Store").
		Preload("Product").
		Preload("Variant").
		Preload("Review").
		First(&item, orderItemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *orderRepository) GetOrderItemsByStoreID(ctx context.Context, storeID uint) ([]domain.OrderItem, error) {
	var items []domain.OrderItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Variant").
		Preload("Review").
		Where("store_id = ?", storeID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *orderRepository) UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, orderID uint, status domain.PaymentStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", status).Error
}

func (r *orderRepository) UpdateOrderAndItemsStatusOnPayment(ctx context.Context, orderID uint, paymentStatus domain.PaymentStatus, itemStatus domain.OrderItemStatus) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Order{}).Where("id = ?", orderID).Update("payment_status", paymentStatus).Error; err != nil {
			return err
		}
		return tx.Model(&domain.OrderItem{}).Where("order_id = ?", orderID).Update("status", itemStatus).Error
	})
}
