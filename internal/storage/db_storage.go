package storage

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
    pool *pgxpool.Pool
}

func NewDBStorage(ctx context.Context, dsn string) (*DBStorage, error) {
    dbStorage := &DBStorage{}

    poolConfig, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }

    if dbStorage.pool, err = pgxpool.NewWithConfig(
        ctx,
        poolConfig,
    ); err != nil {
        return nil, err
    }

    return dbStorage, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
    return s.pool.Ping(ctx)
}
