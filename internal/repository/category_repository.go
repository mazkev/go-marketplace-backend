package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
