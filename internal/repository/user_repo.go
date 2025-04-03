package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"github.com/qtj4/user-service/internal/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Username, user.Email, user.PasswordHash, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, "UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4",
		user.Username, user.Email, now, user.ID)
	if err != nil {
		return err
	}
	user.UpdatedAt = now
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}