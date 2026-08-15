package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	StoreID     uint             `gorm:"not null;index" json:"store_id"`
	CategoryID  uint             `gorm:"not null;index" json:"category_id"`
	Name        string           `gorm:"size:255;not null;index" json:"name"`
	Description string           `gorm:"type:text" json:"description"`
	ImageURL    string           `gorm:"type:text" json:"image_url"`
	Price       float64          `gorm:"type:decimal(15,2);not null" json:"price"`
	Stock       int              `gorm:"not null;default:0" json:"stock"`
	Weight      int              `gorm:"not null;default:0" json:"weight"` // in grams
	RatingAvg   float64          `gorm:"type:decimal(3,2);default:0.00" json:"rating_avg"`
	RatingCount int              `gorm:"default:0" json:"rating_count"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Associations
	Store    *Store           `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Category *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Reviews  []Review         `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

type ProductVariant struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	ProductID     uint     `gorm:"not null;index" json:"product_id"`
	VariantName   string   `gorm:"size:100;not null" json:"variant_name"`
	PriceOverride *float64 `gorm:"type:decimal(15,2)" json:"price_override"`
	Stock         int      `gorm:"not null;default:0" json:"stock"`
}

type CreateProductVariantRequest struct {
	VariantName   string   `json:"variant_name" binding:"required"`
	PriceOverride *float64 `json:"price_override"`
	Stock         int      `json:"stock" binding:"required,min=0"`
}

type CreateProductRequest struct {
	CategoryID  uint                          `json:"category_id" binding:"required"`
	Name        string                        `json:"name" binding:"required,min=3,max=200"`
	Description string                        `json:"description"`
	ImageURL    string                        `json:"image_url"`
	Price       float64                       `json:"price" binding:"required,gt=0"`
	Stock       int                           `json:"stock" binding:"min=0"`
	Weight      int                           `json:"weight" binding:"required,gt=0"`
	Variants    []CreateProductVariantRequest `json:"variants"`
}

type UpdateProductRequest struct {
	CategoryID  *uint    `json:"category_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	ImageURL    *string  `json:"image_url"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	Weight      *int     `json:"weight"`
}

type ProductFilter struct {
	Search     string  `form:"search"`
	CategoryID *uint   `form:"category_id"`
	StoreID    *uint   `form:"store_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Page       int     `form:"page,default=1"`
	Limit      int     `form:"limit,default=10"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uint) (*Product, error)
	List(ctx context.Context, filter *ProductFilter) ([]Product, int64, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	DeductStock(ctx context.Context, productID uint, variantID *uint, quantity int) error
	AddStock(ctx context.Context, productID uint, variantID *uint, quantity int) error
	UpdateRating(ctx context.Context, productID uint, avgRating float64, count int) error
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, userID uint, req *CreateProductRequest) (*Product, error)
	GetProductByID(ctx context.Context, id uint) (*Product, error)
	ListProducts(ctx context.Context, filter *ProductFilter) ([]Product, int64, error)
	UpdateProduct(ctx context.Context, userID uint, productID uint, req *UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, userID uint, productID uint) error
}
