package storage

import (
    "encoding/json"
    "errors"
    "os"

    "m5s/domain"
)

type FileStorage struct {
    fileStoragePath string
}

func NewFileStorage(fileStoragePath string) *FileStorage {
    return &FileStorage{
        fileStoragePath: fileStoragePath,
    }
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
