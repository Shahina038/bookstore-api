package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookstore/internal/handler"
	"bookstore/internal/middleware"
	"bookstore/internal/model"
	"bookstore/internal/repository/mocks"
)

func orderContextRequest(method, url, userID string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestCheckout_Success(t *testing.T) {
	mockRepo := &mocks.MockOrderRepo{
		CreateOrderFn: func(userID string) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: "order-1"},
				UserID:      userID,
				TotalAmount: 89.98,
				Status:      "pending",
			}, nil
		},
	}

	h := handler.NewOrderHandler(mockRepo)
	req := orderContextRequest(http.MethodPost, "/orders", "user-1")
	w := httptest.NewRecorder()

	h.Checkout(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestCheckout_EmptyCart(t *testing.T) {
	mockRepo := &mocks.MockOrderRepo{
		CreateOrderFn: func(userID string) (*model.Order, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := handler.NewOrderHandler(mockRepo)
	req := orderContextRequest(http.MethodPost, "/orders", "user-1")
	w := httptest.NewRecorder()

	h.Checkout(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetOrders_Success(t *testing.T) {
	mockRepo := &mocks.MockOrderRepo{
		GetOrdersByUserIDFn: func(userID string) ([]model.Order, error) {
			return []model.Order{
				{Base: model.Base{ID: "order-1"}, TotalAmount: 89.98, Status: "pending"},
				{Base: model.Base{ID: "order-2"}, TotalAmount: 45.99, Status: "pending"},
			}, nil
		},
	}

	h := handler.NewOrderHandler(mockRepo)
	req := orderContextRequest(http.MethodGet, "/orders", "user-1")
	w := httptest.NewRecorder()

	h.GetOrders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 orders, got %d", len(data))
	}
}

func TestGetOrderByID_Success(t *testing.T) {
	mockRepo := &mocks.MockOrderRepo{
		GetOrderByIDFn: func(id string, userID string) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				UserID:      userID,
				TotalAmount: 89.98,
				Status:      "pending",
			}, nil
		},
	}

	h := handler.NewOrderHandler(mockRepo)
	req := orderContextRequest(http.MethodGet, "/orders/order-1", "user-1")
	w := httptest.NewRecorder()

	h.GetOrderByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetOrderByID_NotFound(t *testing.T) {
	mockRepo := &mocks.MockOrderRepo{
		GetOrderByIDFn: func(id string, userID string) (*model.Order, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := handler.NewOrderHandler(mockRepo)
	req := orderContextRequest(http.MethodGet, "/orders/wrong-id", "user-1")
	w := httptest.NewRecorder()

	h.GetOrderByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
