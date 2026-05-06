// package handler

// import (
// 	"encoding/json"
// 	"net/http"

// 	"bookstore/internal/model"
// 	"bookstore/internal/repository"

// 	"golang.org/x/crypto/bcrypt"
// )

// type UserHandler struct {
// 	userRepo repository.UserRepo
// }

// func NewUserHandler(userRepo repository.UserRepo) *UserHandler {
// 	return &UserHandler{userRepo: userRepo}
// }

// type registerRequest struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	var req registerRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		errorResponse(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	if req.Name == "" || req.Email == "" || req.Password == "" {
// 		errorResponse(w, http.StatusBadRequest, "Name, email and password are required")
// 		return
// 	}

// 	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		errorResponse(w, http.StatusInternalServerError, "Failed to process password")
// 		return
// 	}

// 	user := &model.User{
// 		Name:     req.Name,
// 		Email:    req.Email,
// 		Password: string(hashed),
// 	}

// 	if err := h.userRepo.Create(user); err != nil {
// 		errorResponse(w, http.StatusConflict, "Email already exists")
// 		return
// 	}

// 	successResponse(w, http.StatusCreated, "Account created successfully", map[string]any{
// 		"id":         user.ID,
// 		"name":       user.Name,
// 		"email":      user.Email,
// 		"created_at": user.CreatedAt,
// 	})
// }
