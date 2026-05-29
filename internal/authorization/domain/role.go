package domain

import "time"

type Role struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	CreatedAt   time.Time
}
