package config

import (
	"log"

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
	SeedComprehensiveData(db)
	return db, nil
}
