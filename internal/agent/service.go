package agent

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "net/http"
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

func (as *Service) StartPoller() {
    as.logger.Info("Starting poller", "duration", as.pollInterval)

    ticker := time.NewTicker(as.pollInterval)

    for range ticker.C {
        as.stat.Refresh()

        for name, value := range as.stat.CurrentValues {
            metric := domain.NewGauge(name, value)
            if err := as.repo.Update(metric); err != nil {
                as.logger.Error("update gauge", "error", err)
            }
        }

        pollCountMetric := domain.NewCounter("PollCount", 1)
        if err := as.repo.Update(pollCountMetric); err != nil {
            as.logger.Error("update counter", "error", err)
        }
    }
}

func (as *Service) StartReporter() {
    var client = &http.Client{
        Timeout:   time.Second * 1,
        Transport: &http.Transport{},
    }

    as.logger.Info("Starting reporter", "duration", as.reportInterval)

    ticker := time.NewTicker(as.reportInterval)

    for range ticker.C {
        for _, metric := range as.repo.GetMetricsList() {
            if err := as.makeRequest(client, metric); err != nil {
                as.logger.Error("make reporter request", "error", err)
            }
        }
    }
}

func (as *Service) makeRequest(client *http.Client, metric *domain.Metric) error {
    modelMetric := &models.Metrics{
        ID:    metric.Name,
        MType: metric.Type.String(),
        Delta: &metric.IntValue,
        Value: &metric.FloatValue,
    }

    bytesMetric, err := json.Marshal(modelMetric)
    if err != nil {
        return err
    }

    var buf bytes.Buffer

    zw := gzip.NewWriter(&buf)
    if _, err = zw.Write(bytesMetric); err != nil {
        return err
    }
    zw.Close()

    request, err := http.NewRequest(
        http.MethodPost,
        fmt.Sprintf("http://%s/update/", as.serverAddr),
        &buf,
    )
    if err != nil {
        return fmt.Errorf("execute http request: %v", err)
    }

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Content-Encoding", "gzip")

    response, err := client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    //if _, err := io.Copy(io.Discard, response.Body); err != nil {
    //    return err
    //}

    return nil
}
