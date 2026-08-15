package domain

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusFailed  PaymentStatus = "FAILED"
	PaymentStatusExpired PaymentStatus = "EXPIRED"
)

type OrderItemStatus string

const (
	OrderItemStatusPending    OrderItemStatus = "PENDING"
	OrderItemStatusProcessing OrderItemStatus = "PROCESSING"
	OrderItemStatusShipped    OrderItemStatus = "SHIPPED"
	OrderItemStatusDelivered  OrderItemStatus = "DELIVERED"
	OrderItemStatusCompleted  OrderItemStatus = "COMPLETED"
	OrderItemStatusCancelled  OrderItemStatus = "CANCELLED"
)

type Order struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	UserID           uint          `gorm:"not null;index" json:"user_id"`
	InvoiceNumber    string        `gorm:"size:50;uniqueIndex;not null" json:"invoice_number"`
	TotalAmount      float64       `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	DiscountAmount   float64       `gorm:"type:decimal(15,2);default:0.00;not null" json:"discount_amount"`
	FinalAmount      float64       `gorm:"type:decimal(15,2);not null" json:"final_amount"`
	VoucherCode      string        `gorm:"size:50" json:"voucher_code,omitempty"`
	PaymentMethod    string        `gorm:"size:50;default:'BCA_VA'" json:"payment_method"`
	VANumber         string        `gorm:"size:50" json:"va_number,omitempty"`
	PaymentStatus    PaymentStatus `gorm:"size:20;not null;default:'PENDING'" json:"payment_status"`
	PaymentExpiredAt *time.Time    `json:"payment_expired_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`

	// Associations
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	OrderID        uint            `gorm:"not null;index" json:"order_id"`
	StoreID        uint            `gorm:"not null;index" json:"store_id"`
	ProductID      uint            `gorm:"not null;index" json:"product_id"`
	VariantID      *uint           `gorm:"index" json:"variant_id"`
	Quantity       int             `gorm:"not null" json:"quantity"`
	Price          float64         `gorm:"type:decimal(15,2);not null" json:"price"`
	ShippingCost   float64         `gorm:"type:decimal(15,2);not null;default:0.00" json:"shipping_cost"`
	CourierName    string          `gorm:"size:50" json:"courier_name"`
	TrackingNumber string          `gorm:"size:100" json:"tracking_number"`
	Status         OrderItemStatus `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	// Associations
	Store   *Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Product *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Variant *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Review  *Review         `gorm:"foreignKey:OrderItemID" json:"review,omitempty"`
}

// Request & Response DTOs
type StoreCheckoutOption struct {
	StoreID     uint   `json:"store_id" binding:"required"`
	CourierName string `json:"courier_name" binding:"required"`
}

type CheckoutRequest struct {
	AddressID     uint                  `json:"address_id" binding:"required"`
	PaymentMethod string                `json:"payment_method"` // default BCA_VA
	VoucherCode   string                `json:"voucher_code"`
	Stores        []StoreCheckoutOption `json:"stores"`
}

type ShipOrderRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required,min=4"`
	CourierName    string `json:"courier_name"`
}

type PaymentWebhookPayload struct {
	InvoiceNumber string  `json:"invoice_number" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentStatus string  `json:"payment_status" binding:"required"` // SETTLEMENT, EXPIRE, CANCEL
	SignatureKey  string  `json:"signature_key"`
}

type OrderRepository interface {
	CreateOrderWithItems(ctx context.Context, order *Order, items []OrderItem) error
	GetByID(ctx context.Context, orderID uint) (*Order, error)
	GetByInvoiceNumber(ctx context.Context, invoiceNo string) (*Order, error)
	GetByUserID(ctx context.Context, userID uint) ([]Order, error)
	GetOrderItemByID(ctx context.Context, orderItemID uint) (*OrderItem, error)
	GetOrderItemsByStoreID(ctx context.Context, storeID uint) ([]OrderItem, error)
	UpdateOrderItem(ctx context.Context, item *OrderItem) error
	UpdatePaymentStatus(ctx context.Context, orderID uint, status PaymentStatus) error
	UpdateOrderAndItemsStatusOnPayment(ctx context.Context, orderID uint, paymentStatus PaymentStatus, itemStatus OrderItemStatus) error
}

type OrderUsecase interface {
	Checkout(ctx context.Context, userID uint, req *CheckoutRequest) (*Order, error)
	GetOrderByID(ctx context.Context, userID, orderID uint) (*Order, error)
	GetUserOrders(ctx context.Context, userID uint) ([]Order, error)
	GetStoreOrders(ctx context.Context, userID uint) ([]OrderItem, error)
	ShipOrderItem(ctx context.Context, sellerUserID, orderItemID uint, req *ShipOrderRequest) (*OrderItem, error)
	CompleteOrderItem(ctx context.Context, buyerUserID, orderItemID uint) (*OrderItem, error)
}

type PaymentUsecase interface {
	HandleWebhook(ctx context.Context, payload *PaymentWebhookPayload) error
}
