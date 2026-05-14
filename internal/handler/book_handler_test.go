package handler_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookstore/internal/handler"
	"bookstore/internal/model"
	"bookstore/internal/repository/mocks"
)

func TestGetAllBooks_Success(t *testing.T) {
	mockRepo := &mocks.MockBookRepo{
		GetAllFn: func(genre string) ([]model.Book, error) {
			return []model.Book{
				{Base: model.Base{ID: "1"}, Title: "Clean Code", Genre: "Technology"},
				{Base: model.Base{ID: "2"}, Title: "Atomic Habits", Genre: "Self-Help"},
			}, nil
		},
	}

	h := handler.NewBookHandler(mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	h.GetAllBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 books, got %d", len(data))
	}
}

func TestGetAllBooks_FilterByGenre(t *testing.T) {
	mockRepo := &mocks.MockBookRepo{
		GetAllFn: func(genre string) ([]model.Book, error) {
			if genre == "Technology" {
				return []model.Book{
					{Base: model.Base{ID: "1"}, Title: "Clean Code", Genre: "Technology"},
				}, nil
			}
			return []model.Book{}, nil
		},
	}

	h := handler.NewBookHandler(mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/books?genre=Technology", nil)
	w := httptest.NewRecorder()

	h.GetAllBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp["data"].([]any)
	if len(data) != 1 {
		t.Errorf("expected 1 book, got %d", len(data))
	}
}

func TestGetBookByID_Success(t *testing.T) {
	mockRepo := &mocks.MockBookRepo{
		GetByIDFn: func(id string) (*model.Book, error) {
			return &model.Book{
				Base:  model.Base{ID: id},
				Title: "Clean Code",
			}, nil
		},
	}

	h := handler.NewBookHandler(mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	w := httptest.NewRecorder()

	h.GetBookByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetBookByID_NotFound(t *testing.T) {
	mockRepo := &mocks.MockBookRepo{
		GetByIDFn: func(id string) (*model.Book, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := handler.NewBookHandler(mockRepo)
	req := httptest.NewRequest(http.MethodGet, "/books/invalid-id", nil)
	w := httptest.NewRecorder()

	h.GetBookByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
