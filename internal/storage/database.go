package storage

import (
    "context"
    "fmt"
    "io/fs"
    "os"
    "strconv"

    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/jackc/tern/v2/migrate"

    "m5s/domain"
)

type DBStorage struct {
    pool           *pgxpool.Pool
    storageType    StorageType
    migrationsPath string
}

var SchemaVersionTable = "schema_version"

func NewDBStorage(
    ctx context.Context,
    dsn string,
    migrationsPath string,
    migrationsSchemaVersion int32,
) (*DBStorage, error) {
    dbStorage := &DBStorage{
        storageType:    TypeDatabase,
        migrationsPath: migrationsPath,
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

    if err := dbStorage.startMigrations(
        ctx,
        SchemaVersionTable,
        migrationsSchemaVersion,
    ); err != nil {
        return nil, err
    }

    return dbStorage, nil
}

func (ids *DBStorage) startMigrations(
    ctx context.Context,
    versionTable string,
    schemaVersion int32,
) error {
    conn, err := ids.pool.Acquire(ctx)
    if err != nil {
        return fmt.Errorf(
            "unable to acquire conn from pool: %w", err,
        )
    }

    migrator, err := migrate.NewMigrator(ctx, conn.Conn(), versionTable)
    if err != nil {
        return err
    }

    fSys, err := fs.Sub(os.DirFS("."), ids.migrationsPath)
    if err != nil {
        return err
    }

    if err := migrator.LoadMigrations(fSys); err != nil {
        return err
    }

    if schemaVersion == 0 {
        schemaVersion = int32(len(migrator.Migrations))
    }

    if err := migrator.MigrateTo(ctx, schemaVersion); err != nil {
        return err
    }

    return nil
}

func ParseMigrationVersion(migrationsVersion string) int32 {
    if migrationsVersion == "" {
        return 0
    }

    msv, err := strconv.ParseInt(migrationsVersion, 10, 32)
    if err != nil {
        return 0
    }

    return int32(msv)
}

func (ids *DBStorage) MyType() StorageType {
    return ids.storageType
}

func (ids *DBStorage) Ping(ctx context.Context) error {
    return ids.pool.Ping(ctx)
}

func (ids *DBStorage) Update(metric *domain.Metric) error {
    //TODO implement me
    panic("implement me 1")
}

func (ids *DBStorage) GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error) {
    //TODO implement me
    panic("implement me 2")
}

func (ids *DBStorage) GetMetricsList() ([]*domain.Metric, error) {
    //TODO implement me
    panic("implement me 3")
}

func (ids *DBStorage) UpdateMetrics(metrics []*domain.Metric) error {
    //TODO implement me
    panic("implement me 4")
}
