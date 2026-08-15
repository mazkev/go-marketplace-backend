package usecase

import (
	"context"
	"errors"

	"go-market/internal/domain"
	"go-market/pkg/hash"
	"go-market/pkg/jwt"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrUserNotFound      = errors.New("user not found")
)

type authUsecase struct {
	userRepo   domain.UserRepository
	jwtService jwt.JWTService
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtService jwt.JWTService) domain.AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existingUser, _ := u.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Role:         domain.RoleBuyer,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := u.jwtService.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredential
	}

	tokenPair, err := u.jwtService.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *authUsecase) AddAddress(ctx context.Context, userID uint, req *domain.CreateAddressRequest) (*domain.UserAddress, error) {
	address := &domain.UserAddress{
		UserID:       userID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		FullAddress:  req.FullAddress,
		CityID:       req.CityID,
		IsPrimary:    req.IsPrimary,
	}

	if err := u.userRepo.CreateAddress(ctx, address); err != nil {
		return nil, err
	}
	return address, nil
}

func (u *authUsecase) GetAddresses(ctx context.Context, userID uint) ([]domain.UserAddress, error) {
	return u.userRepo.GetAddressesByUserID(ctx, userID)
}
