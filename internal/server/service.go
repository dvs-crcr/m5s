package server

import (
    "context"
    "errors"
    "time"

    "m5s/domain"
    "m5s/internal/storage"
    "m5s/pkg/logger"
)

type Storage interface {
    Update(metric *domain.Metric) error
    GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList() ([]*domain.Metric, error)
    UpdateMetrics(metrics []*domain.Metric) error
    MyType() storage.StorageType
    Ping(ctx context.Context) error
}

type Config struct {
    storeInterval time.Duration
    restore       bool
}

type Service struct {
    cache   Storage
    storage Storage
    logger  logger.Logger
    config  Config
}

type Option func(*Service)

var (
    ErrDatabaseNoInit = errors.New(
        "database instance has not been initialized",
    )
)

func NewServerService(cache Storage, options ...Option) *Service {
    service := &Service{
        cache: cache,
    }

    for _, opt := range options {
        opt(service)
    }

    service.RestoreMetrics()
    go service.StartStoreTicker()

    return service
}

func WithLogger(logger logger.Logger) Option {
    return func(service *Service) {
        service.logger = logger
    }
}

func WithStorage(
    ctx context.Context,
    restore bool,
    fileStoragePath string,
    storeInterval time.Duration,
    dsn string,
    migrationsPath string,
    migrationsVersion string,
) Option {
    return func(service *Service) {
        service.config.restore = restore
        service.config.storeInterval = storeInterval

        switch {
        case dsn != "":
            service.storage = storage.NewFileStorage(fileStoragePath)

            var err error

            _, err = storage.NewDBStorage(
                ctx,
                dsn,
                migrationsPath,
                storage.ParseMigrationVersion(migrationsVersion),
            )
            if err != nil {
                service.logger.Fatal(
                    "unable to create new db instance",
                    "error", err,
                )
            }
        case fileStoragePath != "":
            service.storage = storage.NewFileStorage(fileStoragePath)
        default:
            service.storage = nil
        }
    }
}

func (ss *Service) Update(
    metricType string,
    name string,
    value string,
) error {
    metric, err := domain.NewMetric(metricType, name, value)
    if err != nil {
        return err
    }

    if err := ss.cache.Update(metric); err != nil {
        return err
    }

    if ss.config.storeInterval == 0 {
        if err := ss.BackupMetrics(); err != nil {
            return err
        }
    }

    return nil
}

func (ss *Service) GetMetric(
    metricType string,
    name string,
) (*domain.Metric, error) {
    switch metricType {
    case domain.MetricTypeGauge.String():
        metric, err := ss.cache.GetMetric(domain.MetricTypeGauge, name)
        if err != nil {
            return nil, err
        }

        return metric, nil
    case domain.MetricTypeCounter.String():
        metric, err := ss.cache.GetMetric(domain.MetricTypeCounter, name)
        if err != nil {
            return nil, err
        }

        return metric, nil
    }

    return nil, domain.ErrInvalidMetricType
}

func (ss *Service) GetMetricValue(
    metricType string,
    name string,
) (string, error) {
    switch metricType {
    case domain.MetricTypeGauge.String():
        metric, err := ss.cache.GetMetric(domain.MetricTypeGauge, name)
        if err != nil {
            return "", err
        }

        return metric.Value(), nil
    case domain.MetricTypeCounter.String():
        metric, err := ss.cache.GetMetric(domain.MetricTypeCounter, name)
        if err != nil {
            return "", err
        }

        return metric.Value(), nil
    }

    return "", domain.ErrInvalidMetricType
}

func (ss *Service) GetMetricsList() string {
    buffer := ""

    metricsList, err := ss.cache.GetMetricsList()
    if err != nil {
        ss.logger.Error("get metrics list", "error", err)
        return ""
    }

    for _, metric := range metricsList {
        buffer += metric.String()
    }

    return buffer
}

// RestoreMetrics is used to restore metrics from storage to memory cache
func (ss *Service) RestoreMetrics() {
    // Check if the storage exists
    if ss.storage == nil {
        return
    }

    storageMetrics, err := ss.storage.GetMetricsList()
    if err != nil {
        ss.logger.Error(
            "restore metrics",
            "error", err,
        )
    }

    if err := ss.cache.UpdateMetrics(storageMetrics); err != nil {
        ss.logger.Error("update metrics", "error", err.Error())
    }

    ss.logger.Info(
        "metrics have been successfully loaded from file",
        "data", storageMetrics,
    )
}

// BackupMetrics uses for extract metrics from memory cache to persistent storage
func (ss *Service) BackupMetrics() error {
    if ss.storage == nil {
        return nil
    }

    metricsList, err := ss.cache.GetMetricsList()
    if err != nil {
        return err
    }

    if err := ss.storage.UpdateMetrics(metricsList); err != nil {
        return err
    }

    return nil
}

func (ss *Service) PingDB(ctx context.Context) error {
    if ss.storage == nil {
        return nil
    }

    if ss.storage.MyType() != storage.TypeDatabase {
        return ErrDatabaseNoInit
    }

    return ss.storage.Ping(ctx)
}
