package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")

    if r.Method != http.MethodPost {
        api.HandleErrors(api.ErrInvalidMethod, w)

        return
    }

    segments := strings.Split(r.URL.Path, "/")

    if len(segments) < 5 {
        api.HandleErrors(api.ErrInvalidSegmentsCount, w)

        return
    }

    if err := h.serverService.Update(segments[2], segments[3], segments[4]); err != nil {
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
