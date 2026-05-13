package handler

import (
	// "encoding/json"
	"net/http"

	// "bookstore/internal/model"
	"bookstore/internal/repository"
)

type BookHandler struct {
	bookRepo repository.BookRepo
}

func NewBookHandler(bookRepo repository.BookRepo) *BookHandler {
	return &BookHandler{bookRepo: bookRepo}
}



func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")

	books, err := h.bookRepo.GetAll(genre)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch books")
		return
	}

	successResponse(w, http.StatusOK, "Books fetched successfully", books)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	book, err := h.bookRepo.GetByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	successResponse(w, http.StatusOK, "Book fetched successfully", book)
}


