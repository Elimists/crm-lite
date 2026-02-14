package db

import (
	"context"
	"crm-lite/models"
	"database/sql"
	"fmt"
	"strings"
)

type TenantQueries struct {
	db *Database
}

func (q *TenantQueries) QueryTenant(ctx context.Context, id int) (*models.Tenant, error) {
	query := `
        SELECT 
            id,
			name,
			slug,
			company_url, 
            description, 
            timezone,
			created_at,
			updated_at 
        FROM tenants 
        WHERE id = $1`

	row := q.db.QueryRow(ctx, query, id)
	var t models.Tenant
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.CompanyUrl,
		&t.Description,
		&t.Timezone,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (q *TenantQueries) QueryTenants(ctx context.Context) ([]*models.Tenant, error) {
	query := `
        SELECT 
            id,
            name,
            slug,
            company_url, 
            description, 
            timezone,
            created_at,
            updated_at 
        FROM tenants`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*models.Tenant

	for rows.Next() {
		var t models.Tenant
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Slug,
			&t.CompanyUrl,
			&t.Description,
			&t.Timezone,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

func (q *TenantQueries) QueryTenantConfigs(ctx context.Context, tenantId int) ([]models.TenantConfig, error) {

	query := `
        SELECT id, tenant_id, display_name, config_name, value, value_type, description, created_at, updated_at
        FROM tenant_configs
        WHERE tenant_id = $1`

	rows, err := q.db.Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Very important to prevent memory leaks!

	var configs []models.TenantConfig

	for rows.Next() {
		var tc models.TenantConfig
		err := rows.Scan(
			&tc.ID,
			&tc.TenantId,
			&tc.DisplayName,
			&tc.ConfigName,
			&tc.Value,
			&tc.ValueType,
			&tc.Description,
			&tc.CreatedAt,
			&tc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, tc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func (q *TenantQueries) InsertTenant(ctx context.Context, t models.Tenant) (int, error) {
	query := `
        INSERT INTO tenants (
            name, slug, company_url, description, timezone
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	var newID int
	err := q.db.QueryRow(
		ctx,
		query,
		t.Name, strings.ToLower(t.Slug), t.CompanyUrl,
		t.Description, t.Timezone,
		t.CreatedAt, t.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (q *TenantQueries) InsertTenantConfigs(ctx context.Context, configs []models.TenantConfig) error {

	if len(configs) == 0 {
		return nil
	}

	query := `
		INSERT INTO tenant_configs (
			tenant_id, display_name, config_name,
			value, value_type, description
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	tx, err := q.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, c := range configs {
		_, err := tx.Exec(ctx, query,
			c.TenantId, c.DisplayName, c.ConfigName,
			c.Value, c.ValueType, c.Description,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (q *TenantQueries) PatchTenant(ctx context.Context, id int, patch models.TenantPatch) error {

	query := "UPDATE tenants SET "
	var args []interface{}
	var updates []string
	argCounter := 1

	if patch.Name != nil && *patch.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *patch.Name)
		argCounter++
	}
	if patch.Slug != nil && *patch.Slug != "" {
		updates = append(updates, fmt.Sprintf("slug = $%d", argCounter))
		args = append(args, strings.ToLower(*patch.Slug))
		argCounter++
	}
	if patch.CompanyUrl != nil {
		updates = append(updates, fmt.Sprintf("company_url = $%d", argCounter))
		args = append(args, *patch.CompanyUrl) // Can be ""
		argCounter++
	}
	if patch.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *patch.Description) // Can be ""
		argCounter++
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	if len(args) == 0 {
		return nil // Nothing to update
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCounter)
	args = append(args, id)

	result, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant with id %d not found", id)
	}

	return nil
}

func (q *TenantQueries) PatchTenantConfig(ctx context.Context, id int, patch models.TenantConfigPatch) error {
	query := "UPDATE tenant_configs SET "
	var args []interface{}
	var updates []string
	argCounter := 1

	if patch.DisplayName != nil && *patch.DisplayName != "" {
		updates = append(updates, fmt.Sprintf("display_name = $%d", argCounter))
		args = append(args, *patch.DisplayName)
		argCounter++
	}
	if patch.ConfigName != nil && *patch.ConfigName != "" {
		updates = append(updates, fmt.Sprintf("config_name = $%d", argCounter))
		args = append(args, *patch.ConfigName)
		argCounter++
	}
	if patch.Value != nil && *patch.Value != "" {
		updates = append(updates, fmt.Sprintf("value = $%d", argCounter))
		args = append(args, *patch.Value)
		argCounter++
	}
	if patch.ValueType != nil && *patch.ValueType != "" {
		updates = append(updates, fmt.Sprintf("value = $%d", argCounter))
		args = append(args, *patch.ValueType)
		argCounter++
	}
	if patch.Description != nil {
		updates = append(updates, fmt.Sprintf("value = $%d", argCounter))
		args = append(args, *patch.Description)
		argCounter++
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	if len(args) == 0 {
		return nil // Nothing to update
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCounter)
	args = append(args, id)

	result, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant with id %d not found", id)
	}

	return nil
}

func (q *TenantQueries) DeleteTenant(ctx context.Context, id int) error {
	query := "DELETE FROM tenants WHERE id = $1"

	result, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant with ID %d not found", id)
	}

	return nil
}

func (q *TenantQueries) DeleteTenantConfig(ctx context.Context, id int) error {
	query := "DELETE FROM tenant_configs WHERE id = $1"

	result, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant config: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tenant config with ID %d not found", id)
	}

	return nil
}
