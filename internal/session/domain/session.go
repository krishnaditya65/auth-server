package domain

import "time"

type Session struct {
	ID               string
	TenantID         string
	IdentityID       string
	UserID           string
	RefreshTokenHash string
	ParentSessionID  *string
	IPAddress        *string
	UserAgent        string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
}
