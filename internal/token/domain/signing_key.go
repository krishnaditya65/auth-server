package domain

import (
	"context"
	"time"
)

type SigningKey struct {
	ID         string
	Algorithm  string
	PublicKey  string
	PrivateKey string
	IsActive   bool
	CreatedAt  time.Time
	RotatedAt  *time.Time
}

type Repository interface {
	GetActive(ctx context.Context) (*SigningKey, error)
	GetByID(ctx context.Context, id string) (*SigningKey, error)
	Create(ctx context.Context, key *SigningKey) error
	DeactivateAll(ctx context.Context) error
	ListPublic(ctx context.Context) ([]*SigningKey, error)
}
