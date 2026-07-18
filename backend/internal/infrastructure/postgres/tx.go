package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBTX is the common interface satisfied by both *pgxpool.Pool and pgx.Tx.
// Repository methods accept this so they can run either directly on the pool
// or within an explicit transaction provided via the context.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// txKey is the context key used to store an active transaction.
type txKey struct{}

// WithTx stores a pgx.Tx in the context so that repository methods can
// participate in the caller's transaction.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromContext retrieves a transaction from the context, if one is present.
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// withTx runs fn inside a transaction. If the context already carries a
// transaction, fn is executed with it and the caller retains commit/rollback
// responsibility. Otherwise a new transaction is begun, committed on success,
// and rolled back on error.
func withTx(ctx context.Context, pool *pgxpool.Pool, fn func(ctx context.Context) error) error {
	if _, ok := TxFromContext(ctx); ok {
		// Already inside a transaction — delegate lifecycle to the caller.
		return fn(ctx)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(WithTx(ctx, tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

// TxManager implements the application-layer TxManager interface using a
// pgxpool.Pool. It begins a transaction, stores it in the context so that
// repository calls participate, and commits or rolls back based on the result
// of the supplied function.
//
// TODO(wire): wire TxManager into wire.go's NewServices so the ticket service
// receives it via SetTxManager.
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager creates a new TxManager backed by the given pool.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// RunInTx executes fn within a database transaction. The transaction is
// committed if fn returns nil, and rolled back otherwise.
func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return withTx(ctx, m.pool, fn)
}
