package handlers

import "net/http"

func (h *Handler) GetMetricsList(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    w.Header().Set("Content-Type", "text/html")

    metricsList := h.serverService.GetMetricsList(ctx)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricsList))
}
