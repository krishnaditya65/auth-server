package domain

import "context"

type PermissionRepository interface {
	Create(
		ctx context.Context,
		permission *Permission,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (*Permission, error)

	GetByName(
		ctx context.Context,
		name string,
	) (*Permission, error)
}
