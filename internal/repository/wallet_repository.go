package repository

import (
	"context"
	"errors"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient store balance for withdrawal")
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWithdrawal(ctx context.Context, withdrawal *domain.Withdrawal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var store domain.Store
		if err := tx.Where("id = ?", withdrawal.StoreID).First(&store).Error; err != nil {
			return err
		}

		if store.Balance < withdrawal.Amount {
			return ErrInsufficientBalance
		}

		// Deduct store balance
		if err := tx.Model(&domain.Store{}).
			Where("id = ? AND balance >= ?", withdrawal.StoreID, withdrawal.Amount).
			Update("balance", gorm.Expr("balance - ?", withdrawal.Amount)).Error; err != nil {
			return err
		}

		// Create withdrawal record
		if err := tx.Create(withdrawal).Error; err != nil {
			return err
		}

		// Record mutation ledger
		newBalance := store.Balance - withdrawal.Amount
		mutation := domain.BalanceMutation{
			StoreID:      withdrawal.StoreID,
			Amount:       withdrawal.Amount,
			Type:         domain.MutationTypeDebit,
			Description:  "Store balance withdrawal to " + withdrawal.BankName + " - " + withdrawal.AccountNumber,
			BalanceAfter: newBalance,
		}

		return tx.Create(&mutation).Error
	})
}

func (r *walletRepository) GetWithdrawalsByStoreID(ctx context.Context, storeID uint) ([]domain.Withdrawal, error) {
	var withdrawals []domain.Withdrawal
	err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Order("id DESC").
		Find(&withdrawals).Error
	return withdrawals, err
}

func (r *walletRepository) RecordMutation(ctx context.Context, mutation *domain.BalanceMutation) error {
	return r.db.WithContext(ctx).Create(mutation).Error
}

func (r *walletRepository) GetMutationsByStoreID(ctx context.Context, storeID uint) ([]domain.BalanceMutation, error) {
	var mutations []domain.BalanceMutation
	err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Order("id DESC").
		Find(&mutations).Error
	return mutations, err
}
