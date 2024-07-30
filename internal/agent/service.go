package agent

import (
    "time"

    "m5s/domain"
    "m5s/internal/models"
    "m5s/internal/repository"
    "m5s/pkg/logger"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetricsList() []*domain.Metric
}

type Service struct {
    repo           Repo
    stat           *models.Statistics
    logger         logger.Logger
    serverAddr     string
    pollInterval   time.Duration
    reportInterval time.Duration
}

type Option func(*Service)

func NewAgentService(options ...Option) *Service {
    repo := repository.NewInMemStorage()
    stat := models.NewStatistics()

    service := &Service{
        repo: repo,
        stat: stat,
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
        service.serverAddr = serverAddr
    }
}

func WithPollInterval(pollInterval time.Duration) Option {
    return func(service *Service) {
        service.pollInterval = pollInterval
    }
}

func WithReportInterval(reportInterval time.Duration) Option {
    return func(service *Service) {
        service.reportInterval = reportInterval
    }
}
