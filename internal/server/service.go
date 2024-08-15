package server

import (
    "context"
    "errors"

    "m5s/domain"
    "m5s/pkg/logger"
)

type Storage interface {
    Update(ctx context.Context, metric *domain.Metric) error
    GetMetric(ctx context.Context, metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList(ctx context.Context) ([]*domain.Metric, error)
    UpdateMetrics(ctx context.Context, metrics []*domain.Metric) error
    Ping(ctx context.Context) error
}

type Service struct {
    storage Storage
    logger  logger.Logger
}

type Option func(*Service)

var (
    ErrDatabaseNoInit = errors.New(
        "database instance has not been initialized",
    )
)

func NewServerService(storage Storage, options ...Option) *Service {
    service := &Service{
        storage: storage,
    }

    for _, opt := range options {
        opt(service)
    }

    return service
}

func WithLogger(logger logger.Logger) Option {
    return func(service *Service) {
        service.logger = logger
    }
}

func (ss *Service) Update(
    ctx context.Context,
    metricType string,
    name string,
    value string,
) error {
    metric, err := domain.NewMetric(metricType, name, value)
    if err != nil {
        return err
    }

    if err := ss.storage.Update(ctx, metric); err != nil {
        return err
    }

    return nil
}

func (ss *Service) GetMetric(
    ctx context.Context,
    metricType string,
    name string,
) (*domain.Metric, error) {
    mt, err := domain.ParseMetricType(metricType)
    if err != nil {
        return nil, err
    }

    metric, err := ss.storage.GetMetric(ctx, mt, name)
    if err != nil {
        return nil, err
    }

    return metric, nil
}

func (ss *Service) GetMetricValue(
    ctx context.Context,
    metricType string,
    name string,
) (string, error) {
    mt, err := domain.ParseMetricType(metricType)
    if err != nil {
        return "", err
    }

    metric, err := ss.storage.GetMetric(ctx, mt, name)
    if err != nil {
        return "", err
    }

    return metric.Value(), nil
}

func (ss *Service) GetMetricsList(ctx context.Context) string {
    var buffer string

    metricsList, err := ss.storage.GetMetricsList(ctx)
    if err != nil {
        ss.logger.Error("get metrics list from cache", "error", err)
        return ""
    }

    for _, metric := range metricsList {
        buffer += metric.String()
    }

    return buffer
}

func (ss *Service) PingDB(ctx context.Context) error {
    return ss.storage.Ping(ctx)
}
