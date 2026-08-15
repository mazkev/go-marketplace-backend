package domain

import (
	"context"
	"time"
)

type WithdrawalStatus string

const (
	WithdrawalStatusPending  WithdrawalStatus = "PENDING"
	WithdrawalStatusApproved WithdrawalStatus = "APPROVED"
	WithdrawalStatusRejected WithdrawalStatus = "REJECTED"
)

type MutationType string

const (
	MutationTypeCredit MutationType = "CREDIT" // Money added (e.g. escrow released)
	MutationTypeDebit  MutationType = "DEBIT"  // Money deducted (e.g. withdrawal)
)

type Withdrawal struct {
	ID            uint             `gorm:"primaryKey" json:"id"`
	StoreID       uint             `gorm:"not null;index" json:"store_id"`
	BankName      string           `gorm:"size:50;not null" json:"bank_name"`
	AccountNumber string           `gorm:"size:50;not null" json:"account_number"`
	AccountHolder string           `gorm:"size:100;not null" json:"account_holder"`
	Amount        float64          `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status        WithdrawalStatus `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`

	// Associations
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

type BalanceMutation struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	StoreID     uint         `gorm:"not null;index" json:"store_id"`
	Amount      float64      `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type        MutationType `gorm:"size:10;not null" json:"type"`
	Description string       `gorm:"type:text" json:"description"`
	BalanceAfter float64     `gorm:"type:decimal(15,2);not null" json:"balance_after"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CreateWithdrawalRequest struct {
	BankName      string  `json:"bank_name" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required,min=5,max=30"`
	AccountHolder string  `json:"account_holder" binding:"required,min=3,max=100"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

type WalletRepository interface {
	CreateWithdrawal(ctx context.Context, withdrawal *Withdrawal) error
	GetWithdrawalsByStoreID(ctx context.Context, storeID uint) ([]Withdrawal, error)
	RecordMutation(ctx context.Context, mutation *BalanceMutation) error
	GetMutationsByStoreID(ctx context.Context, storeID uint) ([]BalanceMutation, error)
}

type WalletUsecase interface {
	RequestWithdrawal(ctx context.Context, userID uint, req *CreateWithdrawalRequest) (*Withdrawal, error)
	GetStoreWithdrawals(ctx context.Context, userID uint) ([]Withdrawal, error)
	GetStoreMutations(ctx context.Context, userID uint) ([]BalanceMutation, error)
}
