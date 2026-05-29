package postgres

import (
	"context"

	authdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(
	db *pgxpool.Pool,
) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(
	ctx context.Context,
	role *authdomain.Role,
) error {

	query := `
INSERT INTO roles (
    id,
    tenant_id,
    name,
    description,
    created_at
)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		role.ID,
		role.TenantID,
		role.Name,
		role.Description,
		role.CreatedAt,
	)

	return err
}

func (r *RoleRepository) GetByID(
	ctx context.Context,
	id string,
) (*authdomain.Role, error) {

	query := `
		SELECT
			id,
			tenant_id,
			name,
			description,
			created_at,
			updated_at
		FROM roles
		WHERE id = $1
	`

	row := r.executor(ctx).QueryRow(
		ctx,
		query,
		id,
	)

	var role authdomain.Role

	err := row.Scan(
		&role.ID,
		&role.TenantID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetByTenantAndName(
	ctx context.Context,
	tenantID string,
	name string,
) (*authdomain.Role, error) {

	query := `
		SELECT
			id,
			tenant_id,
			name,
			description,
			created_at,
		FROM roles
		WHERE tenant_id = $1
		  AND name = $2
	`

	row := r.executor(ctx).QueryRow(
		ctx,
		query,
		tenantID,
		name,
	)

	var role authdomain.Role

	err := row.Scan(
		&role.ID,
		&role.TenantID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

type roleExecutor interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)

	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
}

func (r *RoleRepository) executor(
	ctx context.Context,
) roleExecutor {

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}
