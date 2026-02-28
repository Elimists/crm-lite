-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    company_url TEXT,
    description TEXT,
    timezone TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slug_lowercase CHECK (slug = LOWER(slug))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants;
-- +goose StatementEnd
