package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/domain"
	"github.com/tritrongnguyen/repo-reviewer.git/internal/repository"
	"github.com/tritrongnguyen/repo-reviewer.git/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	ResetPassword(ctx context.Context, email, newPassword string) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(u repository.UserRepository, s repository.SessionRepository) AuthService {
	return &authService{
		userRepo:    u,
		sessionRepo: s,
	}
}

func (s *authService) SignUp(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         0,
	}

	err = s.userRepo.Create(ctx, user)

	if err != nil {
		if isDuplicateEmailExists(err) {
			return ErrEmailAlreadyExists
		}
	}

	return nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("Invalid credentials when get email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("Invalid credentials when compare password")
	}

	sessionID := uuid.NewString()
	exp := time.Now().Add(24 * time.Hour)

	err = s.sessionRepo.Create(ctx, sessionID, user.ID, exp)
	logger.Log.Debug("err", zap.Error(err))
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *authService) ResetPassword(ctx context.Context, email, newPassword string) error {
	return nil
}

func isDuplicateEmailExists(err error) bool {

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key"
	}

	return false
}
