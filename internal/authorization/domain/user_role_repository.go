package domain

import "context"

type UserRoleRepository interface {
	AssignRole(
		ctx context.Context,
		userRole *UserRole,
	) error

	GetRolesForUser(
		ctx context.Context,
		userID string,
	) ([]string, error)
}
