package api

import (
    "encoding/json"
    "net/http"
    "strings"

    "m5s/internal/models"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodPost {
        handleErrors(ErrInvalidMethod, w)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 5 {
        handleErrors(ErrInvalidSegmentsCount, w)

        return
    }

    if err := h.serverService.Update(segments[2], segments[3], segments[4]); err != nil {
        handleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
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

    if err := h.serverService.Update(
        metric.MType,
        metric.ID,
        metric.String(),
    ); err != nil {
        handleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
