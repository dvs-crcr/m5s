package domain

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricName  = errors.New("invalid metric name")
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrNoSuchMetric       = errors.New("no such metric")
)

type MetricType int

const (
	MetricTypeGauge MetricType = iota
	MetricTypeCounter
)

type Metric struct {
	Name       string
	Type       MetricType
	FloatValue float64
	IntValue   int64
}

func (mt MetricType) String() string {
	return [...]string{"gauge", "counter"}[mt]
}

// NewMetric uses to create new Metric instance.
func NewMetric(
	metricType string,
	name string,
	value string,
) (*Metric, error) {
	switch metricType {
	case MetricTypeCounter.String():
		parsedValue, err := validateCounter(name, value)
		if err != nil {
			return nil, err
		}

		return &Metric{
			Name:       name,
			Type:       MetricTypeCounter,
			FloatValue: 0,
			IntValue:   parsedValue,
		}, nil
	case MetricTypeGauge.String():
		parsedValue, err := validateGauge(name, value)
		if err != nil {
			return nil, err
		}

		return &Metric{
			Name:       name,
			Type:       MetricTypeGauge,
			FloatValue: parsedValue,
			IntValue:   0,
		}, nil
	}

	return nil, ErrInvalidMetricType
}

func NewGauge(name string, value float64) *Metric {
	return &Metric{
		Name:       name,
		Type:       MetricTypeGauge,
		FloatValue: value,
		IntValue:   0,
	}
}

func NewCounter(name string, value int64) *Metric {
	return &Metric{
		Name:       name,
		Type:       MetricTypeCounter,
		FloatValue: 0,
		IntValue:   value,
	}
}

func (m *Metric) String() string {
	switch m.Type {
	case MetricTypeGauge:
		return strconv.FormatFloat(m.FloatValue, 'g', -1, 64)
	case MetricTypeCounter:
		return strconv.FormatInt(m.IntValue, 10)
	default:
		return ""
	}
}

func validateCounter(name string, value string) (int64, error) {
	if name == "" {
		return 0, ErrInvalidMetricName
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return parsedValue, fmt.Errorf("%w: %w", ErrInvalidMetricValue, err)
	}

	return parsedValue, nil
}

func validateGauge(name string, value string) (float64, error) {
	if name == "" {
		return 0, ErrInvalidMetricName
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return parsedValue, fmt.Errorf("%w: %w", ErrInvalidMetricValue, err)
	}

	return parsedValue, nil
}
