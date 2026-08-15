package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type CartHandler struct {
	cartUsecase domain.CartUsecase
}

func NewCartHandler(cartUsecase domain.CartUsecase) *CartHandler {
	return &CartHandler{cartUsecase: cartUsecase}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid cart request", err.Error())
		return
	}

	cart, err := h.cartUsecase.AddToCart(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Item added to cart", cart)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	summary, err := h.cartUsecase.GetCartGroupedByStore(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve cart items", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Cart items retrieved", summary)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	cartID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cart item ID", nil)
		return
	}

	var req domain.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid quantity", err.Error())
		return
	}

	if err := h.cartUsecase.UpdateCartItem(c.Request.Context(), userID, uint(cartID), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Cart item updated", nil)
}

func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	cartID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cart item ID", nil)
		return
	}

	if err := h.cartUsecase.DeleteCartItem(c.Request.Context(), userID, uint(cartID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Cart item deleted", nil)
}
