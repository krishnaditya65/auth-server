package domain

import (
	"context"
	"time"
)

type Credential struct {
	ID           string
	IdentityID   string
	CredentialID []byte
	PublicKey    []byte
	Attestation  string
	AAGUID       []byte
	SignCount    uint32
	Transports   []string
	UserHandle   []byte
	Label        string
	CreatedAt    time.Time
	LastUsedAt   *time.Time
}

type Repository interface {
	Create(ctx context.Context, cred *Credential) error
	GetByCredentialID(ctx context.Context, credID []byte) (*Credential, error)
	ListByIdentity(ctx context.Context, identityID string) ([]*Credential, error)
	ListByUserHandle(ctx context.Context, handle []byte) ([]*Credential, error)
	UpdateSignCount(ctx context.Context, id string, count uint32) error
	Delete(ctx context.Context, id string) error
}
