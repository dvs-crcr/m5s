package repository

import (
    "m5s/domain"
)

type InFileStorage struct {
    fileStoragePath string
}

// TODO: implement this

func NewInFileStorage(fileStoragePath string) *InFileStorage {
    return &InFileStorage{
        fileStoragePath: fileStoragePath,
    }
}

func (ifs *InFileStorage) Update(metric *domain.Metric) error {
    return nil
}

func (ifs *InFileStorage) GetMetric(
    metricType domain.MetricType,
    name string,
) (*domain.Metric, error) {
    return nil, nil
}

func (ifs *InFileStorage) GetMetricsList() []*domain.Metric {
    return nil
}
