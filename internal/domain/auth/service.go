package auth

import (
	"context"
	repo "crm-lite/internal/adapters/storage/postgresql/sqlc"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Authenticate(ctx context.Context, email, password string) (repo.GetUserRow, error)
	// Register(ctx context.Context, params RegisterParams) (repo.User, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(db repo.DBTX) Service {
	return &svc{
		repo: repo.New(db),
	}
}

func (s *svc) Authenticate(ctx context.Context, username string, password string) (repo.GetUserRow, error) {

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return repo.GetUserRow{}, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return repo.GetUserRow{}, errors.New("invalid credentials")
	}

	return user, nil
}
