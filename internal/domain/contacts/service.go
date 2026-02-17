package contacts

import "context"

type Service interface {
	CreateContact(ctx context.Context) error
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}

func (s *svc) CreateContact(ctx context.Context) error {
	return nil
}
