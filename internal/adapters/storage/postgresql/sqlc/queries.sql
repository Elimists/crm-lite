-- name: CreateContact :one
INSERT INTO contacts (
    name, email, phone, message, source_domain, status
) VALUES (
    $1, $2, $3, $4, $5, 'new'
)
RETURNING *;
