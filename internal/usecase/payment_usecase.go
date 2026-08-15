package usecase

import (
	"context"
	"errors"
	"strings"

	"go-market/internal/domain"
)

var (
	ErrInvalidPaymentPayload = errors.New("invalid payment webhook payload")
)

type paymentUsecase struct {
	orderRepo domain.OrderRepository
}

func NewPaymentUsecase(orderRepo domain.OrderRepository) domain.PaymentUsecase {
	return &paymentUsecase{orderRepo: orderRepo}
}

func (u *paymentUsecase) HandleWebhook(ctx context.Context, payload *domain.PaymentWebhookPayload) error {
	order, err := u.orderRepo.GetByInvoiceNumber(ctx, payload.InvoiceNumber)
	if err != nil {
		return ErrOrderNotFound
	}

	statusUpper := strings.ToUpper(payload.PaymentStatus)
	switch statusUpper {
	case "SETTLEMENT", "PAID", "SUCCESS":
		if order.PaymentStatus == domain.PaymentStatusPaid {
			return nil // Idempotent: already paid
		}
		return u.orderRepo.UpdateOrderAndItemsStatusOnPayment(
			ctx,
			order.ID,
			domain.PaymentStatusPaid,
			domain.OrderItemStatusProcessing,
		)
	case "EXPIRE", "EXPIRED":
		if order.PaymentStatus == domain.PaymentStatusPaid {
			return nil
		}
		return u.orderRepo.UpdateOrderAndItemsStatusOnPayment(
			ctx,
			order.ID,
			domain.PaymentStatusExpired,
			domain.OrderItemStatusCancelled,
		)
	case "CANCEL", "CANCELLED", "DENY":
		if order.PaymentStatus == domain.PaymentStatusPaid {
			return nil
		}
		return u.orderRepo.UpdateOrderAndItemsStatusOnPayment(
			ctx,
			order.ID,
			domain.PaymentStatusFailed,
			domain.OrderItemStatusCancelled,
		)
	default:
		return ErrInvalidPaymentPayload
	}
}
