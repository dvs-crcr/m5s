package storage

import (
    "context"
    "encoding/json"
    "errors"
    "os"

    "m5s/domain"
)

type FileStorage struct {
    fileStoragePath string
    storageType     StorageType
}

func NewFileStorage(fileStoragePath string) *FileStorage {
    return &FileStorage{
        fileStoragePath: fileStoragePath,
        storageType:     TypeFile,
    }
}

func (ifs *FileStorage) MyType() StorageType {
    return ifs.storageType
}

func (ifs *FileStorage) UpdateMetrics(metrics []*domain.Metric) error {
    metricsBytes, err := json.Marshal(metrics)
    if err != nil {
        return err
    }

    if _, err := os.Stat(ifs.fileStoragePath); err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            return err
        }

        if _, err := os.Create(ifs.fileStoragePath); err != nil {
            return err
        }
    }

    if err := os.WriteFile(
        ifs.fileStoragePath,
        metricsBytes,
        0666,
    ); err != nil {
        return err
    }

    return nil
}

func (ifs *FileStorage) GetMetricsList() ([]*domain.Metric, error) {
    metrics := make([]*domain.Metric, 0)

    fileBytes, err := os.ReadFile(ifs.fileStoragePath)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(fileBytes, &metrics); err != nil {
        return nil, err
    }

    return metrics, nil
}

func (ifs *FileStorage) Update(metric *domain.Metric) error {
    metricsList, err := ifs.GetMetricsList()
    if err != nil {
        return err
    }

    for _, m := range metricsList {
        if m.Name == metric.Name && m.Type == metric.Type {
            m = metric
        }
    }

    if err := ifs.UpdateMetrics(metricsList); err != nil {
        return err
    }

    return nil
}

func (ifs *FileStorage) GetMetric(
    metricType domain.MetricType,
    name string,
) (*domain.Metric, error) {
    metricsList, err := ifs.GetMetricsList()
    if err != nil {
        return nil, err
    }

    for _, m := range metricsList {
        if m.Name == name && m.Type == metricType {
            return m, nil
        }
    }

    return nil, nil
}

func (ifs *FileStorage) Ping(_ context.Context) error {
    return nil
}
