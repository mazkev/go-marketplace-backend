package usecase

import (
	"context"

	"go-market/internal/domain"
)

type wishlistUsecase struct {
	wishlistRepo domain.WishlistRepository
	productRepo  domain.ProductRepository
}

func NewWishlistUsecase(wishlistRepo domain.WishlistRepository, productRepo domain.ProductRepository) domain.WishlistUsecase {
	return &wishlistUsecase{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

func (u *wishlistUsecase) AddToWishlist(ctx context.Context, userID, productID uint) error {
	_, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}

	return u.wishlistRepo.Add(ctx, userID, productID)
}

func (u *wishlistUsecase) RemoveFromWishlist(ctx context.Context, userID, productID uint) error {
	return u.wishlistRepo.Remove(ctx, userID, productID)
}

func (u *wishlistUsecase) GetUserWishlist(ctx context.Context, userID uint) ([]domain.Product, error) {
	wishlists, err := u.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, w := range wishlists {
		if w.Product != nil {
			products = append(products, *w.Product)
		}
	}

	if products == nil {
		products = []domain.Product{}
	}

	return products, nil
}
