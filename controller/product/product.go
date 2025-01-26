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

// handleGetProducts retrieves all products (public access)
func (h *Handler) handleGetProducts(c *gin.Context) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(c.Writer, http.StatusOK, products)
}

// handleCreateProduct creates a new product (admin only)
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

// handleUpdateProduct updates an existing product (admin only)
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

// handleDeleteProduct deletes a product (admin only)
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
