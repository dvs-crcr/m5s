package rest

import (
    "errors"
    "net/http"
    "strings"

    "m5s/internal/application"
)

func (s *MetricsServer) Update(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 5 {
        w.WriteHeader(http.StatusNotFound)

        return
    }

    if err := s.storage.Update(segments[2], segments[3], segments[4]); err != nil {
        if errors.Is(err, application.ErrInvalidMetricType) {
            w.WriteHeader(http.StatusBadRequest)

            return
        } else if errors.Is(err, application.ErrInvalidMetricName) {
            w.WriteHeader(http.StatusNotFound)

            return
        } else if errors.Is(err, application.ErrInvalidMetricValue) {
            w.WriteHeader(http.StatusBadRequest)

            return
        }
    }

    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
}
