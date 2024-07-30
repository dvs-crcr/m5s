package server

import (
    "time"

    "m5s/domain"
    "m5s/internal/repository"
    "m5s/pkg/logger"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList() []*domain.Metric
    UpdateMetrics(metrics []*domain.Metric) error
}

type Storage interface {
    GetMetricsList() ([]*domain.Metric, error)
    UpdateMetrics(metrics []*domain.Metric) error
}

type Config struct {
    storeInterval time.Duration
    restore       bool
}

type Service struct {
    repo    Repo
    storage Storage
    logger  logger.Logger
    config  Config
}

type Option func(*Service)

func NewServerService(repo Repo, options ...Option) *Service {
    service := &Service{
        repo: repo,
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

func WithStoreInterval(storeInterval time.Duration) Option {
    return func(service *Service) {
        service.config.storeInterval = storeInterval
    }
}

func WithStorage(fileStoragePath string) Option {
    return func(service *Service) {
        service.storage = repository.NewInFileStorage(fileStoragePath)
    }
}

func WithRestore(restore bool) Option {
    return func(service *Service) {
        service.config.restore = restore
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

    if err := ss.repo.Update(metric); err != nil {
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
        metric, err := ss.repo.GetMetric(domain.MetricTypeGauge, name)
        if err != nil {
            return nil, err
        }

        return metric, nil
    case domain.MetricTypeCounter.String():
        metric, err := ss.repo.GetMetric(domain.MetricTypeCounter, name)
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
        metric, err := ss.repo.GetMetric(domain.MetricTypeGauge, name)
        if err != nil {
            return "", err
        }

        return metric.Value(), nil
    case domain.MetricTypeCounter.String():
        metric, err := ss.repo.GetMetric(domain.MetricTypeCounter, name)
        if err != nil {
            return "", err
        }

        return metric.Value(), nil
    }

    return "", domain.ErrInvalidMetricType
}

func (ss *Service) GetMetricsList() string {
    buffer := ""

    metricsList := ss.repo.GetMetricsList()
    for _, metric := range metricsList {
        buffer += metric.String()
    }

    return buffer
}

func (ss *Service) RestoreMetrics() {
    if !ss.config.restore {
        return
    }

    storageMetrics, err := ss.storage.GetMetricsList()
    if err != nil {
        ss.logger.Error(
            "restore metrics",
            "error", err,
        )
    }

    if err := ss.repo.UpdateMetrics(storageMetrics); err != nil {
        ss.logger.Error(err.Error())
    }
}

func (ss *Service) BackupMetrics() error {
    metrics := ss.repo.GetMetricsList()

    if err := ss.storage.UpdateMetrics(metrics); err != nil {
        return err
    }

    return nil
}
