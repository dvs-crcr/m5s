package server

import (
    "time"

    "m5s/domain"
    "m5s/pkg/logger"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList() []*domain.Metric
}

type Config struct {
    storeInterval   time.Duration
    fileStoragePath string
    restore         bool
}

type Service struct {
    repo   Repo
    logger logger.Logger
    config Config
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

func WithFileStoragePath(fileStoragePath string) Option {
    return func(service *Service) {
        service.config.fileStoragePath = fileStoragePath
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

    return ss.repo.Update(metric)
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

func (ss *Service) RestoreData() {
    if !ss.config.restore {
        return
    }

    ss.logger.Info(
        "Restoring data", "src", ss.config.fileStoragePath,
        // TODO: implement backup
    )
}
