package usecase

import (
	"context"
	"errors"

	"go-market/internal/domain"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrUnauthorizedStore    = errors.New("user does not have a registered store")
	ErrProductNotOwnedByStore = errors.New("product does not belong to your store")
)

type productUsecase struct {
	productRepo domain.ProductRepository
	storeRepo   domain.StoreRepository
}

func NewProductUsecase(productRepo domain.ProductRepository, storeRepo domain.StoreRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, userID uint, req *domain.CreateProductRequest) (*domain.Product, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	product := &domain.Product{
		StoreID:     store.ID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Weight:      req.Weight,
	}

	if len(req.Variants) > 0 {
		var variants []domain.ProductVariant
		for _, v := range req.Variants {
			variants = append(variants, domain.ProductVariant{
				VariantName:   v.VariantName,
				PriceOverride: v.PriceOverride,
				Stock:         v.Stock,
			})
		}
		product.Variants = variants
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return u.productRepo.GetByID(ctx, product.ID)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (u *productUsecase) ListProducts(ctx context.Context, filter *domain.ProductFilter) ([]domain.Product, int64, error) {
	return u.productRepo.List(ctx, filter)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, userID uint, productID uint, req *domain.UpdateProductRequest) (*domain.Product, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	product, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.StoreID != store.ID {
		return nil, ErrProductNotOwnedByStore
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return u.productRepo.GetByID(ctx, productID)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, userID uint, productID uint) error {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return ErrUnauthorizedStore
	}

	product, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.StoreID != store.ID {
		return ErrProductNotOwnedByStore
	}

	return u.productRepo.Delete(ctx, productID)
}
