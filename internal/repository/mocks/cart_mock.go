package mocks

import "bookstore/internal/model"

type MockCartRepo struct {
	AddItemFn          func(item *model.CartItem) error
	GetCartByUserIDFn  func(userID string) ([]model.CartItem, error)
	UpdateQuantityFn   func(id string, userID string, quantity int) error
	RemoveItemFn       func(id string, userID string) error
	GetItemByIDFn      func(id string, userID string) (*model.CartItem, error)
}

func (m *MockCartRepo) AddItem(item *model.CartItem) error {
	return m.AddItemFn(item)
}
func (m *MockCartRepo) GetCartByUserID(userID string) ([]model.CartItem, error) {
	return m.GetCartByUserIDFn(userID)
}
func (m *MockCartRepo) UpdateQuantity(id string, userID string, quantity int) error {
	return m.UpdateQuantityFn(id, userID, quantity)
}
func (m *MockCartRepo) RemoveItem(id string, userID string) error {
	return m.RemoveItemFn(id, userID)
}
func (m *MockCartRepo) GetItemByID(id string, userID string) (*model.CartItem, error) {
	return m.GetItemByIDFn(id, userID)
}
