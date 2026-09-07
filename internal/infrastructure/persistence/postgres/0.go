package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maadiii/taskmanager/config"
)

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.PgDb.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed on pg pool creation: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed on pg pool ping: %w", err)
	}

	return pool, nil
}
