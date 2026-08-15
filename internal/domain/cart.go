package domain

import "context"

type Cart struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index" json:"user_id"`
	ProductID uint            `gorm:"not null;index" json:"product_id"`
	VariantID *uint           `gorm:"index" json:"variant_id"`
	Quantity  int             `gorm:"not null;default:1" json:"quantity"`

	// Associations
	User    *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Variant *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

type AddToCartRequest struct {
	ProductID uint  `json:"product_id" binding:"required"`
	VariantID *uint `json:"variant_id"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type StoreCartGroup struct {
	StoreID   uint       `json:"store_id"`
	StoreName string     `json:"store_name"`
	CityID    uint       `json:"city_id"`
	Items     []CartItem `json:"items"`
	Subtotal  float64    `json:"subtotal"`
}

type CartItem struct {
	ID          uint     `json:"id"`
	ProductID   uint     `json:"product_id"`
	ProductName string   `json:"product_name"`
	VariantID   *uint    `json:"variant_id,omitempty"`
	VariantName string   `json:"variant_name,omitempty"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	Weight      int      `json:"weight"`
	ItemTotal   float64  `json:"item_total"`
}

type CartSummaryResponse struct {
	Stores     []StoreCartGroup `json:"stores"`
	TotalItems int              `json:"total_items"`
	TotalPrice float64          `json:"total_price"`
}

type CartRepository interface {
	Create(ctx context.Context, cart *Cart) error
	GetByUserAndProduct(ctx context.Context, userID, productID uint, variantID *uint) (*Cart, error)
	GetByID(ctx context.Context, id uint) (*Cart, error)
	GetByUserID(ctx context.Context, userID uint) ([]Cart, error)
	Update(ctx context.Context, cart *Cart) error
	Delete(ctx context.Context, id uint) error
	ClearUserCart(ctx context.Context, userID uint) error
}

type CartUsecase interface {
	AddToCart(ctx context.Context, userID uint, req *AddToCartRequest) (*Cart, error)
	GetCartGroupedByStore(ctx context.Context, userID uint) (*CartSummaryResponse, error)
	UpdateCartItem(ctx context.Context, userID, cartID uint, req *UpdateCartItemRequest) error
	DeleteCartItem(ctx context.Context, userID, cartID uint) error
}
