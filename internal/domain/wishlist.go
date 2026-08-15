package domain

import (
	"context"
	"time"
)

type Wishlist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_product" json:"user_id"`
	ProductID uint      `gorm:"not null;uniqueIndex:idx_user_product" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`

	// Associations
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type WishlistRepository interface {
	Add(ctx context.Context, userID, productID uint) error
	Remove(ctx context.Context, userID, productID uint) error
	GetByUserID(ctx context.Context, userID uint) ([]Wishlist, error)
	IsWishlisted(ctx context.Context, userID, productID uint) (bool, error)
}

type WishlistUsecase interface {
	AddToWishlist(ctx context.Context, userID, productID uint) error
	RemoveFromWishlist(ctx context.Context, userID, productID uint) error
	GetUserWishlist(ctx context.Context, userID uint) ([]Product, error)
}
