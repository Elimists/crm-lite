package contacts

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type createContactParams struct {
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Message      string      `json:"message"`
	SourceDomain string      `json:"source_domain"`
}
