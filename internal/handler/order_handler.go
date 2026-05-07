package handler

import (
	"database/sql"
	"net/http"

	"bookstore/internal/middleware"
	"bookstore/internal/repository"
	"bookstore/internal/model"
)

type OrderHandler struct {
	orderRepo repository.OrderRepo
}

func NewOrderHandler(orderRepo repository.OrderRepo) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	order, err := h.orderRepo.CreateOrder(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusBadRequest, "Cart is empty")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Failed to place order")
		return
	}

	successResponse(w, http.StatusCreated, "Order placed successfully", order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	orders, err := h.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	if orders == nil {
		orders = []model.Order{}
	}

	successResponse(w, http.StatusOK, "Orders fetched successfully", orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id := r.PathValue("id")

	order, err := h.orderRepo.GetOrderByID(id, userID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Order not found")
		return
	}

	successResponse(w, http.StatusOK, "Order fetched successfully", order)
}
