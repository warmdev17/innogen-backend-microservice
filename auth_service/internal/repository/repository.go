package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"innogen-backend/shared/models"
)

// UserRepository handles database queries for the users table.
type UserRepository struct {
	pool *pgxpool.Pool
}

// New creates a new UserRepository with the given connection pool.
func New(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const selectUserColumns = `SELECT id, email, password, username, full_name, role, is_active, created_at, updated_at FROM users`

// FindByEmail looks up a user by their email address.
// Returns nil, nil if no user is found.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx, selectUserColumns+" WHERE email = $1", email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.FullName,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindByEmail: %w", err)
	}
	return u, nil
}

// FindByID looks up a user by their ID.
// Returns nil, nil if no user is found.
func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx, selectUserColumns+" WHERE id = $1", id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.FullName,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindByID: %w", err)
	}
	return u, nil
}
