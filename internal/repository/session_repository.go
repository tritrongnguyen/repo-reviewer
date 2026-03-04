package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	Create(ctx context.Context, sessionID, userID string, expires time.Time) error
	Delete(ctx context.Context, sessionID string) error
}

type sessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &sessionRepo{
		db: db,
	}
}

func (r *sessionRepo) Create(ctx context.Context, sessionID, userID string, expires time.Time) error {
	_, err := r.db.Exec(ctx,
		`
			INSERT INTO sessions (id, user_id, expires_at)
			VALUES ($1, $2, $3)
		`,
		sessionID, userID, expires,
	)
	return err
}

func (r *sessionRepo) Delete(ctx context.Context, sessionID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM sessions WHERE id=$1`,
		sessionID)
	return err
}
