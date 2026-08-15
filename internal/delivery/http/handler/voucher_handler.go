package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type VoucherHandler struct {
	voucherUsecase domain.VoucherUsecase
}

func NewVoucherHandler(voucherUsecase domain.VoucherUsecase) *VoucherHandler {
	return &VoucherHandler{voucherUsecase: voucherUsecase}
}

func (h *VoucherHandler) ListAvailable(c *gin.Context) {
	var storeIDPtr *uint
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		if id, err := strconv.ParseUint(storeIDStr, 10, 32); err == nil {
			val := uint(id)
			storeIDPtr = &val
		}
	}

	vouchers, err := h.voucherUsecase.GetAvailableVouchers(c.Request.Context(), storeIDPtr)
	if err != nil {
		response.InternalServerError(c, "Failed to get vouchers", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Available vouchers retrieved", vouchers)
}

func (h *VoucherHandler) CreateStoreVoucher(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid voucher payload", err.Error())
		return
	}

	voucher, err := h.voucherUsecase.CreateStoreVoucher(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Store voucher created successfully", voucher)
}

func (h *VoucherHandler) ApplyVoucher(c *gin.Context) {
	var req domain.ApplyVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid apply voucher payload", err.Error())
		return
	}

	res, err := h.voucherUsecase.ApplyVoucher(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Voucher applied successfully", res)
}
