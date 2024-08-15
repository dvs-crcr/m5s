package domain

import (
    "database/sql/driver"
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
    MetricTypeUnknown MetricType = iota
    MetricTypeGauge
    MetricTypeCounter
)

type Metric struct {
    Name       string     `json:"name"`
    Type       MetricType `json:"type"`
    FloatValue float64    `json:"float_value"`
    IntValue   int64      `json:"int_value"`
}

func (mt MetricType) String() string {
    return [...]string{"unknown", "gauge", "counter"}[mt]
}

func ParseMetricType(strMetricType string) (MetricType, error) {
    switch strMetricType {
    case "gauge", "GAUGE":
        return MetricTypeGauge, nil
    case "counter", "COUNTER":
        return MetricTypeCounter, nil
    default:
        return MetricTypeUnknown, ErrInvalidMetricType
    }
}

func (mt *MetricType) Scan(value any) error {
    if value == nil {
        return nil
    }

    sv, err := driver.String.ConvertValue(value)
    if err != nil {
        return fmt.Errorf("cannot scan value. %w", err)
    }

    v, ok := sv.(string)
    if !ok {
        return errors.New("cannot scan value. cannot convert value to string")
    }

    metricType, err := ParseMetricType(v)
    if err != nil {
        return errors.New("cannot parse string to MetricType")
    }

    *mt = metricType

    return nil
}

// NewMetric uses to create new Metric instance.
func NewMetric(
    metricType string,
    name string,
    value string,
) (*Metric, error) {
    mt, err := ParseMetricType(metricType)
    if err != nil {
        return nil, err
    }

    switch mt {
    case MetricTypeCounter:
        parsedValue, err := validateCounter(name, value)
        if err != nil {
            return nil, err
        }

        return NewCounter(name, parsedValue), nil
    case MetricTypeGauge:
        parsedValue, err := validateGauge(name, value)
        if err != nil {
            return nil, err
        }

        return NewGauge(name, parsedValue), nil
    default:
        return nil, ErrInvalidMetricType
    }
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

func (m Metric) Value() string {
    switch m.Type {
    case MetricTypeGauge:
        return strconv.FormatFloat(m.FloatValue, 'g', -1, 64)
    case MetricTypeCounter:
        return strconv.FormatInt(m.IntValue, 10)
    default:
        return ""
    }
}

func (m Metric) String() string {
    return fmt.Sprintf(
        "%s(%s)=%s\n", m.Name, m.Type, m.Value(),
    )
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
