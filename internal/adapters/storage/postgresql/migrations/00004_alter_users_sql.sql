-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ADD COLUMN tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE;
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_tenant_id;
ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
-- +goose StatementEnd
