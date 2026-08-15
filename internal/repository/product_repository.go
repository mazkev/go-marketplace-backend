package repository

import (
	"context"
	"errors"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrInsufficientStock = errors.New("insufficient product stock")
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Store").
		Preload("Category").
		Preload("Variants").
		Preload("Reviews.User").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, filter *domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalRows int64

	db := r.db.WithContext(ctx).Model(&domain.Product{}).
		Preload("Store").
		Preload("Category").
		Preload("Variants")

	if filter.Search != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.CategoryID != nil {
		db = db.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.StoreID != nil {
		db = db.Where("store_id = ?", *filter.StoreID)
	}
	if filter.MinPrice != nil {
		db = db.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		db = db.Where("price <= ?", *filter.MaxPrice)
	}

	if err := db.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&products).Error
	return products, totalRows, err
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepository) DeductStock(ctx context.Context, productID uint, variantID *uint, quantity int) error {
	if variantID != nil && *variantID > 0 {
		res := r.db.WithContext(ctx).
			Model(&domain.ProductVariant{}).
			Where("id = ? AND product_id = ? AND stock >= ?", *variantID, productID, quantity).
			Update("stock", gorm.Expr("stock - ?", quantity))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrInsufficientStock
		}
		return nil
	}

	res := r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (r *productRepository) AddStock(ctx context.Context, productID uint, variantID *uint, quantity int) error {
	if variantID != nil && *variantID > 0 {
		return r.db.WithContext(ctx).
			Model(&domain.ProductVariant{}).
			Where("id = ? AND product_id = ?", *variantID, productID).
			Update("stock", gorm.Expr("stock + ?", quantity)).Error
	}

	return r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepository) UpdateRating(ctx context.Context, productID uint, avgRating float64, count int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"rating_avg":   avgRating,
			"rating_count": count,
		}).Error
}
