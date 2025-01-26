package cart

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/youngprinnce/go-ecom/controller/auth"
	"github.com/youngprinnce/go-ecom/types"
	"github.com/youngprinnce/go-ecom/utils"
)

type Handler struct {
	orderStore types.OrderStore
	productStore types.ProductStore
	userStore types.UserStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	productIds := getCartItemsProductIDs(cart.Items)
	product, err := h.productStore.GetProductsByIDs(productIds)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(product, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"orderID":    orderID,
		"totalPrice": totalPrice,
	})
}
