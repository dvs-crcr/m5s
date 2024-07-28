package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "m5s/domain"
    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodGet {
        api.HandleErrors(api.ErrInvalidMethod, w)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 4 {
        api.HandleErrors(api.ErrInvalidSegmentsCount, w)

        return
    }

    metricValue, err := h.serverService.GetMetricValue(segments[2], segments[3])
    if err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricValue))
}

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
    if !api.CheckContentType(
        r.Header.Get("Content-Type"),
        "application/json",
    ) {
        api.HandleErrors(api.ErrInvalidJSONContentType, w)

        return
    }

    if r.Method != http.MethodPost {
        api.HandleErrors(api.ErrInvalidMethod, w)

        return
    }

    w.Header().Set("Content-Type", "application/json")

    metric := &models.Metrics{}

    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&metric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }

    domainMetric, err := h.serverService.GetMetric(metric.MType, metric.ID)
    if err != nil {
        api.HandleErrors(err, w)

        return
    }

    modelMetric := &models.Metrics{}

    switch domainMetric.Type {
    case domain.MetricTypeGauge:
        modelMetric = &models.Metrics{
            ID:    domainMetric.Name,
            MType: domainMetric.Type.String(),
            Value: &domainMetric.FloatValue,
        }
    case domain.MetricTypeCounter:
        modelMetric = &models.Metrics{
            ID:    domainMetric.Name,
            MType: domainMetric.Type.String(),
            Delta: &domainMetric.IntValue,
        }
    }

    if err := json.NewEncoder(w).Encode(modelMetric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }
}
