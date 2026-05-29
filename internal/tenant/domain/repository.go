package domain

import "context"

type Repository interface {
	Create(ctx context.Context, tenant *Tenant) error

	GetByID(
		ctx context.Context,
		id string,
	) (*Tenant, error)

	GetBySlug(
		ctx context.Context,
		slug string,
	) (*Tenant, error)

	ExistsBySlug(
		ctx context.Context,
		slug string,
	) (bool, error)
}
