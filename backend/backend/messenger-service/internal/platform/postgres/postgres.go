package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Ping failed: %w", err)
	}

	return pool, nil
}