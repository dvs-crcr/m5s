package repository

import (
    "encoding/json"
    "errors"
    "io"
    "os"

    "m5s/domain"
)

type InFileStorage struct {
    fileStoragePath string
}

func NewInFileStorage(fileStoragePath string) *InFileStorage {
    return &InFileStorage{
        fileStoragePath: fileStoragePath,
    }
}

func (ifs *InFileStorage) UpdateMetrics(metrics []*domain.Metric) error {
    file, err := os.OpenFile(
        ifs.fileStoragePath,
        os.O_CREATE|os.O_WRONLY,
        0666,
    )
    if err != nil {
        return err
    }

    metricsBytes, err := json.Marshal(metrics)
    if err != nil {
        return err
    }

    if _, err := file.Write(metricsBytes); err != nil {
        return err
    }

    return nil
}

func (ifs *InFileStorage) GetMetricsList() ([]*domain.Metric, error) {
    metrics := make([]*domain.Metric, 0)

    file, err := os.OpenFile(
        ifs.fileStoragePath,
        os.O_RDONLY,
        0666,
    )
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return metrics, nil
        }
        return nil, err
    }

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(fileBytes, &metrics); err != nil {
        return nil, err
    }

    return metrics, nil
}
