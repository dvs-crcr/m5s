package api

import (
    "errors"
    "net/http"
    "slices"

    "m5s/domain"
)

var (
    ErrInvalidJSONContentType = errors.New(
        "invalid request Content-Type, \"application/json\" expected",
    )
    ErrInvalidJSONStruct    = errors.New("invalid JSON structure")
    ErrInvalidSegmentsCount = errors.New("invalid segments count")
    ErrInvalidMethod        = errors.New("invalid method")
    ErrInternal             = errors.New("internal server error")
)

func HandleErrors(err error, w http.ResponseWriter) {
    switch {
    case errors.Is(err, ErrInternal):
        w.WriteHeader(http.StatusInternalServerError)
    case errors.Is(err, ErrInvalidMethod):
        w.WriteHeader(http.StatusMethodNotAllowed)
    case errors.Is(err, ErrInvalidSegmentsCount):
        w.WriteHeader(http.StatusNotFound)
    case errors.Is(err, ErrInvalidJSONContentType):
        w.WriteHeader(http.StatusBadRequest)
    case errors.Is(err, domain.ErrInvalidMetricType):
        w.WriteHeader(http.StatusBadRequest)
    case errors.Is(err, domain.ErrInvalidMetricName):
        w.WriteHeader(http.StatusNotFound)
    case errors.Is(err, domain.ErrInvalidMetricValue):
        w.WriteHeader(http.StatusBadRequest)
    case errors.Is(err, domain.ErrNoSuchMetric):
        w.WriteHeader(http.StatusNotFound)
    }
}

func CheckContentType(contentType string, availableTypes ...string) bool {
    return slices.Contains(availableTypes, contentType)
}
