package repository

import (
    "sync"

    "m5s/domain"
)

type InMemStorage struct {
    sync.RWMutex
    store map[string]*domain.Metric
}

func NewInMemStorage() *InMemStorage {
    return &InMemStorage{
        store: make(map[string]*domain.Metric),
    }
}

func (ims *InMemStorage) Update(metric *domain.Metric) error {
    ims.Lock()
    defer ims.Unlock()

    if metric.Type == domain.MetricTypeCounter {
        newValue := metric.IntValue

        prevMetric, ok := ims.store[metric.Name]
        if ok {
            newValue += prevMetric.IntValue

            metric.IntValue = newValue

            ims.store[metric.Name] = metric

            return nil
        }
    }

    ims.store[metric.Name] = metric

    return nil
}

func (ims *InMemStorage) GetMetric(
    metricType domain.MetricType,
    name string,
) (*domain.Metric, error) {
    ims.RLock()
    defer ims.RUnlock()

    if name == "" {
        return nil, domain.ErrInvalidMetricName
    }

    metric, ok := ims.store[name]
    if !ok {
        return nil, domain.ErrNoSuchMetric
    }

    if metric.Type != metricType {
        return nil, domain.ErrInvalidMetricType
    }

    return metric, nil
}

func (ims *InMemStorage) GetMetricsList() []*domain.Metric {
    ims.RLock()
    defer ims.RUnlock()

    metrics := make([]*domain.Metric, len(ims.store))

    i := 0

    for _, metric := range ims.store {
        metrics[i] = metric
        i++
    }

    return metrics
}
