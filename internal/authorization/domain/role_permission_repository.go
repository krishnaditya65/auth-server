package domain

import "context"

type RolePermissionRepository interface {
	AssignPermission(
		ctx context.Context,
		rolePermission *RolePermission,
	) error

	GetPermissionsForRole(
		ctx context.Context,
		roleID string,
	) ([]string, error)

	GetPermissionsForUser(
		ctx context.Context,
		userID string,
	) ([]string, error)
}
