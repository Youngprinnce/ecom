package order

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/youngprinnce/go-ecom/middleware"
	"github.com/youngprinnce/go-ecom/types"
	"github.com/youngprinnce/go-ecom/utils"
)

type Handler struct {
	productStore types.ProductStore
	orderStore   types.OrderStore
	userStore    types.UserStore
}

func NewHandler(productStore types.ProductStore, orderStore types.OrderStore, userStore types.UserStore) *Handler {
	return &Handler{
		productStore: productStore,
		orderStore:   orderStore,
		userStore:    userStore,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Apply JWTAuth middleware to all order routes
	orderRouter := router.Group("/orders")
	orderRouter.Use(middleware.JWTAuth())

	// Authenticated routes
	orderRouter.GET("", h.handleGetOrders)
	orderRouter.POST("", h.handleCreateOrder)
	orderRouter.DELETE("/:id", h.handleCancelOrder)

	// Admin-only route
	adminRouter := orderRouter.Group("/:id/status")
	adminRouter.Use(middleware.AdminOnly())
	adminRouter.PUT("", h.handleUpdateOrderStatus)
}

// handleCheckout handles the checkout process for the cart.
func (h *Handler) handleCreateOrder(c *gin.Context) {
	// Retrieve userID from the request context
	userID, exists := c.Get(string(middleware.UserKey))
	if !exists {
		utils.WriteError(c.Writer, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	// Parse the request payload
	var payload types.CartCheckoutPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if len(payload.Items) == 0 {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("cart is empty"))
		return
	}

	// Get product IDs from the cart items
	productIDs := getCartItemsProductIDs(payload.Items)

	// Fetch products from the database
	products, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Create the order
	orderID, total, err := h.createOrder(products, payload.Items, userID.(int))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Return the order ID and total price
	utils.WriteJSON(c.Writer, http.StatusOK, map[string]interface{}{
		"orderID":    orderID,
		"totalPrice": total,
	})
}

func (h *Handler) handleGetOrders(c *gin.Context) {
	// Retrieve userID from the request context
	userID, exists := c.Get(string(middleware.UserKey))
	if !exists {
		utils.WriteError(c.Writer, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	// Fetch orders for the user
	orders, err := h.orderStore.GetOrdersByUserID(userID.(int))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Return the orders
	utils.WriteJSON(c.Writer, http.StatusOK, orders)
}

func (h *Handler) handleUpdateOrderStatus(c *gin.Context) {
	// Extract order ID from the URL
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid order ID"))
		return
	}

	// Parse the request payload
	var payload types.UpdateOrderStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	// Update the order status
	if err := h.orderStore.UpdateOrderStatus(orderID, payload.Status); err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	utils.WriteJSON(c.Writer, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) handleCancelOrder(c *gin.Context) {
	// Retrieve userID from the request context
	userID, exists := c.Get(string(middleware.UserKey))
	if !exists {
		utils.WriteError(c.Writer, http.StatusUnauthorized, fmt.Errorf("invalid user ID"))
		return
	}

	// Extract order ID from the URL
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid order ID"))
		return
	}

	order, err := h.orderStore.GetOrderByID(orderID)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("order not found"))
		return
	}

	if order.Status != "pending" {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("order is not pending, can't cancel"))
		return
	}

	// restore product quantities
	orderItems, err := h.orderStore.GetOrderItemsByOrderID(orderID)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Create a map of products for quick lookup
	productMap := make(map[int]*types.Product)
	for _, item := range orderItems {
		product, err := h.productStore.GetProductByID(item.ProductID)
		if err != nil {
			utils.WriteError(c.Writer, http.StatusInternalServerError, err)
			return
		}
		productMap[product.ID] = product
	}

	// Update product quantities
	for _, item := range orderItems {
		product := productMap[item.ProductID]
		product.Quantity += item.Quantity

		if err := h.productStore.UpdateProduct(*product); err != nil {
			utils.WriteError(c.Writer, http.StatusInternalServerError, err)
			return
		}
	}

	// Cancel the order
	if err := h.orderStore.CancelOrder(orderID, userID.(int)); err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

// getCartItemsProductIDs extracts product IDs from the cart items.
func getCartItemsProductIDs(cartItems []types.CartCheckoutItem) []int {
	productIDs := make([]int, len(cartItems))
	for i, item := range cartItems {
		productIDs[i] = item.ProductID
	}
	return productIDs
}

// createOrder creates an order and updates product quantities.
func (h *Handler) createOrder(products []types.Product, items []types.CartCheckoutItem, userID int) (int, float64, error) {
	// Create a map of products for quick lookup
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Validate product availability
	if err := checkIfProductIsInStock(productMap, items); err != nil {
		return 0, 0, err
	}

	// Calculate the total price
	totalPrice := calculateTotalPrice(productMap, items)

	// Update product quantities
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Create the order in the database
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "default address", // Replace with actual address logic
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range items {
		product := productMap[item.ProductID]
		if err := h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}); err != nil {
			return 0, 0, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return orderID, totalPrice, nil
}

// checkIfProductIsInStock ensures all products in the cart are in stock.
func checkIfProductIsInStock(productMap map[int]types.Product, cartItems []types.CartCheckoutItem) error {
	for _, item := range cartItems {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d not found", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is out of stock", product.ID)
		}
	}
	return nil
}

// calculateTotalPrice calculates the total price of the cart.
func calculateTotalPrice(productMap map[int]types.Product, cartItems []types.CartCheckoutItem) float64 {
	total := 0.0
	for _, item := range cartItems {
		product := productMap[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}
	return total
}
