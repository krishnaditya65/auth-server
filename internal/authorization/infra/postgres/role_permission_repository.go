package postgres

import (
	"context"

	authdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RolePermissionRepository struct {
	db *pgxpool.Pool
}

func NewRolePermissionRepository(
	db *pgxpool.Pool,
) *RolePermissionRepository {
	return &RolePermissionRepository{
		db: db,
	}
}

func (r *RolePermissionRepository) AssignPermission(
	ctx context.Context,
	rolePermission *authdomain.RolePermission,
) error {

	query := `
		INSERT INTO role_permissions (
			role_id,
			permission_id
		)
		VALUES ($1,$2)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		rolePermission.RoleID,
		rolePermission.PermissionID,
	)

	return err
}

func (r *RolePermissionRepository) GetPermissionsForRole(
	ctx context.Context,
	roleID string,
) ([]string, error) {

	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p
			ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.executor(ctx).Query(
		ctx,
		query,
		roleID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions []string

	for rows.Next() {

		var permission string

		err := rows.Scan(
			&permission,
		)

		if err != nil {
			return nil, err
		}

		permissions = append(
			permissions,
			permission,
		)
	}

	return permissions, nil
}

func (r *RolePermissionRepository) GetPermissionsForUser(
	ctx context.Context,
	userID string,
) ([]string, error) {

	query := `
		SELECT DISTINCT p.name
		FROM user_roles ur

		JOIN role_permissions rp
			ON rp.role_id = ur.role_id

		JOIN permissions p
			ON p.id = rp.permission_id

		WHERE ur.user_id = $1
	`

	rows, err := r.executor(ctx).Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions []string

	for rows.Next() {

		var permission string

		err := rows.Scan(
			&permission,
		)

		if err != nil {
			return nil, err
		}

		permissions = append(
			permissions,
			permission,
		)
	}

	return permissions, nil
}

type rolePermissionExecutor interface {
	Exec(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgconn.CommandTag, error)

	Query(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgx.Rows, error)
}

func (r *RolePermissionRepository) executor(
	ctx context.Context,
) rolePermissionExecutor {

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}
