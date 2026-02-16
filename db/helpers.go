package db

import "time"

// USER specific helpers
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

// TENANT specific helpers
func (t *Tenant) IsSystem() bool {
	return t.ID == 1 || t.Slug == "proreact"
}

func (t *Tenant) Location() *time.Location {
	loc, err := time.LoadLocation(t.Timezone.String)
	if err != nil {
		return time.UTC
	}
	return loc
}

// CONTACT specific helpers
func (c *Contact) IsNew() bool {
	return c.Status == "new"
}

func (c *Contact) IsValid() bool {
	return c.Name != "" && c.Email != "" && c.Message != ""
}

func (c *Contact) IsFromDomain(domain string) bool {
	return c.SourceDomain == domain
}
