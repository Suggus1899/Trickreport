package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolConfig holds tunable parameters for a pgx connection pool.
type PoolConfig struct {
	MaxConns          int
	MinConns          int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration

	// StatementTimeout and IdleInTxTimeout bound how long a single connection
	// can be held by a runaway or lock-blocked query/transaction, so one bad
	// query can't starve the whole pool. Zero disables the guard.
	StatementTimeout time.Duration
	IdleInTxTimeout  time.Duration
}

// DefaultPoolConfig returns a sensible default pool configuration.
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxConns:          20,
		MinConns:          4,
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: 1 * time.Minute,
		StatementTimeout:  30 * time.Second,
		IdleInTxTimeout:   30 * time.Second,
	}
}

// NewPool creates and validates a pgx connection pool configured with the
// provided PoolConfig.
func NewPool(ctx context.Context, databaseURL string, config PoolConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	if config.MaxConns > 0 {
		cfg.MaxConns = int32(config.MaxConns)
	}
	if config.MinConns > 0 {
		cfg.MinConns = int32(config.MinConns)
	}
	if config.MaxConnLifetime > 0 {
		cfg.MaxConnLifetime = config.MaxConnLifetime
	}
	if config.MaxConnIdleTime > 0 {
		cfg.MaxConnIdleTime = config.MaxConnIdleTime
	}
	if config.HealthCheckPeriod > 0 {
		cfg.HealthCheckPeriod = config.HealthCheckPeriod
	}
	if config.StatementTimeout > 0 {
		cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprintf("%d", config.StatementTimeout.Milliseconds())
	}
	if config.IdleInTxTimeout > 0 {
		cfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = fmt.Sprintf("%d", config.IdleInTxTimeout.Milliseconds())
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
