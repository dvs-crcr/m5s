package api

import (
    "errors"
    "net/http"
    "strings"

    "m5s/domain"
    "m5s/internal/repository"
    "m5s/server"
)

type Handler struct {
    serverService *server.ServerService
    Mux           *http.ServeMux
}

func NewHandler() *Handler {
    handler := &Handler{
        serverService: server.NewServerService(
            repository.NewInMemStorage(),
        ),
        Mux: http.NewServeMux(),
    }

    handler.Mux.HandleFunc("/update/", handler.Update)

    return handler
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
        if errors.Is(err, domain.ErrInvalidMetricType) {
            w.WriteHeader(http.StatusBadRequest)

            return
        } else if errors.Is(err, domain.ErrInvalidMetricName) {
            w.WriteHeader(http.StatusNotFound)

            return
        } else if errors.Is(err, domain.ErrInvalidMetricValue) {
            w.WriteHeader(http.StatusBadRequest)

            return
        }
    }

    w.WriteHeader(http.StatusOK)
}
