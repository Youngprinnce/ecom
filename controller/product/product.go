package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Public routes (accessible to all users)
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)

	// Protected routes (admin only)
	router.Use(middleware.JWTAuth)             // Require JWT authentication
	router.Use(middleware.AdminOnly)           // Require admin privileges

	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", h.handleUpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/products/{id}", h.handleDeleteProduct).Methods(http.MethodDelete)
}

// handleGetProducts retrieves all products (public access)
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

// handleCreateProduct creates a new product (admin only)
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p types.CreateProductPayload
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Create the product
	if err := h.store.CreateProduct(p); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Log the creation of a new product
	utils.Log.WithFields(logrus.Fields{
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"quantity":    p.Quantity,
	}).Info("New product created")

	utils.WriteJSON(w, http.StatusCreated, p)
}

// handleUpdateProduct updates an existing product (admin only)
func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
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
		utils.WriteError(w, http.StatusInternalServerError, err)
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

	utils.WriteJSON(w, http.StatusOK, product)
}

// handleDeleteProduct deletes a product (admin only)
func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Delete the product
	if err := h.store.DeleteProduct(productID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Log the deletion of a product
	utils.Log.WithFields(logrus.Fields{
		"productID": productID,
	}).Info("Product deleted")

	w.WriteHeader(http.StatusNoContent)
}
