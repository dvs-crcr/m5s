package handlers

import "net/http"

func (h *Handler) GetMetricsList(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    metricsList := h.serverService.GetMetricsList()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricsList))
}
