package usecase

import (
	"context"

	"go-market/internal/domain"
)

type walletUsecase struct {
	walletRepo domain.WalletRepository
	storeRepo  domain.StoreRepository
}

func NewWalletUsecase(walletRepo domain.WalletRepository, storeRepo domain.StoreRepository) domain.WalletUsecase {
	return &walletUsecase{
		walletRepo: walletRepo,
		storeRepo:  storeRepo,
	}
}

func (u *walletUsecase) RequestWithdrawal(ctx context.Context, userID uint, req *domain.CreateWithdrawalRequest) (*domain.Withdrawal, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	withdrawal := &domain.Withdrawal{
		StoreID:       store.ID,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountHolder: req.AccountHolder,
		Amount:        req.Amount,
		Status:        domain.WithdrawalStatusPending,
	}

	if err := u.walletRepo.CreateWithdrawal(ctx, withdrawal); err != nil {
		return nil, err
	}

	return withdrawal, nil
}

func (u *walletUsecase) GetStoreWithdrawals(ctx context.Context, userID uint) ([]domain.Withdrawal, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	return u.walletRepo.GetWithdrawalsByStoreID(ctx, store.ID)
}

func (u *walletUsecase) GetStoreMutations(ctx context.Context, userID uint) ([]domain.BalanceMutation, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	return u.walletRepo.GetMutationsByStoreID(ctx, store.ID)
}
