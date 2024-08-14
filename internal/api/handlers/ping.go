package handlers

import (
    "net/http"

    "m5s/internal/api"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
    if err := h.serverService.PingDB(r.Context()); err != nil {
        api.HandleErrors(api.ErrInternal, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
