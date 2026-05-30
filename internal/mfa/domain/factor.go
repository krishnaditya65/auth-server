package domain

import (
	"context"
	"time"
)

type Factor struct {
	ID              string
	IdentityID      string
	FactorType      string
	SecretEncrypted string
	Label           string
	Verified        bool
	CreatedAt       time.Time
}

type Repository interface {
	Create(ctx context.Context, factor *Factor) error
	GetByID(ctx context.Context, id string) (*Factor, error)
	GetByIdentity(ctx context.Context, identityID string) ([]*Factor, error)
	GetVerifiedByIdentity(ctx context.Context, identityID string) ([]*Factor, error)
	MarkVerified(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
