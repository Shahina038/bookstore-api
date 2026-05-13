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

	adminHandler := handler.NewAdminHandler(userRepo, orderRepo, bookRepo)

admin := middleware.AdminOnly(userRepo)

// Admin routes (JWT + admin role required)

// mux.Handle("PATCH /admin/users/{id}/make-admin", auth(admin(http.HandlerFunc(adminHandler.MakeAdmin))))


// Book management — admin only


	auth := middleware.Auth(cfg.JWTSecret)

	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle("GET /user", auth(http.HandlerFunc(userHandler.GetProfile)))
	mux.Handle("DELETE /user-delete", auth(http.HandlerFunc(userHandler.Delete)))
	mux.Handle("GET /admin/users", auth(admin(http.HandlerFunc(adminHandler.GetAllUsers))))
	mux.Handle("GET /admin/users/{id}", auth(admin(http.HandlerFunc(adminHandler.GetUserByID))))

	// Book routes
	mux.Handle("GET /admin/books", auth(admin(http.HandlerFunc(adminHandler.GetAllBooks))))
	mux.Handle("GET /admin/books/{id}", auth(admin(http.HandlerFunc(adminHandler.GetBookByID))))
	mux.Handle("POST /admin/books", auth(admin(http.HandlerFunc(adminHandler.CreateBook))))
	mux.Handle("PATCH /admin/books/{id}", auth(admin(http.HandlerFunc(adminHandler.UpdateBook))))
	mux.Handle("DELETE /admin/books/{id}", auth(admin(http.HandlerFunc(adminHandler.DeleteBook))))

	// for User
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
	mux.Handle("GET /admin/orders", auth(admin(http.HandlerFunc(adminHandler.GetAllOrders))))

	return mux
}
