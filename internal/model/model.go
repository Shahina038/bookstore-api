package model

import "time"

type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Base
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `json:"is_admin"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Book struct {
	Base
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

type CartItem struct {
	Base
	UserID   string `json:"user_id"`
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
	Book     *Book  `json:"book,omitempty"`
}

type Order struct {
	Base
	UserID      string      `json:"user_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	BookID    string  `json:"book_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Book      *Book   `json:"book,omitempty"`
}
