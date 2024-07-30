package agent

import (
    "time"

    "m5s/domain"
    "m5s/pkg/logger"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetricsList() []*domain.Metric
}

type Config struct {
    serverAddr     string
    pollInterval   time.Duration
    reportInterval time.Duration
}

type Service struct {
    repo   Repo
    logger logger.Logger
    config Config
}

type Option func(*Service)

func NewAgentService(repo Repo, options ...Option) *Service {
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
