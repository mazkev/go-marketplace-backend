package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/middleware"
	"go-market/internal/domain"
	"go-market/pkg/response"
)

type ProductHandler struct {
	productUsecase  domain.ProductUsecase
	categoryRepo    domain.CategoryRepository
}

func NewProductHandler(productUsecase domain.ProductUsecase, categoryRepo domain.CategoryRepository) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		categoryRepo:   categoryRepo,
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var filter domain.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	products, totalRows, err := h.productUsecase.ListProducts(c.Request.Context(), &filter)
	if err != nil {
		response.InternalServerError(c, "Failed to list products", err.Error())
		return
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	}
	totalPage := int(math.Ceil(float64(totalRows) / float64(limit)))

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Products retrieved successfully", products, meta)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	product, err := h.productUsecase.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product details retrieved", product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid product payload", err.Error())
		return
	}

	product, err := h.productUsecase.CreateProduct(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid product update payload", err.Error())
		return
	}

	product, err := h.productUsecase.UpdateProduct(c.Request.Context(), userID, uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	if err := h.productUsecase.DeleteProduct(c.Request.Context(), userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryRepo.GetAll(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to get categories", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Categories retrieved successfully", categories)
}
