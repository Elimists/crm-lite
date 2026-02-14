package models

import "time"

type User struct {
	ID           int       `json:"-"`
	TenantId     int       `json:"tenant_id"`
	UserName     string    `json:"user_name"`
	Fname        string    `json:"fname"`
	Lname        string    `json:"lname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	Scopes       []string  `json:"scopes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) IsSystemAdmin() bool {
	for _, role := range u.Roles {
		if role == "sysadmin" || role == "superuser" {
			return true
		}
	}
	return false
}

func (u *User) HasScope(required string) bool {
	for _, s := range u.Scopes {
		if s == "*" || s == required {
			return true
		}
	}
	return false
}

type UserPatch struct {
	TenantId *int      `json:"tenant_id"`
	UserName *string   `json:"user_name"`
	Fname    *string   `json:"fname"`
	Lname    *string   `json:"lname"`
	Email    *string   `json:"email"`
	Password *string   `json:"password"`
	Roles    *[]string `json:"roles"`
	Scopes   *[]string `json:"scopes"`
}
