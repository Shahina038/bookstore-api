package mocks

import "bookstore/internal/model"

type MockUserRepo struct {
	CreateFn     func(user *model.User) error
	GetByEmailFn func(email string) (*model.User, error)
	GetByIDFn    func(id string) (*model.User, error)
	SoftDeleteFn func(id string) error
}

func (m *MockUserRepo) Create(user *model.User) error {
	return m.CreateFn(user)
}

func (m *MockUserRepo) GetByEmail(email string) (*model.User, error) {
	return m.GetByEmailFn(email)
}

func (m *MockUserRepo) GetByID(id string) (*model.User, error) {
	return m.GetByIDFn(id)
}

func (m *MockUserRepo) SoftDelete(id string) error {
	return m.SoftDeleteFn(id)
}
