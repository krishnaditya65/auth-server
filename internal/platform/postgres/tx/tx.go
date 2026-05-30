package tx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Executor is the shared interface satisfied by both *pgxpool.Pool and pgx.Tx.
// All postgres repositories use this instead of defining local executor interfaces.
type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type contextKey struct{}

type Manager struct {
	db *pgxpool.Pool
}

func NewManager(db *pgxpool.Pool) *Manager {
	return &Manager{db: db}
}

func (m *Manager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txCtx := context.WithValue(ctx, contextKey{}, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// FromContext returns the active pgx.Tx from ctx, if one was injected by WithinTransaction.
func FromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(contextKey{}).(pgx.Tx)
	return tx, ok
}
