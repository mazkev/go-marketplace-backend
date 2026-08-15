package domain

import (
	"context"
	"time"
)

type Store struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	StoreName  string    `gorm:"size:255;not null" json:"store_name"`
	DomainSlug string    `gorm:"size:255;uniqueIndex;not null" json:"domain_slug"`
	CityID     uint      `gorm:"not null" json:"city_id"`
	Balance    float64   `gorm:"type:decimal(15,2);default:0.00;not null" json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Associations
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Products []Product `gorm:"foreignKey:StoreID" json:"products,omitempty"`
}

type CreateStoreRequest struct {
	StoreName  string `json:"store_name" binding:"required,min=3,max=100"`
	DomainSlug string `json:"domain_slug" binding:"required,min=3,max=50,alphanum"`
	CityID     uint   `json:"city_id" binding:"required"`
}

type StoreResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	StoreName  string    `json:"store_name"`
	DomainSlug string    `json:"domain_slug"`
	CityID     uint      `json:"city_id"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *Store) ToResponse() StoreResponse {
	return StoreResponse{
		ID:         s.ID,
		UserID:     s.UserID,
		StoreName:  s.StoreName,
		DomainSlug: s.DomainSlug,
		CityID:     s.CityID,
		Balance:    s.Balance,
		CreatedAt:  s.CreatedAt,
	}
}

type StoreRepository interface {
	Create(ctx context.Context, store *Store) error
	GetByID(ctx context.Context, id uint) (*Store, error)
	GetByUserID(ctx context.Context, userID uint) (*Store, error)
	GetBySlug(ctx context.Context, slug string) (*Store, error)
	AddBalance(ctx context.Context, storeID uint, amount float64) error
	Update(ctx context.Context, store *Store) error
}

type StoreUsecase interface {
	CreateStore(ctx context.Context, userID uint, req *CreateStoreRequest) (*StoreResponse, error)
	GetStoreProfile(ctx context.Context, storeID uint) (*Store, error)
	GetStoreByUserID(ctx context.Context, userID uint) (*Store, error)
	GetStoreBalance(ctx context.Context, userID uint) (float64, error)
}
