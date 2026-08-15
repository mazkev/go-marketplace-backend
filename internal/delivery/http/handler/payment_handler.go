package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type PaymentHandler struct {
	paymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase domain.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var payload domain.PaymentWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, "Invalid payment notification payload", err.Error())
		return
	}

	if err := h.paymentUsecase.HandleWebhook(c.Request.Context(), &payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Payment notification processed successfully", gin.H{
		"invoice_number": payload.InvoiceNumber,
		"payment_status": payload.PaymentStatus,
	})
}
