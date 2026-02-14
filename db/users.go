package db

import (
	"context"
	"crm-lite/models"
	"fmt"
	"strings"
)

type UserQueries struct {
	db *Database
}

func (q *UserQueries) findUser(ctx context.Context, column string, value interface{}) (*models.User, error) {
	query := fmt.Sprintf(`
        SELECT id, user_name, fname, lname, email, password_hash, tenant_id, roles, scopes 
        FROM users 
        WHERE %s = $1`, column)

	var u models.User
	err := q.db.QueryRow(ctx, query, value).Scan(
		&u.ID, &u.UserName, &u.Fname, &u.Lname,
		&u.Email, &u.PasswordHash, &u.TenantId, &u.Roles, &u.Scopes,
	)

	return &u, err
}

func (q *UserQueries) QueryUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return q.findUser(ctx, "email", strings.ToLower(email))
}

func (q *UserQueries) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return q.findUser(ctx, "id", id)
}

func (q *UserQueries) GetUserByUserName(ctx context.Context, user_name string) (*models.User, error) {
	return q.findUser(ctx, "user_name", user_name)
}

func (q *UserQueries) InsertUser(ctx context.Context, u models.User) (int, error) {
	query := `
        INSERT INTO users (
            user_name, fname, lname, email, password_hash, 
            tenant_id, roles, scopes, created_at, updated_at
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id`

	var newID int

	err := q.db.QueryRow(
		ctx,
		query,
		u.UserName,
		u.Fname,
		u.Lname,
		u.Email,
		u.PasswordHash,
		u.TenantId,
		u.Roles,
		u.Scopes,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (q *UserQueries) PatchUser(ctx context.Context, id int, patch models.UserPatch) error {
	query := "UPDATE users SET "
	var args []interface{}
	var updates []string
	argCounter := 1

	if patch.UserName != nil && *patch.UserName != "" {
		updates = append(updates, fmt.Sprintf("user_name = $%d", argCounter))
		args = append(args, *patch.UserName)
		argCounter++
	}
	if patch.Fname != nil {
		updates = append(updates, fmt.Sprintf("fname = $%d", argCounter))
		args = append(args, *patch.Fname)
		argCounter++
	}
	if patch.Lname != nil {
		updates = append(updates, fmt.Sprintf("lname = $%d", argCounter))
		args = append(args, *patch.Lname)
		argCounter++
	}
	if patch.Email != nil && *patch.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argCounter))
		args = append(args, *patch.Email)
		argCounter++
	}
	if patch.Password != nil && *patch.Password != "" {
		updates = append(updates, fmt.Sprintf("password_hash = $%d", argCounter))
		args = append(args, *patch.Password)
		argCounter++
	}
	if patch.TenantId != nil && *patch.TenantId != 0 {
		updates = append(updates, fmt.Sprintf("tenant_id = $%d", argCounter))
		args = append(args, *patch.TenantId)
		argCounter++
	}
	if patch.Roles != nil {
		updates = append(updates, fmt.Sprintf("roles = $%d", argCounter))
		args = append(args, *patch.Roles)
		argCounter++
	}
	if patch.Scopes != nil {
		updates = append(updates, fmt.Sprintf("scopes = $%d", argCounter))
		args = append(args, *patch.Scopes)
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

func (q *UserQueries) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}
