package api

import (
    "encoding/json"
    "net/http"
    "strings"

    "m5s/internal/models"
)

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodGet {
        handleErrors(ErrInvalidMethod, w)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 4 {
        handleErrors(ErrInvalidSegmentsCount, w)

        return
    }

    metricValue, err := h.serverService.GetMetricValue(segments[2], segments[3])
    if err != nil {
        handleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricValue))
}

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "application/json" {
        handleErrors(ErrInvalidJSONContentType, w)

        return
    }

    if r.Method != http.MethodPost {
        handleErrors(ErrInvalidMethod, w)

        return
    }

    w.Header().Set("Content-Type", "application/json")

    metric := &models.Metrics{}

    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&metric); err != nil {
        handleErrors(ErrInvalidJSONStruct, w)

        return
    }

    domainMetric, err := h.serverService.GetMetric(metric.MType, metric.ID)
    if err != nil {
        handleErrors(err, w)

        return
    }

    modelMetric := &models.Metrics{
        ID:    domainMetric.Name,
        MType: domainMetric.Type.String(),
        Delta: &domainMetric.IntValue,
        Value: &domainMetric.FloatValue,
    }

    if err := json.NewEncoder(w).Encode(modelMetric); err != nil {
        handleErrors(ErrInvalidJSONStruct, w)

        return
    }
}
