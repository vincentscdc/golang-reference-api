package sqlc

import (
	"context"
	"fmt"

	"golangreferenceapi/internal/api/configuration"
	"golangreferenceapi/internal/db"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monacohq/golang-common/database/pginit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newConnPool(ctx context.Context, cfg *configuration.Database) (*pgxpool.Pool, error) {
	pgi, err := pginit.New(
		&pginit.Config{
			Host:         cfg.Host,
			Port:         cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			Database:     cfg.Database,
			MaxConns:     cfg.MaxConns,
			MaxIdleConns: cfg.MaxIdleConns,
			MaxLifeTime:  cfg.MaxLifeTime,
		},
		pginit.WithLogLevel(zerolog.WarnLevel),
		pginit.WithLogger(&log.Logger, "request-id"),
		pginit.WithDecimalType(),
		pginit.WithUUIDType(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PGInit: %w", err)
	}

	pool, err := pgi.ConnPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func NewRepo(ctx context.Context, cfg *configuration.Database) (*Repo, error) {
	pool, err := newConnPool(ctx, cfg)
	if err != nil {
		return &Repo{}, fmt.Errorf("failed to init pgx: %w", err)
	}

	return NewSQLCRepository(db.New(pool)), nil
}
