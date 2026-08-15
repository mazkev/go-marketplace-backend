package usecase

import (
	"context"
	"errors"
	"time"

	"go-market/internal/domain"
)

var (
	ErrVoucherNotFound     = errors.New("voucher not found or inactive")
	ErrVoucherExpired      = errors.New("voucher is not within valid active period")
	ErrVoucherQuotaFull    = errors.New("voucher quota has been exhausted")
	ErrVoucherMinSpend     = errors.New("order amount does not meet voucher minimum spend requirement")
	ErrVoucherStoreMismatch = errors.New("voucher is not valid for this store")
)

type voucherUsecase struct {
	voucherRepo domain.VoucherRepository
	storeRepo   domain.StoreRepository
}

func NewVoucherUsecase(voucherRepo domain.VoucherRepository, storeRepo domain.StoreRepository) domain.VoucherUsecase {
	return &voucherUsecase{
		voucherRepo: voucherRepo,
		storeRepo:   storeRepo,
	}
}

func (u *voucherUsecase) CreateStoreVoucher(ctx context.Context, userID uint, req *domain.CreateVoucherRequest) (*domain.Voucher, error) {
	store, err := u.storeRepo.GetByUserID(ctx, userID)
	if err != nil || store == nil {
		return nil, ErrUnauthorizedStore
	}

	voucher := &domain.Voucher{
		Code:            req.Code,
		StoreID:         &store.ID,
		VoucherType:     req.VoucherType,
		DiscountPercent: req.DiscountPercent,
		DiscountAmount:  req.DiscountAmount,
		MaxDiscount:     req.MaxDiscount,
		MinSpend:        req.MinSpend,
		Quota:           req.Quota,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsActive:        true,
	}

	if err := u.voucherRepo.Create(ctx, voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (u *voucherUsecase) CreatePlatformVoucher(ctx context.Context, req *domain.CreateVoucherRequest) (*domain.Voucher, error) {
	voucher := &domain.Voucher{
		Code:            req.Code,
		StoreID:         nil, // platform-wide
		VoucherType:     req.VoucherType,
		DiscountPercent: req.DiscountPercent,
		DiscountAmount:  req.DiscountAmount,
		MaxDiscount:     req.MaxDiscount,
		MinSpend:        req.MinSpend,
		Quota:           req.Quota,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsActive:        true,
	}

	if err := u.voucherRepo.Create(ctx, voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (u *voucherUsecase) GetAvailableVouchers(ctx context.Context, storeID *uint) ([]domain.Voucher, error) {
	return u.voucherRepo.ListAvailable(ctx, storeID)
}

func (u *voucherUsecase) ApplyVoucher(ctx context.Context, req *domain.ApplyVoucherRequest) (*domain.ApplyVoucherResponse, error) {
	voucher, err := u.voucherRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, ErrVoucherNotFound
	}

	now := time.Now()
	if now.Before(voucher.StartDate) || now.After(voucher.EndDate) {
		return nil, ErrVoucherExpired
	}

	if voucher.UsedCount >= voucher.Quota {
		return nil, ErrVoucherQuotaFull
	}

	if req.TotalAmount < voucher.MinSpend {
		return nil, ErrVoucherMinSpend
	}

	if voucher.StoreID != nil {
		if req.StoreID == nil || *voucher.StoreID != *req.StoreID {
			return nil, ErrVoucherStoreMismatch
		}
	}

	var discount float64
	if voucher.VoucherType == domain.VoucherTypePercentage && voucher.DiscountPercent != nil {
		discount = (req.TotalAmount * (*voucher.DiscountPercent)) / 100.0
		if voucher.MaxDiscount != nil && discount > *voucher.MaxDiscount {
			discount = *voucher.MaxDiscount
		}
	} else if voucher.VoucherType == domain.VoucherTypeFixed && voucher.DiscountAmount != nil {
		discount = *voucher.DiscountAmount
		if discount > req.TotalAmount {
			discount = req.TotalAmount
		}
	}

	finalAmount := req.TotalAmount - discount
	if finalAmount < 0 {
		finalAmount = 0
	}

	return &domain.ApplyVoucherResponse{
		Valid:          true,
		VoucherCode:    voucher.Code,
		DiscountAmount: discount,
		FinalAmount:    finalAmount,
		Message:        "Voucher applied successfully",
		Voucher:        voucher,
	}, nil
}
