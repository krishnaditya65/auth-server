package tx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	db *pgxpool.Pool
}

func NewManager(db *pgxpool.Pool) *Manager {
	return &Manager{
		db: db,
	}
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

	txCtx := context.WithValue(ctx, "tx", tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
