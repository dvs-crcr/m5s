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

func (ims *InMemStorage) GetAll() []*domain.Metric {
    metrics := make([]*domain.Metric, len(ims.store))

    i := 0
    for _, metric := range ims.store {
        metrics[i] = metric
        i++
    }

    return metrics
}
