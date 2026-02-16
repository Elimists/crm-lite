-- name: GetTenantBySlug :one
SELECT * FROM tenants
WHERE slug = $1 LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants;

-- name: GetConfigsByTenantSlug :many
SELECT tc.* FROM tenant_configs tc
JOIN tenants t ON tc.tenant_id = t.id
WHERE t.slug = $1
ORDER BY tc.config_name;

-- name: GetUserForAuth :one
SELECT id, user_name, password_hash, tenant_id, roles, scopes
FROM users
where email = $1 LIMIT 1;

-- name: GetUserWithTenantByUserName :one
SELECT 
    u.id, u.user_name, u.fname, u.lname, u.email, u.roles, u.scopes,
    t.slug AS tenant_slug,
    t.name AS tenant_name
FROM users u
JOIN tenants t ON u.tenant_id = t.id
WHERE user_name = $1 LIMIT 1;

-- name: ListUsersByTenantSlug :many
SELECT
    u.id, u.user_name, u.fname, u.lname, u.email, u.roles, u.scopes, u.created_at
FROM users u
JOIN tenants t ON u.tenant_id = t.id
WHERE t.slug = $1
ORDER BY u.created_at DESC;

-- name: ListContactsByTenantSlug :many
SELECT
    c.name, c.email, c.phone, c.message, c.status, c.created_at, c.updated_at
FROM contacts c
JOIN tenants t ON c.tenant_id = t.id
WHERE t.slug = $1
ORDER BY c.created_at DESC;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: InsertContact :one
INSERT INTO contacts (
    name, email, phone, message, source_domain, 
    tenant_id, 
    status
) VALUES (
    $1, $2, $3, $4, $5, 
    (SELECT id FROM tenants WHERE slug = $6), 
    'new'
)
RETURNING *;

-- name: UpdateContactStatus :one
UPDATE contacts
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;