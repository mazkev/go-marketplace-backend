package repository

import (
	"context"
	"time"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) domain.VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(ctx context.Context, voucher *domain.Voucher) error {
	return r.db.WithContext(ctx).Create(voucher).Error
}

func (r *voucherRepository) GetByCode(ctx context.Context, code string) (*domain.Voucher, error) {
	var voucher domain.Voucher
	err := r.db.WithContext(ctx).
		Preload("Store").
		Where("code = ? AND is_active = ?", code, true).
		First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) ListAvailable(ctx context.Context, storeID *uint) ([]domain.Voucher, error) {
	var vouchers []domain.Voucher
	now := time.Now()

	db := r.db.WithContext(ctx).
		Where("is_active = ? AND start_date <= ? AND end_date >= ? AND used_count < quota", true, now, now)

	if storeID != nil {
		db = db.Where("store_id = ? OR store_id IS NULL", *storeID)
	} else {
		db = db.Where("store_id IS NULL")
	}

	err := db.Order("id DESC").Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepository) IncrementUsage(ctx context.Context, voucherID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Voucher{}).
		Where("id = ?", voucherID).
		Update("used_count", gorm.Expr("used_count + 1")).Error
}
