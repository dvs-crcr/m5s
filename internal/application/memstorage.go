package application

import (
    "errors"
    "fmt"
    "strconv"
    "sync"
)

var (
    ErrInvalidMetricType  = errors.New("invalid metric type")
    ErrInvalidMetricName  = errors.New("invalid metric name")
    ErrInvalidMetricValue = errors.New("invalid metric value")
)

type MemStorage struct {
    sync.RWMutex
    Gauge   map[string]float64
    Counter map[string]int64
}

func NewMemStorage() *MemStorage {
    return &MemStorage{
        Gauge:   make(map[string]float64),
        Counter: make(map[string]int64),
    }
}

func (ms *MemStorage) validateMetricType(metricType string) bool {
    switch metricType {
    case "gauge":
        return true
    case "counter":
        return true
    }

    return false
}

func (ms *MemStorage) validateCounter(name string, value string) (int64, error) {
    if name == "" {
        return 0, ErrInvalidMetricName
    }

    parsedValue, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        return parsedValue, fmt.Errorf("%w: %w", ErrInvalidMetricValue, err)
    }

    return parsedValue, nil
}

func (ms *MemStorage) validateGauge(name string, value string) (float64, error) {
    if name == "" {
        return 0, ErrInvalidMetricName
    }

    parsedValue, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return parsedValue, fmt.Errorf("%w: %w", ErrInvalidMetricValue, err)
    }

    return parsedValue, nil
}

func (ms *MemStorage) updateCounter(name string, value string) error {
    parsedValue, err := ms.validateCounter(name, value)
    if err != nil {
        return err
    }

    ms.Counter[name] = parsedValue

    return nil
}

func (ms *MemStorage) updateGauge(name string, value string) error {
    parsedValue, err := ms.validateGauge(name, value)
    if err != nil {
        return err
    }

    ms.Gauge[name] = parsedValue

    return nil
}

func (ms *MemStorage) Update(
    metricType string,
    name string,
    value string,
) error {
    ms.RLock()
    defer ms.RUnlock()

    switch metricType {
    case "counter":
        return ms.updateCounter(name, value)
    case "gauge":
        return ms.updateGauge(name, value)
    }

    return ErrInvalidMetricType
}
