package repository

import (
	"bookstore/internal/model"
	"database/sql"
)

type BookRepo interface {
	Create(book *model.Book) error
	GetAll(genre string) ([]model.Book, error)
	GetByID(id string) (*model.Book, error)
	Update(book *model.Book) error
	Delete(id string) error
	UpdateStock(tx *sql.Tx, bookID string, quantity int) error
}

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *model.Book) error {
	query := `
		INSERT INTO books (title, author, genre, price, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, book.Title, book.Author, book.Genre, book.Price, book.Stock).
		Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func (r *BookRepository) GetAll(genre string) ([]model.Book, error) {
	query := `
		SELECT id, title, author, genre, price, stock, created_at, updated_at
		FROM books
		WHERE ($1 = '' OR genre = $1)
		ORDER BY title ASC`

	rows, err := r.db.Query(query, genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var b model.Book
		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.Genre,
			&b.Price, &b.Stock, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetByID(id string) (*model.Book, error) {
	book := &model.Book{}
	query := `
		SELECT id, title, author, genre, price, stock, created_at, updated_at
		FROM books
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.Genre,
		&book.Price, &book.Stock, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) UpdateStock(tx *sql.Tx, bookID string, quantity int) error {
	query := `
		UPDATE books
		SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2 AND stock >= $1`
	result, err := tx.Exec(query, quantity, bookID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *BookRepository) Update(book *model.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, genre = $3, price = $4, stock = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(query, book.Title, book.Author, book.Genre, book.Price, book.Stock, book.ID).
		Scan(&book.UpdatedAt)
}

func (r *BookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
