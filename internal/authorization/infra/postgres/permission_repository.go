package postgres

import (
	"context"

	authdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionRepository struct {
	db *pgxpool.Pool
}

func NewPermissionRepository(
	db *pgxpool.Pool,
) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) Create(
	ctx context.Context,
	permission *authdomain.Permission,
) error {

	query := `
		INSERT INTO permissions (
			id,
			name,
			description,
			created_at
		)
		VALUES ($1,$2,$3,$4)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.CreatedAt,
	)

	return err
}

func (r *PermissionRepository) GetByID(
	ctx context.Context,
	id string,
) (*authdomain.Permission, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at
		FROM permissions
		WHERE id = $1
	`

	row := r.executor(ctx).QueryRow(
		ctx,
		query,
		id,
	)

	var permission authdomain.Permission

	err := row.Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetByName(
	ctx context.Context,
	name string,
) (*authdomain.Permission, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at
		FROM permissions
		WHERE name = $1
	`

	row := r.executor(ctx).QueryRow(
		ctx,
		query,
		name,
	)

	var permission authdomain.Permission

	err := row.Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

type permissionExecutor interface {
	Exec(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgconn.CommandTag, error)

	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
}

func (r *PermissionRepository) executor(
	ctx context.Context,
) permissionExecutor {

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}
