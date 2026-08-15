package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) GetByUserAndProduct(ctx context.Context, userID, productID uint, variantID *uint) (*domain.Cart, error) {
	var cart domain.Cart
	query := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID)
	if variantID != nil && *variantID > 0 {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByID(ctx context.Context, id uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.WithContext(ctx).
		Preload("Product.Store").
		Preload("Variant").
		First(&cart, id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.WithContext(ctx).
		Preload("Product.Store").
		Preload("Variant").
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&carts).Error
	return carts, err
}

func (r *cartRepository) Update(ctx context.Context, cart *domain.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Cart{}, id).Error
}

func (r *cartRepository) ClearUserCart(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Cart{}).Error
}
