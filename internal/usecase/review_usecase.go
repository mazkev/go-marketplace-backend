package usecase

import (
	"context"
	"errors"

	"go-market/internal/domain"
)

var (
	ErrOrderNotCompleted   = errors.New("reviews can only be submitted for completed order items")
	ErrReviewAlreadyExists = errors.New("review already submitted for this order item")
)

type reviewUsecase struct {
	reviewRepo  domain.ReviewRepository
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewReviewUsecase(
	reviewRepo domain.ReviewRepository,
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
) domain.ReviewUsecase {
	return &reviewUsecase{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *reviewUsecase) CreateReview(ctx context.Context, userID uint, req *domain.CreateReviewRequest) (*domain.Review, error) {
	// 1. Fetch order item and check status
	orderItem, err := u.orderRepo.GetOrderItemByID(ctx, req.OrderItemID)
	if err != nil {
		return nil, ErrOrderItemNotFound
	}

	order, err := u.orderRepo.GetByID(ctx, orderItem.OrderID)
	if err != nil || order.UserID != userID {
		return nil, ErrNotOrderOwner
	}

	if orderItem.Status != domain.OrderItemStatusCompleted {
		return nil, ErrOrderNotCompleted
	}

	// 2. Check if already reviewed
	existingReview, _ := u.reviewRepo.GetByOrderItemID(ctx, req.OrderItemID)
	if existingReview != nil {
		return nil, ErrReviewAlreadyExists
	}

	review := &domain.Review{
		OrderItemID: req.OrderItemID,
		UserID:      userID,
		ProductID:   orderItem.ProductID,
		Rating:      req.Rating,
		Comment:     req.Comment,
	}

	if err := u.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// 3. Recalculate and update product rating stats
	avgRating, count, err := u.reviewRepo.GetProductRatingStats(ctx, orderItem.ProductID)
	if err == nil {
		_ = u.productRepo.UpdateRating(ctx, orderItem.ProductID, avgRating, count)
	}

	return review, nil
}

func (u *reviewUsecase) GetProductReviews(ctx context.Context, productID uint) ([]domain.Review, error) {
	return u.reviewRepo.GetByProductID(ctx, productID)
}
