package contacts

import (
	"context"
	"time"
)

// Contact model
type Contact struct {
	ID           int32
	Name         string
	Email        string
	Phone        string
	Message      string
	SourceDomain string
	CreatedAt    time.Time
}

// CreateContactParams defines the parameters required to create a new contact
type CreateContactParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Message      string `json:"message"`
	SourceDomain string `json:"source_domain"`
}

// Service defines the interface for contact-related operations
type Service interface {
	CreateContact(ctx context.Context, tempContact CreateContactParams) (Contact, error)
	GetContact(ctx context.Context, id int32) (Contact, error)
}
