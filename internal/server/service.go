package server

import (
    "fmt"

    "m5s/domain"
)

type Repo interface {
    Update(metric *domain.Metric) error
    GetMetric(metricType domain.MetricType, name string) (*domain.Metric, error)
    GetMetricsList() []*domain.Metric
}

type Service struct {
    repo Repo
}

func NewServerService(repo Repo) *Service {
    return &Service{
        repo: repo,
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

func (ss *Service) GetMetricsList() string {
    buffer := ""

    metricsList := ss.repo.GetMetricsList()
    for _, metric := range metricsList {
        buffer += fmt.Sprintf(
            "%s(%s)=%s\n", metric.Name, metric.Type, metric,
        )
    }

    return buffer
}
