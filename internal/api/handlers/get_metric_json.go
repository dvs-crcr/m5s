package handlers

import (
    "encoding/json"
    "net/http"

    "m5s/domain"
    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if !api.CheckContentType(
        r.Header.Get("Content-Type"),
        "application/json",
    ) {
        api.HandleErrors(api.ErrInvalidJSONContentType, w)

        return
    }

    metric := &models.Metrics{}

    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&metric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }

    domainMetric, err := h.serverService.GetMetric(ctx, metric.MType, metric.ID)
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(modelMetric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }
}
