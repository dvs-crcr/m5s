package api

import (
    "errors"
    "net/http"
    "strings"

    "m5s/domain"
    "m5s/internal/repository"
    "m5s/internal/server"
)

type Handler struct {
    serverService *server.Service
    Mux           *http.ServeMux
}

func NewHandler() *Handler {
    return &Handler{
        serverService: server.NewServerService(
            repository.NewInMemStorage(),
        ),
    }
}

func checkDefaultErrors(err error, w http.ResponseWriter) {
    if errors.Is(err, domain.ErrInvalidMetricType) {
        w.WriteHeader(http.StatusBadRequest)

        return
    } else if errors.Is(err, domain.ErrInvalidMetricName) {
        w.WriteHeader(http.StatusNotFound)

        return
    } else if errors.Is(err, domain.ErrInvalidMetricValue) {
        w.WriteHeader(http.StatusBadRequest)

        return
    } else if errors.Is(err, domain.ErrNoSuchMetric) {
        w.WriteHeader(http.StatusNotFound)

        return
    }
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 5 {
        w.WriteHeader(http.StatusNotFound)

        return
    }

    if err := h.serverService.Update(segments[2], segments[3], segments[4]); err != nil {
        checkDefaultErrors(err, w)
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 4 {
        w.WriteHeader(http.StatusNotFound)

        return
    }

    metricValue, err := h.serverService.GetMetric(segments[2], segments[3])
    if err != nil {
        checkDefaultErrors(err, w)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricValue))
}

func (h *Handler) GetMetricsList(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")

    metricsList := h.serverService.GetMetricsList()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricsList))
}
