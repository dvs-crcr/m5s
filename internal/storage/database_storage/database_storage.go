package databaseStorage

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

func (ids *DBStorage) Ping(ctx context.Context) error {
    return ids.pool.Ping(ctx)
}

func (ids *DBStorage) Update(ctx context.Context, metric *domain.Metric) error {
    //TODO implement me
    panic("implement me 1")
}

func (ids *DBStorage) GetMetric(ctx context.Context, metricType domain.MetricType, name string) (*domain.Metric, error) {
    //TODO implement me
    panic("implement me 2")
}

func (ids *DBStorage) GetMetricsList(ctx context.Context) ([]*domain.Metric, error) {
    fmt.Println("CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")

    metrics := make([]*domain.Metric, 0)

    rows, err := ids.pool.Query(
        ctx,
        `SELECT
            id,
            metric_type,
            delta,
            value
        FROM metrics.metrics;`)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var row *domain.Metric

        if err := rows.Scan(
            &row.Name, &row.Type, &row.FloatValue, &row.IntValue,
        ); err != nil {
            return nil, err
        }

        metrics = append(metrics, row)
    }

    fmt.Println("AAAAAAAAAAAAAAAAAAAA", metrics)

    return metrics, err
}

func (ids *DBStorage) UpdateMetrics(ctx context.Context, metrics []*domain.Metric) error {
    //TODO implement me
    panic("implement me 4")
}