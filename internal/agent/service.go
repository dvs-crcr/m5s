package agent

import (
    "fmt"
    "net/http"
    "time"

    "m5s/domain"
    "m5s/internal/logger"
    "m5s/internal/repository"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetricsList() []*domain.Metric
}

type Service struct {
    repo           Repo
    stat           *domain.Statistics
    logger         logger.Logger
    serverAddr     string
    pollInterval   time.Duration
    reportInterval time.Duration
}

type Option func(*Service)

func NewAgentService(options ...Option) *Service {
    repo := repository.NewInMemStorage()

    service := &Service{
        repo: repo,
        stat: domain.NewStatistics(),
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

func (as *Service) StartPoller() {
    as.logger.Info("Start poller", "duration", as.pollInterval)

    ticker := time.NewTicker(as.pollInterval)

    for range ticker.C {
        as.stat.Refresh()

        for name, value := range as.stat.CurrentValues {
            metric := domain.NewGauge(name, value)
            if err := as.repo.Update(metric); err != nil {
                as.logger.Fatal("update gauge", "error", err)
            }
        }

        pollCountMetric := domain.NewCounter("PollCount", 1)
        if err := as.repo.Update(pollCountMetric); err != nil {
            as.logger.Fatal("update counter", "error", err)
        }
    }
}

func (as *Service) StartReporter() {
    as.logger.Info("Start reporter", "duration", as.reportInterval)

    ticker := time.NewTicker(as.reportInterval)

    for range ticker.C {
        for _, metric := range as.repo.GetMetricsList() {
            if err := as.makeRequest(metric); err != nil {
                as.logger.Fatal("make reporter request", "error", err)
            }
        }
    }
}

func (as *Service) makeRequest(metric *domain.Metric) error {
    request, err := http.NewRequest(
        http.MethodPost,
        fmt.Sprintf(
            "http://%s/update/%s/%s/%s",
            as.serverAddr, metric.Type, metric.Name, metric,
        ),
        nil,
    )
    if err != nil {
        return fmt.Errorf("execute http request: %v", err)
    }

    request.Header.Set("Content-Type", "text/plain")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    return nil
}
