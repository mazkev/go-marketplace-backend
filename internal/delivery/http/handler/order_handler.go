package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type OrderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(orderUsecase domain.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid checkout payload", err.Error())
		return
	}

	order, err := h.orderUsecase.Checkout(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Order created successfully", order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid order ID", nil)
		return
	}

	order, err := h.orderUsecase.GetOrderByID(c.Request.Context(), userID, uint(orderID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Order details retrieved", order)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	orders, err := h.orderUsecase.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve orders", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User orders retrieved", orders)
}

func (h *OrderHandler) GetStoreOrders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	items, err := h.orderUsecase.GetStoreOrders(c.Request.Context(), userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Store orders retrieved", items)
}

func (h *OrderHandler) ShipOrderItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	orderItemID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid order item ID", nil)
		return
	}

	var req domain.ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid shipping input", err.Error())
		return
	}

	item, err := h.orderUsecase.ShipOrderItem(c.Request.Context(), userID, uint(orderItemID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Order item shipped successfully", item)
}

func (h *OrderHandler) CompleteOrderItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	orderItemID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid order item ID", nil)
		return
	}

	item, err := h.orderUsecase.CompleteOrderItem(c.Request.Context(), userID, uint(orderItemID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Order item confirmed and completed. Escrow released to seller.", item)
}
