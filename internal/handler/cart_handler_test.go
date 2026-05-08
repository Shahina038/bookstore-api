package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookstore/internal/handler"
	"bookstore/internal/middleware"
	"bookstore/internal/model"
	"bookstore/internal/repository/mocks"
)

func cartContextRequest(method, url, body, userID string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestGetCart_Success(t *testing.T) {
	mockRepo := &mocks.MockCartRepo{
		GetCartByUserIDFn: func(userID string) ([]model.CartItem, error) {
			return []model.CartItem{
				{Base: model.Base{ID: "item-1"}, BookID: "book-1", Quantity: 2},
			}, nil
		},
	}

	h := handler.NewCartHandler(mockRepo)
	req := cartContextRequest(http.MethodGet, "/cart", "", "user-1")
	w := httptest.NewRecorder()

	h.GetCartItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].([]any)
	if len(data) != 1 {
		t.Errorf("expected 1 cart item, got %d", len(data))
	}
}

func TestGetCart_Empty(t *testing.T) {
	mockRepo := &mocks.MockCartRepo{
		GetCartByUserIDFn: func(userID string) ([]model.CartItem, error) {
			return nil, nil
		},
	}

	h := handler.NewCartHandler(mockRepo)
	req := cartContextRequest(http.MethodGet, "/cart", "", "user-1")
	w := httptest.NewRecorder()

	h.GetCartItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAddItem_Success(t *testing.T) {
	mockRepo := &mocks.MockCartRepo{
		AddItemFn: func(item *model.CartItem) error {
			item.ID = "new-item-id"
			return nil
		},
	}

	h := handler.NewCartHandler(mockRepo)
	body := `{"book_id":"book-1","quantity":2}`
	req := cartContextRequest(http.MethodPost, "/cart/items", body, "user-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AddCartItem(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestAddItem_MissingBookID(t *testing.T) {
	mockRepo := &mocks.MockCartRepo{}
	h := handler.NewCartHandler(mockRepo)

	body := `{"book_id":"","quantity":2}`
	req := cartContextRequest(http.MethodPost, "/cart/items", body, "user-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AddCartItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRemoveItem_Success(t *testing.T) {
	mockRepo := &mocks.MockCartRepo{
		RemoveItemFn: func(id string, userID string) error {
			return nil
		},
	}

	h := handler.NewCartHandler(mockRepo)
	req := cartContextRequest(http.MethodDelete, "/cart/items/item-1", "", "user-1")
	w := httptest.NewRecorder()

	h.RemoveItem(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
