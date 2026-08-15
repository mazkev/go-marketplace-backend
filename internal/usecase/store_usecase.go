package usecase

import (
	"context"
	"errors"

	"go-market/internal/domain"
)

var (
	ErrStoreAlreadyExists = errors.New("user already owns a store")
	ErrSlugAlreadyTaken   = errors.New("store domain slug is already taken")
	ErrStoreNotFound      = errors.New("store not found")
	ErrUnauthorizedRole   = errors.New("akun dengan role pembeli (buyer) tidak diizinkan membuka toko seller, silakan daftar sebagai akun penjual (seller)")
)

type storeUsecase struct {
	storeRepo domain.StoreRepository
	userRepo  domain.UserRepository
}

func NewStoreUsecase(storeRepo domain.StoreRepository, userRepo domain.UserRepository) domain.StoreUsecase {
	return &storeUsecase{
		storeRepo: storeRepo,
		userRepo:  userRepo,
	}
}

func (u *storeUsecase) CreateStore(ctx context.Context, userID uint, req *domain.CreateStoreRequest) (*domain.StoreResponse, error) {
	// Validate that the user has seller role
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrStoreNotFound
	}
	if user.Role != domain.RoleSeller {
		return nil, ErrUnauthorizedRole
	}

	// Check if user already has a store
	existingStore, _ := u.storeRepo.GetByUserID(ctx, userID)
	if existingStore != nil {
		return nil, ErrStoreAlreadyExists
	}

	// Check domain slug
	existingSlug, _ := u.storeRepo.GetBySlug(ctx, req.DomainSlug)
	if existingSlug != nil {
		return nil, ErrSlugAlreadyTaken
	}

	store := &domain.Store{
		UserID:     userID,
		StoreName:  req.StoreName,
		DomainSlug: req.DomainSlug,
		CityID:     req.CityID,
		Balance:    0.0,
	}

	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	// Update user role to seller if currently buyer
	_ = u.userRepo.UpdateRole(ctx, userID, domain.RoleSeller)

	resp := store.ToResponse()
	return &resp, nil
}

func (u *storeUsecase) GetStoreProfile(ctx context.Context, storeID uint) (*domain.Store, error) {
	store, err := u.storeRepo.GetByID(ctx, storeID)
	if err != nil {
		return nil, ErrStoreNotFound
	}
	return store, nil
}

func (u *storeUsecase) GetStoreByUserID(ctx context.Context, userID uint) (*domain.Store, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrStoreNotFound
	}
	return store, nil
}

func (u *storeUsecase) GetStoreBalance(ctx context.Context, userID uint) (float64, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, ErrStoreNotFound
	}
	return store.Balance, nil
}
