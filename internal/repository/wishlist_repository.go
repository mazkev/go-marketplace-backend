package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) domain.WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Add(ctx context.Context, userID, productID uint) error {
	item := domain.Wishlist{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.WithContext(ctx).Where(item).FirstOrCreate(&item).Error
}

func (r *wishlistRepository) Remove(ctx context.Context, userID, productID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&domain.Wishlist{}).Error
}

func (r *wishlistRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Wishlist, error) {
	var list []domain.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product.Store").
		Preload("Product.Category").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&list).Error
	return list, err
}

func (r *wishlistRepository) IsWishlisted(ctx context.Context, userID, productID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Wishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}
