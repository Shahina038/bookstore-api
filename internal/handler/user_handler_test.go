package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bookstore/internal/handler"
	"bookstore/internal/model"
	"bookstore/internal/repository/mocks"
)

func TestRegister_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{
		CreateFn: func(user *model.User) error {
			user.ID = "test-uuid"
			user.CreatedAt = time.Now()
			return nil
		},
	}

	h := handler.NewUserHandler(mockRepo, "test-secret")

	body := `{"name":"John","email":"john@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestRegister_MissingFields(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{}
	h := handler.NewUserHandler(mockRepo, "test-secret")

	body := `{"name":"","email":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{
		GetByEmailFn: func(email string) (*model.User, error) {
			// bcrypt hash of "secret123"
			return &model.User{
				Base:     model.Base{ID: "test-uuid"},
				Name:     "John",
				Email:    email,
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			}, nil
		},
	}

	h := handler.NewUserHandler(mockRepo, "test-secret")

	body := `{"email":"john@example.com","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].(map[string]any)
	if data["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{
		GetByEmailFn: func(email string) (*model.User, error) {
			return &model.User{
				Base:     model.Base{ID: "test-uuid"},
				Email:    email,
				Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			}, nil
		},
	}

	h := handler.NewUserHandler(mockRepo, "test-secret")

	body := `{"email":"john@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
