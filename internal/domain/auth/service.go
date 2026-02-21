package auth

import (
	"context"
	repo "crm-lite/internal/adapters/storage/postgresql/sqlc"
	"crm-lite/internal/jwt"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	// Register(ctx context.Context, params RegisterParams) (repo.User, error)
}

type svc struct {
	repo *repo.Queries
	auth *jwt.Authenticator
}

func NewService(db *pgxpool.Pool, auth *jwt.Authenticator) Service {
	return &svc{
		repo: repo.New(db),
		auth: auth,
	}
}

func (s *svc) Authenticate(ctx context.Context, username string, password string) (string, error) {

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.auth.CreateToken(user.UserName, user.TenantSlug, user.Roles, user.Scopes)
	if err != nil {
		return "", errors.New("failed to create token")
	}

	return token, nil
}
