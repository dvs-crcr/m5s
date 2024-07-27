package api

import "net/http"

func (h *Handler) GetMetricsList(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")

    metricsList := h.serverService.GetMetricsList()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(metricsList))
}
