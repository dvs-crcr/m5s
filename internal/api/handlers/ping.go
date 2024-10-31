package handlers

import (
    "net/http"

    "m5s/internal/api"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if err := h.serverService.PingDB(ctx); err != nil {
        api.HandleErrors(api.ErrInternal, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
