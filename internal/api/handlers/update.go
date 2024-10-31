package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    w.Header().Set("Content-Type", "text/plain")

    metricType := chi.URLParam(r, "metricType")
    metricName := chi.URLParam(r, "metricName")
    metricValue := chi.URLParam(r, "metricValue")

    if err := h.serverService.Update(ctx, metricType, metricName, metricValue); err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
