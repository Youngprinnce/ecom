package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/youngprinnce/go-ecom/types"
	"github.com/youngprinnce/go-ecom/utils"
	"github.com/youngprinnce/go-ecom/middleware"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Create a subrouter for product routes
	productRouter := router.Group("/products")
	productRouter.Use(middleware.JWTAuth(), middleware.AdminOnly()) // Require JWT authentication with admin privileges

	productRouter.GET("", h.handleGetProducts)
	productRouter.POST("", h.handleCreateProduct)
	productRouter.PUT("/:id", h.handleUpdateProduct)
	productRouter.DELETE("/:id", h.handleDeleteProduct)
}

// handleGetProducts retrieves all products.
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} types.Product "list of products"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /products [get]
func (h *Handler) handleGetProducts(c *gin.Context) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(c.Writer, http.StatusOK, products)
}

// handleCreateProduct creates a new product.
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body types.CreateProductPayload true "Product payload"
// @Success 201 {object} types.CreateProductPayload "created product"
// @Failure 400 {object} map[string]string "invalid request payload"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /products [post]
func (h *Handler) handleCreateProduct(c *gin.Context) {
	var p types.CreateProductPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Create the product
	if err := h.store.CreateProduct(p); err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Log the creation of a new product
	utils.Log.WithFields(logrus.Fields{
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"quantity":    p.Quantity,
	}).Info("New product created")

	utils.WriteJSON(c.Writer, http.StatusCreated, p)
}

// handleUpdateProduct updates an existing product.
// @Summary Update a product
// @Description Update an existing product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Product ID"
// @Param payload body types.CreateProductPayload true "Product payload"
// @Success 200 {object} types.Product "updated product"
// @Failure 400 {object} map[string]string "invalid product ID or payload"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /products/{id} [put]
func (h *Handler) handleUpdateProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	var payload types.CreateProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Update the product
	product := types.Product{
		ID:          productID,
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	}

	if err := h.store.UpdateProduct(product); err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Log the update of a product
	utils.Log.WithFields(logrus.Fields{
		"productID":   productID,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"quantity":    product.Quantity,
	}).Info("Product updated")

	utils.WriteJSON(c.Writer, http.StatusOK, product)
}

// handleDeleteProduct deletes a product.
// @Summary Delete a product
// @Description Delete a product (admin only)
// @Tags products
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Product ID"
// @Success 204 "no content"
// @Failure 400 {object} map[string]string "invalid product ID"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /products/{id} [delete]
func (h *Handler) handleDeleteProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Delete the product
	if err := h.store.DeleteProduct(productID); err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Log the deletion of a product
	utils.Log.WithFields(logrus.Fields{
		"productID": productID,
	}).Info("Product deleted")

	c.Writer.WriteHeader(http.StatusNoContent)
}
