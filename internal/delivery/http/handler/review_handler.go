package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type ReviewHandler struct {
	reviewUsecase domain.ReviewUsecase
}

func NewReviewHandler(reviewUsecase domain.ReviewUsecase) *ReviewHandler {
	return &ReviewHandler{reviewUsecase: reviewUsecase}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid review payload", err.Error())
		return
	}

	review, err := h.reviewUsecase.CreateReview(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Review submitted successfully", review)
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	reviews, err := h.reviewUsecase.GetProductReviews(c.Request.Context(), uint(productID))
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve reviews", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Reviews retrieved successfully", reviews)
}
