package domain

import (
	"context"
	"time"
)

type Review struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderItemID uint      `gorm:"uniqueIndex;not null" json:"order_item_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	Rating      int       `gorm:"not null" json:"rating"` // 1-5
	Comment     string    `gorm:"type:text" json:"comment"`
	CreatedAt   time.Time `json:"created_at"`

	// Associations
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type CreateReviewRequest struct {
	OrderItemID uint   `json:"order_item_id" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Comment     string `json:"comment"`
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	GetByOrderItemID(ctx context.Context, orderItemID uint) (*Review, error)
	GetByProductID(ctx context.Context, productID uint) ([]Review, error)
	GetProductRatingStats(ctx context.Context, productID uint) (avgRating float64, count int, err error)
}

type ReviewUsecase interface {
	CreateReview(ctx context.Context, userID uint, req *CreateReviewRequest) (*Review, error)
	GetProductReviews(ctx context.Context, productID uint) ([]Review, error)
}
