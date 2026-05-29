package domain

import "time"

type User struct {
	ID          string
	TenantID    string
	IdentityID  string
	Username    string
	DisplayName string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
