package domain

import "context"

type Repository interface {
	Create(ctx context.Context, credential *Credential) error
	GetByIdentityID(
		ctx context.Context,
		identityID string,
	) (*Credential, error)
}
