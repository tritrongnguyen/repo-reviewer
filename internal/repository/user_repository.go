package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserDetailsByEmail(ctx context.Context, email string) (*domain.UserDetails, error)
	UpdatePassword(ctx context.Context, userID, hash string) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	sql := `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, sql,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	return err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sql := `
		SELECT id, email, password_hash, role, created_at
		FROM users WHERE email=$1
	`
	row := r.db.QueryRow(ctx, sql, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetUserDetailsByEmail(ctx context.Context, email string) (*domain.UserDetails, error) {
	return nil, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID, hash string) error {
	return nil
}
