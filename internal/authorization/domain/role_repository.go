package domain

import "context"

type RoleRepository interface {
	Create(
		ctx context.Context,
		role *Role,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (*Role, error)

	GetByTenantAndName(
		ctx context.Context,
		tenantID string,
		name string,
	) (*Role, error)
}
