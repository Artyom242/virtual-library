package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTransactor struct {
	pool *pgxpool.Pool
}

func NewPgxTransactor(pool *pgxpool.Pool) *PgxTransactor {
	return &PgxTransactor{pool}
}

func (t *PgxTransactor) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.BeginFunc(ctx, t.pool, func(tx pgx.Tx) error {
		return fn(context.WithValue(ctx, "txKey{}", tx))
	})
}
