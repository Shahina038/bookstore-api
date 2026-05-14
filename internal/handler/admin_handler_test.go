package handler_test


import (
	"bytes"
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

func adminContextRequest(method, url, userID string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestAdminGetAllUsers_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepo{
		GetAllFn: func() ([]*model.User, error) {
			return []*model.User{
				{Base: model.Base{ID: "1"}, Name: "John", Email: "john@example.com"},
				{Base: model.Base{ID: "2"}, Name: "Jane", Email: "jane@example.com"},
			}, nil
		},
	}

	h := handler.NewAdminHandler(mockUserRepo, &mocks.MockOrderRepo{}, &mocks.MockBookRepo{})
	req := adminContextRequest(http.MethodGet, "/admin/users", "admin-1")
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 users, got %d", len(data))
	}
}

func TestAdminGetUserByID_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepo{
		GetByIDFn: func(id string) (*model.User, error) {
			return &model.User{
				Base:  model.Base{ID: id},
				Name:  "John",
				Email: "john@example.com",
			}, nil
		},
	}

	h := handler.NewAdminHandler(mockUserRepo, &mocks.MockOrderRepo{}, &mocks.MockBookRepo{})
	req := adminContextRequest(http.MethodGet, "/admin/users/1", "admin-1")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminGetUserByID_NotFound(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepo{
		GetByIDFn: func(id string) (*model.User, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := handler.NewAdminHandler(mockUserRepo, &mocks.MockOrderRepo{}, &mocks.MockBookRepo{})
	req := adminContextRequest(http.MethodGet, "/admin/users/invalid", "admin-1")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// func TestAdminMakeAdmin_Success(t *testing.T) {
// 	mockUserRepo := &mocks.MockUserRepo{
// 		MakeAdminFn: func(id string) error {
// 			return nil
// 		},
// 	}

// 	h := handler.NewAdminHandler(mockUserRepo, &mocks.MockOrderRepo{}, &mocks.MockBookRepo{})
// 	req := adminContextRequest(http.MethodPatch, "/admin/users/1/make-admin", "admin-1")
// 	w := httptest.NewRecorder()

// 	h.MakeAdmin(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("expected 200, got %d", w.Code)
// 	}
// }

func TestAdminGetAllOrders_Success(t *testing.T) {
	mockOrderRepo := &mocks.MockOrderRepo{
		GetAllOrdersFn: func() ([]model.Order, error) {
			return []model.Order{
				{Base: model.Base{ID: "order-1"}, TotalAmount: 49.99, Status: "pending"},
				{Base: model.Base{ID: "order-2"}, TotalAmount: 89.99, Status: "pending"},
			}, nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, mockOrderRepo, &mocks.MockBookRepo{})
	req := adminContextRequest(http.MethodGet, "/admin/orders", "admin-1")
	w := httptest.NewRecorder()

	h.GetAllOrders(w, req)

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

func TestAdminGetAllBooks_Success(t *testing.T) {
	mockBookRepo := &mocks.MockBookRepo{
		GetAllFn: func(genre string) ([]model.Book, error) {
			return []model.Book{
				{Base: model.Base{ID: "1"}, Title: "Clean Code", Genre: "Technology"},
			}, nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, mockBookRepo)
	req := adminContextRequest(http.MethodGet, "/admin/books", "admin-1")
	w := httptest.NewRecorder()

	h.GetAllBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminGetBookByID_Success(t *testing.T) {
	mockBookRepo := &mocks.MockBookRepo{
		GetByIDFn: func(id string) (*model.Book, error) {
			return &model.Book{
				Base:  model.Base{ID: id},
				Title: "Clean Code",
			}, nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, mockBookRepo)
	req := adminContextRequest(http.MethodGet, "/admin/books/1", "admin-1")
	w := httptest.NewRecorder()

	h.GetBookByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminCreateBook_Success(t *testing.T) {
	mockBookRepo := &mocks.MockBookRepo{
		CreateFn: func(book *model.Book) error {
			book.ID = "new-uuid"
			return nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, mockBookRepo)

	body := `{"title":"Learning Go","author":"Jon Bodner","genre":"Technology","price":45.99,"stock":30}`
	req := httptest.NewRequest(http.MethodPost, "/admin/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateBook(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestAdminCreateBook_MissingFields(t *testing.T) {
	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, &mocks.MockBookRepo{})

	body := `{"title":"","author":"","price":0}`
	req := httptest.NewRequest(http.MethodPost, "/admin/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAdminUpdateBook_Success(t *testing.T) {
	mockBookRepo := &mocks.MockBookRepo{
		GetByIDFn: func(id string) (*model.Book, error) {
			return &model.Book{
				Base:   model.Base{ID: id},
				Title:  "Old Title",
				Author: "Old Author",
				Price:  39.99,
				Stock:  10,
			}, nil
		},
		UpdateFn: func(book *model.Book) error {
			return nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, mockBookRepo)

	body := `{"title":"New Title","price":49.99}`
	req := httptest.NewRequest(http.MethodPatch, "/admin/books/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminDeleteBook_Success(t *testing.T) {
	mockBookRepo := &mocks.MockBookRepo{
		DeleteFn: func(id string) error {
			return nil
		},
	}

	h := handler.NewAdminHandler(&mocks.MockUserRepo{}, &mocks.MockOrderRepo{}, mockBookRepo)
	req := httptest.NewRequest(http.MethodDelete, "/admin/books/1", nil)
	w := httptest.NewRecorder()

	h.DeleteBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
