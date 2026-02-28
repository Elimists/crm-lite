-- name: CreateContact :one
INSERT INTO contacts (
    name, email, phone, message, source_domain, status
) VALUES (
    $1, $2, $3, $4, $5, 'new'
)
RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1;

-- name: GetUser :one
SELECT 
    u.*, 
    t.slug AS tenant_slug 
FROM users u
JOIN tenants t ON u.tenant_id = t.id
WHERE u.email = $1 LIMIT 1;