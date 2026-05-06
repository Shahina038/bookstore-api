package handler

import (
	"encoding/json"
	"net/http"

	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/internal/utils"
	"bookstore/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepo
	jwtSecret string
}

func NewUserHandler(userRepo repository.UserRepo, jwtSecret string) *UserHandler {
	return &UserHandler{userRepo: userRepo, jwtSecret: jwtSecret}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		errorResponse(w, http.StatusBadRequest, "Name, email and password are required")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := h.userRepo.Create(user); err != nil {
		errorResponse(w, http.StatusConflict, "Email already exists")
		return
	}

	successResponse(w, http.StatusCreated, "Account created successfully", map[string]any{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		errorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	successResponse(w, http.StatusOK, "Login successful", map[string]any{
		"token": token,
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	successResponse(w, http.StatusOK, "User profile", map[string]any{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.userRepo.SoftDelete(userID); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	successResponse(w, http.StatusOK, "Account deleted successfully", nil)
}
