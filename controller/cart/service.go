package cart

import (
	"fmt"

	"github.com/youngprinnce/go-ecom/types"
)


func getCartItemsProductIDs(cartItems []types.CartCheckoutItem) []int {
	productId := []int{}
	for _, item := range cartItems {
		if item.Quantity <= 0 {
			continue
		}
		productId = append(productId, item.ProductID)
	}
	return productId
}

func (h *Handler) createOrder(products []types.Product, items []types.CartCheckoutItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}
	if err := checkIfProductIsInStock(productMap, items); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(productMap, items)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, fmt.Errorf("could not update product: %w", err)
		}
	}
	
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID: userId,
		Total:  totalPrice,
		Status: "pending",
		Address: "address",
	})

	if err != nil {
		return 0, 0, fmt.Errorf("could not create order: %w", err)
	}

	return orderID, totalPrice, nil
}

func checkIfProductIsInStock(product map[int]types.Product, cartItems []types.CartCheckoutItem) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range cartItems {
		product, ok := product[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d not found", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is out of stock", product.ID)
		}
	}
	return nil
}

func calculateTotalPrice(product map[int]types.Product, cartItems []types.CartCheckoutItem) float64 {
	totalPrice := 0.0
	for _, item := range cartItems {
		product := product[item.ProductID]
		totalPrice += product.Price * float64(item.Quantity)
	}
	return totalPrice
}
