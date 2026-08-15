package config

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"go-market/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// Auto-Migrate all domain models
	err = db.AutoMigrate(
		&domain.User{},
		&domain.UserAddress{},
		&domain.Store{},
		&domain.Category{},
		&domain.Product{},
		&domain.ProductVariant{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Review{},
		&domain.Voucher{},
		&domain.Withdrawal{},
		&domain.BalanceMutation{},
		&domain.Wishlist{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database migration completed successfully.")
	seedInitialData(db)
	return db, nil
}

func seedInitialData(db *gorm.DB) {
	var count int64
	db.Model(&domain.Category{}).Count(&count)
	if count == 0 {
		categories := []domain.Category{
			{Name: "Elektronik", Slug: "elektronik"},
			{Name: "Fashion Pria", Slug: "fashion-pria"},
			{Name: "Fashion Wanita", Slug: "fashion-wanita"},
			{Name: "Handphone & Tablet", Slug: "handphone-tablet"},
			{Name: "Komputer & Laptop", Slug: "komputer-laptop"},
			{Name: "Makanan & Minuman", Slug: "makanan-minuman"},
		}
		for _, cat := range categories {
			db.Create(&cat)
		}
		log.Println("Seeded initial categories.")
	}

	var voucherCount int64
	db.Model(&domain.Voucher{}).Count(&voucherCount)
	if voucherCount == 0 {
		percent := 10.0
		maxDisc := 50000.0
		now := time.Now()
		vouchers := []domain.Voucher{
			{
				Code:            "DISKON10",
				VoucherType:     domain.VoucherTypePercentage,
				DiscountPercent: &percent,
				MaxDiscount:     &maxDisc,
				MinSpend:        100000.0,
				Quota:           100,
				StartDate:       now,
				EndDate:         now.AddDate(1, 0, 0),
				IsActive:        true,
			},
		}
		for _, v := range vouchers {
			db.Create(&v)
		}
		log.Println("Seeded initial platform vouchers.")
	}
}
