package handler

import (
	"encoding/json"
	"net/http"

	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type AdminHandler struct {
	userRepo  repository.UserRepo
	orderRepo repository.OrderRepo
	bookRepo  repository.BookRepo
}

func NewAdminHandler(userRepo repository.UserRepo, orderRepo repository.OrderRepo, bookRepo repository.BookRepo) *AdminHandler {
	return &AdminHandler{userRepo: userRepo, orderRepo: orderRepo, bookRepo: bookRepo}
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	successResponse(w, http.StatusOK, "Users fetched successfully", users)
}

func (h *AdminHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	successResponse(w, http.StatusOK, "User fetched successfully", user)
}

// func (h *AdminHandler) MakeAdmin(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")

// 	if err := h.userRepo.MakeAdmin(id); err != nil {
// 		errorResponse(w, http.StatusInternalServerError, "Failed to update user")
// 		return
// 	}

// 	successResponse(w, http.StatusOK, "User promoted to admin", nil)
// }

func (h *AdminHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderRepo.GetAllOrders()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	successResponse(w, http.StatusOK, "Orders fetched successfully", orders)
}

type bookRequest struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

func (h *AdminHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req bookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Author == "" || req.Price <= 0 {
		errorResponse(w, http.StatusBadRequest, "Title, author and price are required")
		return
	}

	book := &model.Book{
		Title:  req.Title,
		Author: req.Author,
		Genre:  req.Genre,
		Price:  req.Price,
		Stock:  req.Stock,
	}

	if err := h.bookRepo.Create(book); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	successResponse(w, http.StatusCreated, "Book created successfully", book)
}

type patchBookRequest struct {
	Title  *string  `json:"title"`
	Author *string  `json:"author"`
	Genre  *string  `json:"genre"`
	Price  *float64 `json:"price"`
	Stock  *int     `json:"stock"`
}

func (h *AdminHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req patchBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	book, err := h.bookRepo.GetByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Genre != nil {
		book.Genre = *req.Genre
	}
	if req.Price != nil {
		book.Price = *req.Price
	}
	if req.Stock != nil {
		book.Stock = *req.Stock
	}

	if err := h.bookRepo.Update(book); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update book")
		return
	}

	successResponse(w, http.StatusOK, "Book updated successfully", book)
}


func (h *AdminHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.bookRepo.Delete(id); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	successResponse(w, http.StatusOK, "Book deleted successfully", nil)
}

func (h *AdminHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")

	books, err := h.bookRepo.GetAll(genre)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch books")
		return
	}

	successResponse(w, http.StatusOK, "Books fetched successfully", books)
}

func (h *AdminHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	book, err := h.bookRepo.GetByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	successResponse(w, http.StatusOK, "Book fetched successfully", book)
}
