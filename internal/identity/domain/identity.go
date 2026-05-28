package domain

import "time"

type Identity struct {
	ID            string
	PrimaryEmail  string
	EmailVerified bool
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
