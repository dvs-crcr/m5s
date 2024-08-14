package storage

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"

    "m5s/domain"
)

type DBStorage struct {
    pool        *pgxpool.Pool
    storageType StorageType
}

func NewDBStorage(ctx context.Context, dsn string) (*DBStorage, error) {
    dbStorage := &DBStorage{
        storageType: TypeDatabase,
    }

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

func (ids *DBStorage) MyType() StorageType {
    return ids.storageType
}

func (ids *DBStorage) Ping(ctx context.Context) error {
    return ids.pool.Ping(ctx)
}

func (ids *DBStorage) Update(metric *domain.Metric) error {
    //TODO implement me
    panic("implement me")
}

func (ids *DBStorage) GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error) {
    //TODO implement me
    panic("implement me")
}

func (ids *DBStorage) GetMetricsList() ([]*domain.Metric, error) {
    //TODO implement me
    panic("implement me")
}

func (ids *DBStorage) UpdateMetrics(metrics []*domain.Metric) error {
    //TODO implement me
    panic("implement me")
}
