package repository

import (
	"database/sql"
	"bookstore/internal/model"
)

type CartRepo interface {
	AddItem(item *model.CartItem) error
	GetCartByUserID(userID string) ([]model.CartItem, error)
	UpdateQuantity(id string, userID string, quantity int) error
	RemoveItem(id string, userID string) error
	GetItemByID(id string, userID string) (*model.CartItem, error)
}

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) AddItem(item *model.CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, book_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET quantity = cart_items.quantity + $3, updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, item.UserID, item.BookID, item.Quantity).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *CartRepository) GetCartByUserID(userID string) ([]model.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.book_id, ci.quantity, ci.created_at, ci.updated_at,
		       b.id, b.title, b.author, b.genre, b.price, b.stock, b.created_at, b.updated_at
		FROM cart_items ci
		JOIN books b ON b.id = ci.book_id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		var book model.Book
		err := rows.Scan(
			&item.ID, &item.UserID, &item.BookID, &item.Quantity,
			&item.CreatedAt, &item.UpdatedAt,
			&book.ID, &book.Title, &book.Author, &book.Genre,
			&book.Price, &book.Stock, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		item.Book = &book
		items = append(items, item)
	}
	return items, nil
}

func (r *CartRepository) GetItemByID(id string, userID string) (*model.CartItem, error) {
	item := &model.CartItem{}
	query := `
		SELECT id, user_id, book_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&item.ID, &item.UserID, &item.BookID, &item.Quantity,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *CartRepository) UpdateQuantity(id string, userID string, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3`

	_, err := r.db.Exec(query, quantity, id, userID)
	return err
}

func (r *CartRepository) RemoveItem(id string, userID string) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, id, userID)
	return err
}
