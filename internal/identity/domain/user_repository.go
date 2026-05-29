package domain

import "context"

type UserRepository interface {
	Create(
		ctx context.Context,
		user *User,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (*User, error)

	GetByTenantAndIdentity(
		ctx context.Context,
		tenantID string,
		identityID string,
	) (*User, error)

	GetByIdentityID(
		ctx context.Context,
		identityID string,
	) (*User, error)
}
