package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type StoreHandler struct {
	storeUsecase domain.StoreUsecase
}

func NewStoreHandler(storeUsecase domain.StoreUsecase) *StoreHandler {
	return &StoreHandler{storeUsecase: storeUsecase}
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid store registration data", err.Error())
		return
	}

	storeResp, err := h.storeUsecase.CreateStore(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Store created successfully", storeResp)
}

func (h *StoreHandler) GetStoreProfile(c *gin.Context) {
	idParam := c.Param("id")
	storeID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid store ID", nil)
		return
	}

	store, err := h.storeUsecase.GetStoreProfile(c.Request.Context(), uint(storeID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Store profile retrieved", store)
}

func (h *StoreHandler) GetMyStore(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	store, err := h.storeUsecase.GetStoreByUserID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "My store retrieved", store)
}

func (h *StoreHandler) GetStoreBalance(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	balance, err := h.storeUsecase.GetStoreBalance(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Store balance retrieved", gin.H{
		"balance": balance,
	})
}
