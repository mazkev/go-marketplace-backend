package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type WalletHandler struct {
	walletUsecase domain.WalletUsecase
}

func NewWalletHandler(walletUsecase domain.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: walletUsecase}
}

func (h *WalletHandler) RequestWithdrawal(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid withdrawal payload", err.Error())
		return
	}

	withdrawal, err := h.walletUsecase.RequestWithdrawal(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Withdrawal request submitted successfully", withdrawal)
}

func (h *WalletHandler) GetStoreWithdrawals(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	withdrawals, err := h.walletUsecase.GetStoreWithdrawals(c.Request.Context(), userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Store withdrawals retrieved", withdrawals)
}

func (h *WalletHandler) GetStoreMutations(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	mutations, err := h.walletUsecase.GetStoreMutations(c.Request.Context(), userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Store balance mutations retrieved", mutations)
}
