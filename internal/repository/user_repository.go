package repository

import (
	"database/sql"
	"bookstore/internal/model"
)

type UserRepo interface {
    Create(user *model.User) error
    GetByEmail(email string) (*model.User, error)
    GetByID(id string) (*model.User, error)
    SoftDelete(id string) error
	GetAll() ([]*model.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, password, is_admin, is_deleted, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_deleted = FALSE`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin,
		&user.IsDeleted, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, is_admin, is_deleted, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.IsAdmin,
		&user.IsDeleted, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (r *UserRepository) SoftDelete(id string) error {
    query := `
        UPDATE users
        SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1`

    _, err := r.db.Exec(query, id)
    return err
}

func (r *UserRepository) GetAll() ([]*model.User, error) {
	query := `
		SELECT id, name, email, is_admin, is_deleted, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.IsDeleted,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
