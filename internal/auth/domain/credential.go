package domain

import "time"

type Credential struct {
	ID             string
	IdentityID     string
	CredentialType string
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
