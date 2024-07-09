package server

import (
    "fmt"

    "m5s/domain"
)

type ServerRepo interface {
    Update(metric *domain.Metric) error
    GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList() []*domain.Metric
}

type ServerService struct {
    repo ServerRepo
}

func NewServerService(repo ServerRepo) *ServerService {
    return &ServerService{
        repo: repo,
    }
}

func (ss *ServerService) Update(
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

func (ss *ServerService) GetMetric(
    metricType string,
    name string,
) (string, error) {
    switch metricType {
    case domain.MetricTypeGauge.String():
        metric, err := ss.repo.GetMetric(domain.MetricTypeGauge, name)
        if err != nil {
            return "", err
        }

        return metric.String(), nil
    case domain.MetricTypeCounter.String():
        metric, err := ss.repo.GetMetric(domain.MetricTypeCounter, name)
        if err != nil {
            return "", err
        }

        return metric.String(), nil
    }

    return "", domain.ErrInvalidMetricType
}

func (ss *ServerService) GetMetricsList() string {
    buffer := ""

    metricsList := ss.repo.GetMetricsList()
    for _, metric := range metricsList {
        buffer += fmt.Sprintf(
            "%s(%s)=%s\n", metric.Name, metric.Type, metric,
        )
    }

    return buffer
}
