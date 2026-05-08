package mocks

import (
	"database/sql"
	"bookstore/internal/model"
)

type MockBookRepo struct {
	CreateFn      func(book *model.Book) error
	GetAllFn      func(genre string) ([]model.Book, error)
	GetByIDFn     func(id string) (*model.Book, error)
	UpdateFn      func(book *model.Book) error
	DeleteFn      func(id string) error
	UpdateStockFn func(tx *sql.Tx, bookID string, quantity int) error
}

func (m *MockBookRepo) Create(book *model.Book) error {
	return m.CreateFn(book)
}
func (m *MockBookRepo) GetAll(genre string) ([]model.Book, error) {
	return m.GetAllFn(genre)
}
func (m *MockBookRepo) GetByID(id string) (*model.Book, error) {
	return m.GetByIDFn(id)
}
func (m *MockBookRepo) Update(book *model.Book) error {
	return m.UpdateFn(book)
}
func (m *MockBookRepo) Delete(id string) error {
	return m.DeleteFn(id)
}
func (m *MockBookRepo) UpdateStock(tx *sql.Tx, bookID string, quantity int) error {
	return m.UpdateStockFn(tx, bookID, quantity)
}
