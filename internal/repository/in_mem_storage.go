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
    ims.RLock()
    defer ims.RUnlock()

    ims.store[metric.Name] = metric

    return nil
}

func (ims *InMemStorage) GetMetric(
    metricType domain.MetricType,
    name string,
) (*domain.Metric, error) {
    if name == "" {
        return nil, domain.ErrInvalidMetricName
    }

    metric, ok := ims.store[name]
    if !ok {
        return nil, domain.ErrNoSuchMetric
    }

    if metric.Type != metricType {
        return nil, domain.ErrWrongMetricType
    }

    return metric, nil
}

func (ims *InMemStorage) GetMetricsList() []*domain.Metric {
    metrics := make([]*domain.Metric, len(ims.store))

    i := 0
    for _, metric := range ims.store {
        metrics[i] = metric
        i++
    }

    return metrics
}
