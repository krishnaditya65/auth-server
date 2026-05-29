package domain

import "time"

type Tenant struct {
	ID          string
	Slug        string
	DisplayName string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
