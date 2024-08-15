package agent

import (
    "context"
    "time"

    "m5s/domain"
    "m5s/pkg/logger"
)

type Storage interface {
    Update(ctx context.Context, metric *domain.Metric) error
    GetMetricsList(ctx context.Context) ([]*domain.Metric, error)
    UpdateMetrics(ctx context.Context, metrics []*domain.Metric) error
}

type Config struct {
    serverAddr     string
    pollInterval   time.Duration
    reportInterval time.Duration
}

type Service struct {
    storage Storage
    logger  logger.Logger
    config  Config
}

type Option func(*Service)

func NewAgentService(storage Storage, options ...Option) *Service {
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

func WithAddress(serverAddr string) Option {
    return func(service *Service) {
        service.config.serverAddr = serverAddr
    }
}

func WithPollInterval(pollInterval time.Duration) Option {
    return func(service *Service) {
        service.config.pollInterval = pollInterval
    }
}

func WithReportInterval(reportInterval time.Duration) Option {
    return func(service *Service) {
        service.config.reportInterval = reportInterval
    }
}
