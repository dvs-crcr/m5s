package fileStorage

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/fs"
    "os"
    "time"

    "m5s/domain"
    memoryStorage "m5s/internal/storage/memory_storage"
    "m5s/pkg/logger"
)

type FileStorage struct {
    cache           *memoryStorage.MemStorage
    fileStoragePath string
    storeInterval   time.Duration
    restore         bool
    logger          logger.Logger
}

func NewFileStorage(
    ctx context.Context,
    logger logger.Logger,
    fileStoragePath string,
    storeInterval time.Duration,
    restore bool,
) (*FileStorage, error) {
    ifs := &FileStorage{
        cache:           memoryStorage.NewMemStorage(logger),
        logger:          logger,
        fileStoragePath: fileStoragePath,
        storeInterval:   storeInterval,
        restore:         restore,
    }

    if err := ifs.restoreMetrics(ctx); err != nil {
        return nil, fmt.Errorf("restore metrics: %w", err)
    }
    go ifs.startStoreTicker(ctx)

    return ifs, nil
}

func (ifs *FileStorage) _(ctx context.Context, metrics []*domain.Metric) error {
    if ifs.storeInterval == 0 {
        if err := ifs.storeMetricsToFile(metrics); err != nil {
            return err
        }
    }

    return ifs.cache.UpdateMetrics(ctx, metrics)
}

// UpdateMetrics uses to store metrics
func (ifs *FileStorage) UpdateMetrics(ctx context.Context, metrics []*domain.Metric) error {
    if err := ifs.cache.UpdateMetrics(ctx, metrics); err != nil {
        return err
    }

    if ifs.storeInterval == 0 {
        if err := ifs.storeMetricsToFile(metrics); err != nil {
            return err
        }
    }

    return nil
}

// storeMetricsToFile uses to store metrics to file
func (ifs *FileStorage) storeMetricsToFile(metrics []*domain.Metric) error {
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

// getMetricsFromFile uses to get metrics from file
func (ifs *FileStorage) getMetricsFromFile() ([]*domain.Metric, error) {
    metrics := make([]*domain.Metric, 0)

    fileBytes, err := os.ReadFile(ifs.fileStoragePath)
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            return nil, nil
        }

        return nil, err
    }

    if err := json.Unmarshal(fileBytes, &metrics); err != nil {
        return nil, err
    }

    return metrics, nil
}

// GetMetricsList uses to get metrics from cache
func (ifs *FileStorage) GetMetricsList(ctx context.Context) ([]*domain.Metric, error) {
    return ifs.cache.GetMetricsList(ctx)
}

func (ifs *FileStorage) Update(ctx context.Context, metric *domain.Metric) error {
    if err := ifs.cache.Update(ctx, metric); err != nil {
        return err
    }

    if ifs.storeInterval == 0 {
        metricsList, err := ifs.cache.GetMetricsList(ctx)
        if err != nil {
            return err
        }

        for _, m := range metricsList {
            if m.Name == metric.Name && m.Type == metric.Type {
                m = metric
            }
        }

        if err := ifs.storeMetricsToFile(metricsList); err != nil {
            return err
        }
    }

    return nil
}

// GetMetric uses to get metric from cache
func (ifs *FileStorage) GetMetric(
    ctx context.Context,
    metricType domain.MetricType,
    name string,
) (*domain.Metric, error) {
    return ifs.cache.GetMetric(ctx, metricType, name)
}

// Ping - storage interface stub
func (ifs *FileStorage) Ping(_ context.Context) error {
    return nil
}

// restoreMetrics is used to restore metrics from storage to memory cache
func (ifs *FileStorage) restoreMetrics(ctx context.Context) error {
    if !ifs.restore {
        return nil
    }

    storageMetrics, err := ifs.getMetricsFromFile()
    if err != nil {
        return err
    }

    if err := ifs.cache.UpdateMetrics(ctx, storageMetrics); err != nil {
        return err
    }

    return nil
}

// backupMetrics uses for extract metrics from memory cache to persistent storage
func (ifs *FileStorage) backupMetrics(ctx context.Context) error {
    metricsList, err := ifs.cache.GetMetricsList(ctx)
    if err != nil {
        return err
    }

    if err := ifs.storeMetricsToFile(metricsList); err != nil {
        return err
    }

    return nil
}
