package contacts

import (
	"context"
	repo "crm-lite/internal/adapters/storage/postgresql/sqlc"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type svc struct {
	db   *pgxpool.Pool
	repo *repo.Queries
}

func NewService(db *pgxpool.Pool) Service {
	return &svc{
		db:   db,
		repo: repo.New(db),
	}
}

func (s *svc) GetContact(ctx context.Context, id int32) (Contact, error) {

	contact, err := s.repo.GetContact(ctx, id)
	if err != nil {
		return Contact{}, err
	}
	return Contact{
		ID:           contact.ID,
		Name:         contact.Name,
		Email:        contact.Email,
		Phone:        contact.Phone.String,
		Message:      contact.Message,
		SourceDomain: contact.SourceDomain,
		CreatedAt:    contact.CreatedAt.Time,
	}, nil
}

func (s *svc) CreateContact(ctx context.Context, c CreateContactParams) (Contact, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Contact{}, errors.New("error acquiring db connection from pool")
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	contact, err := qtx.CreateContact(ctx, repo.CreateContactParams{
		Name:         c.Name,
		Email:        c.Email,
		Phone:        pgtype.Text{String: c.Phone, Valid: c.Phone != ""},
		Message:      c.Message,
		SourceDomain: c.SourceDomain,
	})
	if err != nil {
		return Contact{}, errors.New("error creating contact")
	}

	if err := tx.Commit(ctx); err != nil {
		return Contact{}, err
	}

	return Contact{
		ID:           contact.ID,
		Name:         contact.Name,
		Email:        contact.Email,
		Phone:        contact.Phone.String,
		Message:      contact.Message,
		SourceDomain: contact.SourceDomain,
		CreatedAt:    contact.CreatedAt.Time,
	}, nil
}
