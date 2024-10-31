package models

import "strconv"

type Metrics struct {
    ID    string   `json:"id"`              // имя метрики
    MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
    Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
    Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) String() string {
    switch m.MType {
    case "gauge", "GAUGE":
        return strconv.FormatFloat(*m.Value, 'g', -1, 64)
    case "counter", "COUNTER":
        return strconv.FormatInt(*m.Delta, 10)
    default:
        return ""
    }
}
