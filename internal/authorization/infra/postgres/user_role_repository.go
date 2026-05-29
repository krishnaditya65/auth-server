package postgres

import (
	"context"

	authdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRoleRepository struct {
	db *pgxpool.Pool
}

func NewUserRoleRepository(
	db *pgxpool.Pool,
) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}

func (r *UserRoleRepository) AssignRole(
	ctx context.Context,
	userRole *authdomain.UserRole,
) error {

	query := `
		INSERT INTO user_roles (
			user_id,
			role_id
		)
		VALUES ($1,$2)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		userRole.UserID,
		userRole.RoleID,
	)

	return err
}

type userRoleExecutor interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)

	Query(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgx.Rows, error)
}

func (r *UserRoleRepository) executor(
	ctx context.Context,
) userRoleExecutor {

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}

func (r *UserRoleRepository) GetRolesForUser(
	ctx context.Context,
	userID string,
) ([]string, error) {

	query := `
		SELECT r.name
		FROM user_roles ur
		JOIN roles r
			ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roles []string

	for rows.Next() {
		var role string

		err := rows.Scan(
			&role,
		)

		if err != nil {
			return nil, err
		}

		roles = append(
			roles,
			role,
		)
	}

	return roles, nil
}
