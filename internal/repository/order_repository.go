package repository

import (
	"database/sql"
	"bookstore/internal/model"
)

type OrderRepo interface {
	CreateOrder(userID string) (*model.Order, error)
	GetOrdersByUserID(userID string) ([]model.Order, error)
	GetOrderByID(id string, userID string) (*model.Order, error)
}

type OrderRepository struct {
	db      *sql.DB
	bookRepo *BookRepository
}

func NewOrderRepository(db *sql.DB, bookRepo *BookRepository) *OrderRepository {
	return &OrderRepository{db: db, bookRepo: bookRepo}
}

func (r *OrderRepository) CreateOrder(userID string) (*model.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Get cart items
	rows, err := tx.Query(`
		SELECT ci.id, ci.book_id, ci.quantity, b.price
		FROM cart_items ci
		JOIN books b ON b.id = ci.book_id
		WHERE ci.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type cartRow struct {
		ID       string
		BookID   string
		Quantity int
		Price    float64
	}

	var cartItems []cartRow
	for rows.Next() {
		var item cartRow
		if err := rows.Scan(&item.ID, &item.BookID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	if len(cartItems) == 0 {
		return nil, sql.ErrNoRows
	}

	// 2. Calculate total
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	// 3. Create order
	order := &model.Order{
		UserID:      userID,
		TotalAmount: total,
		Status:      "pending",
	}
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total_amount, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		order.UserID, order.TotalAmount, order.Status).
		Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 4. Create order items and reduce stock
	for _, item := range cartItems {
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, book_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)`,
			order.ID, item.BookID, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}

		if err := r.bookRepo.UpdateStock(tx, item.BookID, item.Quantity); err != nil {
			return nil, err
		}
	}

	// 5. Clear cart
	_, err = tx.Exec(`DELETE FROM cart_items WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID string) ([]model.Order, error) {
	query := `
		SELECT id, user_id, total_amount, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderByID(id string, userID string) (*model.Order, error) {
	order := &model.Order{}
	err := r.db.QueryRow(`
		SELECT id, user_id, total_amount, status, created_at
		FROM orders
		WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT oi.id, oi.order_id, oi.book_id, oi.quantity, oi.unit_price,
		       b.title, b.author
		FROM order_items oi
		JOIN books b ON b.id = oi.book_id
		WHERE oi.order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		var book model.Book
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.BookID, &item.Quantity, &item.UnitPrice,
			&book.Title, &book.Author,
		)
		if err != nil {
			return nil, err
		}
		item.Book = &book
		order.Items = append(order.Items, item)
	}

	return order, nil
}
