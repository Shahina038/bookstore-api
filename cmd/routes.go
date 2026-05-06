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

	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	auth := middleware.Auth(cfg.JWTSecret)

	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle("GET /user", auth(http.HandlerFunc(userHandler.GetProfile)))
	mux.Handle("DELETE /user-delete", auth(http.HandlerFunc(userHandler.Delete)))

	// Book routes
	mux.Handle("POST /books", auth(http.HandlerFunc(bookHandler.CreateBook)))
	mux.HandleFunc("GET /books", bookHandler.GetAllBooks)
	mux.HandleFunc("GET /book/{id}", bookHandler.GetBookByID)
	mux.Handle("PUT /books/{id}", auth(http.HandlerFunc(bookHandler.UpdateBook)))
	mux.Handle("DELETE /books/{id}", auth(http.HandlerFunc(bookHandler.DeleteBook)))

	return mux
}
