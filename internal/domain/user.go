package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:255;not null" json:"name"`
	Email        string        `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string        `gorm:"size:255;not null" json:"-"`
	Phone        string        `gorm:"size:30" json:"phone"`
	Role         Role          `gorm:"size:20;not null;default:'buyer'" json:"role"`
	Store        *Store        `gorm:"foreignKey:UserID" json:"store,omitempty"`
	Addresses    []UserAddress `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type UserAddress struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	ReceiverName string    `gorm:"size:255;not null" json:"receiver_name"`
	Phone        string    `gorm:"size:30;not null" json:"phone"`
	FullAddress  string    `gorm:"type:text;not null" json:"full_address"`
	CityID       uint      `gorm:"not null" json:"city_id"`
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Request & Response DTOs
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Role     Role   `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

type CreateAddressRequest struct {
	ReceiverName string `json:"receiver_name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	FullAddress  string `json:"full_address" binding:"required"`
	CityID       uint   `json:"city_id" binding:"required"`
	IsPrimary    bool   `json:"is_primary"`
}

// Repository Interface
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdateRole(ctx context.Context, userID uint, role Role) error

	// Address
	CreateAddress(ctx context.Context, address *UserAddress) error
	GetAddressesByUserID(ctx context.Context, userID uint) ([]UserAddress, error)
	GetPrimaryAddress(ctx context.Context, userID uint) (*UserAddress, error)
	SetPrimaryAddress(ctx context.Context, userID uint, addressID uint) error
}

// Usecase Interface
type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID uint) (*User, error)
	AddAddress(ctx context.Context, userID uint, req *CreateAddressRequest) (*UserAddress, error)
	GetAddresses(ctx context.Context, userID uint) ([]UserAddress, error)
}
