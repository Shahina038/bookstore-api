package handler

import (
	"encoding/json"
	"net/http"

	"bookstore/internal/middleware"
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type CartHandler struct {
	cartRepo repository.CartRepo
}

func NewCartHandler(cartRepo repository.CartRepo) *CartHandler {
	return &CartHandler{cartRepo: cartRepo}
}

func (h *CartHandler) GetCartItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	items, err := h.cartRepo.GetCartByUserID(userID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch cart")
		return
	}

	if items == nil {
		items = []model.CartItem{}
	}

	successResponse(w, http.StatusOK, "Cart fetched successfully", items)
}

func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		BookID   string `json:"book_id"`
		Quantity int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.BookID == "" {
		errorResponse(w, http.StatusBadRequest, "book_id is required")
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	item := &model.CartItem{
		UserID:   userID,
		BookID:   req.BookID,
		Quantity: req.Quantity,
	}

	if err := h.cartRepo.AddItem(item); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to add item to cart")
		return
	}

	successResponse(w, http.StatusCreated, "Item added to cart", item)
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id := r.PathValue("id")

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Quantity <= 0 {
		errorResponse(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	if err := h.cartRepo.UpdateQuantity(id, userID, req.Quantity); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update cart item")
		return
	}

	successResponse(w, http.StatusOK, "Cart item updated", nil)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id := r.PathValue("id")

	if err := h.cartRepo.RemoveItem(id, userID); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to remove item")
		return
	}

	successResponse(w, http.StatusOK, "Item removed from cart", nil)
}
