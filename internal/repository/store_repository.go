package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) domain.StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *domain.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) GetByID(ctx context.Context, id uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.WithContext(ctx).
		Preload("Products").
		First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetBySlug(ctx context.Context, slug string) (*domain.Store, error) {
	var store domain.Store
	err := r.db.WithContext(ctx).
		Where("domain_slug = ?", slug).
		First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) AddBalance(ctx context.Context, storeID uint, amount float64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Store{}).
		Where("id = ?", storeID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *storeRepository) Update(ctx context.Context, store *domain.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}
