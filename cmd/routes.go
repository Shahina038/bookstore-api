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
	orderRepo := repository.NewOrderRepository(db, bookRepo)
	orderHandler := handler.NewOrderHandler(orderRepo)

	cartRepo := repository.NewCartRepository(db)
	cartHandler := handler.NewCartHandler(cartRepo)

	auth := middleware.Auth(cfg.JWTSecret)

	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle("GET /user", auth(http.HandlerFunc(userHandler.GetProfile)))
	mux.Handle("DELETE /user-delete", auth(http.HandlerFunc(userHandler.Delete)))

	// Book routes
	mux.HandleFunc("POST /books", bookHandler.CreateBook)
	mux.HandleFunc("PATCH /books/{id}", bookHandler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", bookHandler.DeleteBook)
	mux.Handle("GET /books", auth(http.HandlerFunc(bookHandler.GetAllBooks)))
	mux.Handle("GET /book/{id}", auth(http.HandlerFunc(bookHandler.GetBookByID)))

	// Cart routes (token required)
	mux.Handle("GET /cart-items", auth(http.HandlerFunc(cartHandler.GetCartItems)))
	mux.Handle("POST /cart-items", auth(http.HandlerFunc(cartHandler.AddCartItem)))
	mux.Handle("PATCH /cart-item/{id}", auth(http.HandlerFunc(cartHandler.UpdateItem)))
	mux.Handle("DELETE /cart-item/{id}", auth(http.HandlerFunc(cartHandler.RemoveItem)))

	// Order routes (token required)
	mux.Handle("POST /orders", auth(http.HandlerFunc(orderHandler.Checkout)))
	mux.Handle("GET /orders", auth(http.HandlerFunc(orderHandler.GetOrders)))
	mux.Handle("GET /order/{id}", auth(http.HandlerFunc(orderHandler.GetOrderByID)))

	return mux
}
