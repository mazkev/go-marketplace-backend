package repository

import (
	"context"

	"go-market/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Store").
		Preload("Addresses").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Store").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) UpdateRole(ctx context.Context, userID uint, role domain.Role) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("role", role).Error
}

func (r *userRepository) CreateAddress(ctx context.Context, address *domain.UserAddress) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if address.IsPrimary {
			if err := tx.Model(&domain.UserAddress{}).Where("user_id = ?", address.UserID).Update("is_primary", false).Error; err != nil {
				return err
			}
		} else {
			// If it's the first address, make it primary
			var count int64
			tx.Model(&domain.UserAddress{}).Where("user_id = ?", address.UserID).Count(&count)
			if count == 0 {
				address.IsPrimary = true
			}
		}
		return tx.Create(address).Error
	})
}

func (r *userRepository) GetAddressesByUserID(ctx context.Context, userID uint) ([]domain.UserAddress, error) {
	var addresses []domain.UserAddress
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_primary DESC, id ASC").
		Find(&addresses).Error
	return addresses, err
}

func (r *userRepository) GetPrimaryAddress(ctx context.Context, userID uint) (*domain.UserAddress, error) {
	var address domain.UserAddress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_primary = ?", userID, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *userRepository) SetPrimaryAddress(ctx context.Context, userID uint, addressID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.UserAddress{}).Where("user_id = ?", userID).Update("is_primary", false).Error; err != nil {
			return err
		}
		return tx.Model(&domain.UserAddress{}).Where("id = ? AND user_id = ?", addressID, userID).Update("is_primary", true).Error
	})
}
