package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	authResp, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", authResp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid credentials format", err.Error())
		return
	}

	authResp, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", authResp)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	user, err := h.authUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved", user)
}

func (h *AuthHandler) AddAddress(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid address data", err.Error())
		return
	}

	address, err := h.authUsecase.AddAddress(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(c, "Failed to create address", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Address added successfully", address)
}

func (h *AuthHandler) GetAddresses(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	addresses, err := h.authUsecase.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get addresses", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Addresses retrieved successfully", addresses)
}
