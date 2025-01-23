package product

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/youngprinnce/go-ecom/types"
	"github.com/youngprinnce/go-ecom/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p types.CreateProductPayload
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.CreateProduct(p); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, p)
}
