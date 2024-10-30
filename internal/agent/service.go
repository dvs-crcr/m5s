package agent

import (
    "context"
    "fmt"
    "time"

    "m5s/domain"
    internalLogger "m5s/pkg/logger"
)

var logger = internalLogger.GetLogger()

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
    config  Config
}

type Option func(*Service)

func NewAgentService(storage Storage, options ...Option) *Service {
    logger = logger.With("package", "agent")

    service := &Service{
        storage: storage,
    }

    for _, opt := range options {
        opt(service)
    }

    logger.Infow(
        "init new agent service",
        "config", fmt.Sprintf("%+v", service.config),
    )

    return service
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
