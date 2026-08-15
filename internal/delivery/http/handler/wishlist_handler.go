package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type WishlistHandler struct {
	wishlistUsecase domain.WishlistUsecase
}

func NewWishlistHandler(wishlistUsecase domain.WishlistUsecase) *WishlistHandler {
	return &WishlistHandler{wishlistUsecase: wishlistUsecase}
}

func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("product_id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	if err := h.wishlistUsecase.AddToWishlist(c.Request.Context(), userID, uint(productID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product added to wishlist", nil)
}

func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("product_id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	if err := h.wishlistUsecase.RemoveFromWishlist(c.Request.Context(), userID, uint(productID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product removed from wishlist", nil)
}

func (h *WishlistHandler) GetUserWishlist(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	products, err := h.wishlistUsecase.GetUserWishlist(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get wishlist", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Wishlist retrieved", products)
}
