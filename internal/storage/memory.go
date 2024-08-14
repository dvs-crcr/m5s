package storage

import (
    "fmt"
    "sync"

    "m5s/domain"
)

type MemStorage struct {
    sync.RWMutex
    store       map[string]*domain.Metric
    storageType StorageType
}

func NewMemStorage() *MemStorage {
    return &MemStorage{
        store:       make(map[string]*domain.Metric),
        storageType: TypeMemory,
    }
}

func (ims *MemStorage) MyType() StorageType {
    return ims.storageType
}

func (ims *MemStorage) Update(metric *domain.Metric) error {
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

func (ims *MemStorage) UpdateMetrics(metrics []*domain.Metric) error {
    for _, metric := range metrics {
        if err := ims.Update(metric); err != nil {
            return fmt.Errorf("unable to restore metric: %v: %w", metric, err)
        }
    }

    return nil
}

func (ims *MemStorage) GetMetric(
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

func (ims *MemStorage) GetMetricsList() ([]*domain.Metric, error) {
    ims.RLock()
    defer ims.RUnlock()

    metrics := make([]*domain.Metric, 0, len(ims.store))

    for _, metric := range ims.store {
        metrics = append(metrics, metric)
    }

    return metrics, nil
}
