package handler

import (
	"net/http"
	"strconv"

	"github.com/Rizwank123/LahziCart/catalog-service/internal/models"
	"github.com/Rizwank123/LahziCart/shared/response"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service models.ProductService
}

func NewProductHandler(service models.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var in models.CreateProductInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	product, err := h.service.SaveProduct(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("product created", product))
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	product, err := h.service.GetProductByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("product not found"))
	}
	return c.JSON(http.StatusOK, response.Success("product found", product))
}

func (h *ProductHandler) ListProducts(c echo.Context) error {
	limit := queryInt(c, "limit", 20)
	page := queryInt(c, "page", 1)
	products, total, err := h.service.GetProducts(limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list products"))
	}
	return c.JSON(http.StatusOK, response.SuccessWithPagination("products found", products, map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}))
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	var in models.UpdateProductInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	product, err := h.service.UpdateProduct(id, in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("product updated", product))
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	if err := h.service.DeleteProdduct(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to delete product"))
	}
	return c.JSON(http.StatusOK, response.Success("product deleted", nil))
}

func (h *ProductHandler) CreateReview(c echo.Context) error {
	var in models.CreateReviewInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	product, err := h.service.SaveReview(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("review created", product))
}

func (h *ProductHandler) ListReviews(c echo.Context) error {
	productID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	reviews, err := h.service.FindReviewForProduct(productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list reviews"))
	}
	return c.JSON(http.StatusOK, response.Success("reviews found", reviews))
}

func (h *ProductHandler) UpdateReview(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid review id"))
	}
	var in models.UpdateReviewInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	review, err := h.service.UpdateReview(id, in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("review updated", review))
}

func (h *ProductHandler) DeleteReview(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid review id"))
	}
	if err := h.service.DeleteReview(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to delete review"))
	}
	return c.JSON(http.StatusOK, response.Success("review deleted", nil))
}

func (h *ProductHandler) CreateImage(c echo.Context) error {
	var in models.CreateProductImageInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	image, err := h.service.SaveProductImage(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success("product image created", image))
}

func (h *ProductHandler) ListImages(c echo.Context) error {
	productID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid product id"))
	}
	images, err := h.service.FindProductImagesByProductId(productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to list images"))
	}
	return c.JSON(http.StatusOK, response.Success("images found", images))
}

func (h *ProductHandler) UpdateImage(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid image id"))
	}
	var in models.UpdateProductImageInput
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}
	image, err := h.service.UpdateProductImage(id, in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success("product image updated", image))
}

func (h *ProductHandler) DeleteImage(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid image id"))
	}
	if err := h.service.DeleteProductImage(id); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to delete image"))
	}
	return c.JSON(http.StatusOK, response.Success("product image deleted", nil))
}

func queryInt(c echo.Context, name string, fallback int) int {
	raw := c.QueryParam(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
