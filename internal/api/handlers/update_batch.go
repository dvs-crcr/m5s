package handlers

import (
    "encoding/json"
    "io"
    "net/http"

    "m5s/internal/api"
    "m5s/internal/models"
)

func (h *Handler) UpdateBatch(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if !api.CheckContentType(
        r.Header.Get("Content-Type"),
        "application/json",
    ) {
        api.HandleErrors(api.ErrInvalidJSONContentType, w)

        return
    }

    metrics := make([]*models.Metrics, 0)

    bodyData, err := io.ReadAll(r.Body)
    if err != nil {
        logger.Errorw(err.Error())

        return
    }

    if err := json.Unmarshal(bodyData, &metrics); err != nil {
        api.HandleErrors(api.ErrInvalidJSONStruct, w)

        return
    }

    if err := h.serverService.UpdateBatch(ctx, metrics); err != nil {
        api.HandleErrors(err, w)

        return
    }

    w.WriteHeader(http.StatusOK)
}
