package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api"
)

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    w.Header().Set("Content-Type", "text/plain")

    metricType := chi.URLParam(r, "metricType")
    metricName := chi.URLParam(r, "metricName")

    metricValue, err := h.serverService.GetMetricValue(ctx, metricType, metricName)
    if err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricValue))
}
