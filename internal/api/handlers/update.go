package handlers

import (
    "encoding/json"
    "io"
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    metricType := chi.URLParam(r, "metricType")
    metricName := chi.URLParam(r, "metricName")
    metricValue := chi.URLParam(r, "metricValue")

    if err := h.serverService.Update(metricType, metricName, metricValue); err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
    if !api.CheckContentType(
        r.Header.Get("Content-Type"),
        "application/json",
    ) {
        api.HandleErrors(api.ErrInvalidJSONContentType, w)

        return
    }

    metric := &models.Metrics{}

    bodyData, err := io.ReadAll(r.Body)
    if err != nil {
        h.logger.Error("read body", "error", err)

        return
    }

    if err := json.Unmarshal(bodyData, &metric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }

    if err := h.serverService.Update(
        metric.MType,
        metric.ID,
        metric.String(),
    ); err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
