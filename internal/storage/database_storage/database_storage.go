package databasestorage

import (
    "context"
    "fmt"
    "io/fs"
    "os"
    "strconv"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/jackc/tern/v2/migrate"

    "m5s/domain"
    internalLogger "m5s/pkg/logger"
)

var logger = internalLogger.NewLogger()

type DBStorage struct {
    pool           *pgxpool.Pool
    migrationsPath string
}

var SchemaVersionTable = "schema_version"

func NewDBStorage(
    ctx context.Context,
    dsn string,
    migrationsPath string,
    migrationsSchemaVersion string,
) (*DBStorage, error) {
    logger = logger.With(
        "package", "storage",
        "type", "db",
    )
    logger.Info("init storage")

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
        parseMigrationVersion(migrationsSchemaVersion),
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

func parseMigrationVersion(migrationsVersion string) int32 {
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

func (ids *DBStorage) GetMetric(ctx context.Context, metricType domain.MetricType, name string) (*domain.Metric, error) {
    metric := &domain.Metric{}

    query := `
        SELECT
            id,
            metric_type,
            delta,
            value
        FROM metrics.metrics
        WHERE
            id=$1
            AND metric_type=$2;
    `

    row := ids.pool.QueryRow(ctx, query, name, metricType.String())

    if err := row.Scan(
        &metric.Name, &metric.Type, &metric.IntValue, &metric.FloatValue,
    ); err != nil {
        logger.Errorw(err.Error())
        return nil, err
    }

    return metric, nil
}

func (ids *DBStorage) GetMetricsList(ctx context.Context) ([]*domain.Metric, error) {
    metrics := make([]*domain.Metric, 0)

    query := `
        SELECT
            id,
            metric_type,
            delta,
            value
        FROM metrics.metrics;
    `

    rows, err := ids.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var row domain.Metric

        if err := rows.Scan(
            &row.Name, &row.Type, &row.IntValue, &row.FloatValue,
        ); err != nil {
            return nil, err
        }

        metrics = append(metrics, &row)
    }

    return metrics, err
}

func (ids *DBStorage) Update(ctx context.Context, metric *domain.Metric) error {
    query := `
        INSERT INTO metrics.metrics(id, metric_type, delta, value)
        VALUES($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            delta = metrics.delta + $3,
            value = $4;
    `

    if _, err := ids.pool.Exec(
        ctx,
        query,
        metric.Name, metric.Type.String(), metric.IntValue, metric.FloatValue,
    ); err != nil {
        return err
    }

    return nil
}

func (ids *DBStorage) UpdateMetrics(ctx context.Context, metrics []*domain.Metric) error {
    logger.Debugw("update metrics", "metrics", fmt.Sprintf("%v", metrics))

    batch := &pgx.Batch{}

    query := `
        INSERT INTO metrics.metrics(id, metric_type, delta, value)
        VALUES($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            delta = metrics.delta + $3,
            value = $4;
    `

    for _, metric := range metrics {
        batch.Queue(query, metric.Name, metric.Type.String(), metric.IntValue, metric.FloatValue)
    }

    br := ids.pool.SendBatch(ctx, batch)

    if _, err := br.Exec(); err != nil {
        return err
    }

    return nil
}
