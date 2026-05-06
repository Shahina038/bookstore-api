// package main

// import (
// 	"database/sql"
// 	"net/http"

// 	"bookstore/internal/handler"
// 	"bookstore/internal/repository"
// )

// func setupRoutes(db *sql.DB) http.Handler {
// 	userRepo := repository.NewUserRepository(db)
// 	userHandler := handler.NewUserHandler(userRepo)

// 	mux := http.NewServeMux()

// 	// Auth routes
// 	mux.HandleFunc("POST /register", userHandler.Register)

// 	return mux
// }
