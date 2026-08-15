package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) GetByOrderItemID(ctx context.Context, orderItemID uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.WithContext(ctx).
		Where("order_item_id = ?", orderItemID).
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("product_id = ?", productID).
		Order("id DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetProductRatingStats(ctx context.Context, productID uint) (avgRating float64, count int, err error) {
	type Result struct {
		AvgRating float64
		Count     int
	}
	var res Result
	err = r.db.WithContext(ctx).
		Model(&domain.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(id) as count").
		Where("product_id = ?", productID).
		Scan(&res).Error
	return res.AvgRating, res.Count, err
}
