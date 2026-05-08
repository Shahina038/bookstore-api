package mocks

import "bookstore/internal/model"

type MockOrderRepo struct {
	CreateOrderFn       func(userID string) (*model.Order, error)
	GetOrdersByUserIDFn func(userID string) ([]model.Order, error)
	GetOrderByIDFn      func(id string, userID string) (*model.Order, error)
}

func (m *MockOrderRepo) CreateOrder(userID string) (*model.Order, error) {
	return m.CreateOrderFn(userID)
}
func (m *MockOrderRepo) GetOrdersByUserID(userID string) ([]model.Order, error) {
	return m.GetOrdersByUserIDFn(userID)
}
func (m *MockOrderRepo) GetOrderByID(id string, userID string) (*model.Order, error) {
	return m.GetOrderByIDFn(id, userID)
}
