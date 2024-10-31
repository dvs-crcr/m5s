package handlers

import (
    "encoding/json"
    "io"
    "net/http"

    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if !api.CheckContentType(
        r.Header.Get("Content-Type"),
        "application/json",
    ) {
        api.HandleErrors(api.ErrInvalidJSONContentType, w)

        return
    }

    metric := &models.Metrics{}

    bodyData, err := io.ReadAll(r.Body)
    if err != nil {
        logger.Errorw(err.Error())

        return
    }

    if err := json.Unmarshal(bodyData, &metric); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }

    if err := h.serverService.Update(
        ctx,
        metric.MType,
        metric.ID,
        metric.String(),
    ); err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
