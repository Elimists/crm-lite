package models

import "time"

type Contact struct {
	ID           int       `json:"-"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Message      string    `json:"message"`
	SourceDomain string    `json:"source_domain"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
