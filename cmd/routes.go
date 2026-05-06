package cmd

import (
	"database/sql"
	"net/http"

	"bookstore/internal/config"
	"bookstore/internal/handler"
	"bookstore/internal/middleware"
	"bookstore/internal/repository"
)

func SetupRoutes(db *sql.DB, cfg *config.Config) http.Handler {
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo, cfg.JWTSecret)

	auth := middleware.Auth(cfg.JWTSecret)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	// Protected routes
	mux.Handle("GET /user", auth(http.HandlerFunc(userHandler.GetProfile)))
	mux.Handle("DELETE /user-delete", auth(http.HandlerFunc(userHandler.Delete)))

	return mux
}
