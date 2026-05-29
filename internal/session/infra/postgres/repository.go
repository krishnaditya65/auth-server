package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	sessiondomain "github.com/krishnaditya65/auth-server/internal/session/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	session *sessiondomain.Session,
) error {
	query := `
		INSERT INTO sessions (
			id,
			tenant_id,
			identity_id,
			user_id,
			refresh_token_hash,
			parent_session_id,
			ip_address,
			user_agent,
			expires_at,
			revoked_at,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		session.ID,
		session.TenantID,
		session.IdentityID,
		session.UserID,
		session.RefreshTokenHash,
		session.ParentSessionID,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
		session.RevokedAt,
		session.CreatedAt,
	)

	return err
}

type executor interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
}

func (r *Repository) executor(ctx context.Context) executor {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}
