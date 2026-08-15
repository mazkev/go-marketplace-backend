package domain

import (
	"context"
	"time"
)

type VoucherType string

const (
	VoucherTypePercentage VoucherType = "percentage"
	VoucherTypeFixed      VoucherType = "fixed"
)

type Voucher struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	Code            string      `gorm:"size:50;uniqueIndex;not null" json:"code"`
	StoreID         *uint       `gorm:"index" json:"store_id,omitempty"` // nil if platform-wide voucher
	VoucherType     VoucherType `gorm:"size:20;not null;default:'percentage'" json:"voucher_type"`
	DiscountPercent *float64    `gorm:"type:decimal(5,2)" json:"discount_percent,omitempty"`
	DiscountAmount  *float64    `gorm:"type:decimal(15,2)" json:"discount_amount,omitempty"`
	MaxDiscount     *float64    `gorm:"type:decimal(15,2)" json:"max_discount,omitempty"`
	MinSpend        float64     `gorm:"type:decimal(15,2);default:0.00;not null" json:"min_spend"`
	Quota           int         `gorm:"not null;default:100" json:"quota"`
	UsedCount       int         `gorm:"not null;default:0" json:"used_count"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	IsActive        bool        `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`

	// Associations
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

type CreateVoucherRequest struct {
	Code            string      `json:"code" binding:"required,min=3,max=30"`
	VoucherType     VoucherType `json:"voucher_type" binding:"required,oneof=percentage fixed"`
	DiscountPercent *float64    `json:"discount_percent"`
	DiscountAmount  *float64    `json:"discount_amount"`
	MaxDiscount     *float64    `json:"max_discount"`
	MinSpend        float64     `json:"min_spend" binding:"min=0"`
	Quota           int         `json:"quota" binding:"required,min=1"`
	StartDate       time.Time   `json:"start_date" binding:"required"`
	EndDate         time.Time   `json:"end_date" binding:"required"`
}

type ApplyVoucherRequest struct {
	Code        string  `json:"code" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
	StoreID     *uint   `json:"store_id"`
}

type ApplyVoucherResponse struct {
	Valid          bool     `json:"valid"`
	VoucherCode    string   `json:"voucher_code"`
	DiscountAmount float64  `json:"discount_amount"`
	FinalAmount    float64  `json:"final_amount"`
	Message        string   `json:"message"`
	Voucher        *Voucher `json:"voucher,omitempty"`
}

type VoucherRepository interface {
	Create(ctx context.Context, voucher *Voucher) error
	GetByCode(ctx context.Context, code string) (*Voucher, error)
	ListAvailable(ctx context.Context, storeID *uint) ([]Voucher, error)
	IncrementUsage(ctx context.Context, voucherID uint) error
}

type VoucherUsecase interface {
	CreateStoreVoucher(ctx context.Context, userID uint, req *CreateVoucherRequest) (*Voucher, error)
	CreatePlatformVoucher(ctx context.Context, req *CreateVoucherRequest) (*Voucher, error)
	GetAvailableVouchers(ctx context.Context, storeID *uint) ([]Voucher, error)
	ApplyVoucher(ctx context.Context, req *ApplyVoucherRequest) (*ApplyVoucherResponse, error)
}
