package domain

import "time"

type RefreshToken struct {
	ID string

	SessionID string

	TokenHash string

	ExpiresAt time.Time

	RevokedAt *time.Time

	CreatedAt time.Time
}
