package models

import (
	"time"
)

type Tenant struct {
	ID          int       `json:"-"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	CompanyUrl  string    `json:"company_url"`
	Description string    `json:"description"`
	Timezone    string    `json:"timezone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Tenant) IsSystem() bool {
	return t.ID == 1 || t.Slug == "proreact"
}

func (t *Tenant) Location() *time.Location {
	loc, err := time.LoadLocation(t.Timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}

type TenantConfig struct {
	ID          int       `json:"-"`
	TenantId    int       `json:"tenant_id"`
	DisplayName string    `json:"display_name"`
	ConfigName  string    `json:"config_name"`
	Value       string    `json:"value"`
	ValueType   string    `json:"value_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TenantPatch struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	CompanyUrl  *string `json:"company_url"`
	Description *string `json:"description"`
	Timezone    *string `json:"timezone"`
}

type TenantConfigPatch struct {
	DisplayName *string `json:"display_name"`
	ConfigName  *string `json:"config_name"`
	Value       *string `json:"value"`
	ValueType   *string `json:"value_type"`
	Description *string `json:"description"`
}
