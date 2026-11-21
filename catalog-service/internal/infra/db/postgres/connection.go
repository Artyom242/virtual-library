package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection opens a concurrency-safe connection pool to Postgres.
// A *pgxpool.Pool (unlike a single *pgx.Conn) is safe for use by multiple
// goroutines, which is required when serving concurrent HTTP requests.
func NewConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// Fail fast if the database is unreachable.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
