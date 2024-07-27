package api

import (
    "errors"
    "net/http"
    "time"

    "m5s/domain"
    "m5s/internal/logger"
    "m5s/internal/repository"
    "m5s/internal/server"
)

type Handler struct {
    serverService *server.Service
    logger        logger.Logger
}

var (
    ErrInvalidJSONContentType = errors.New(
        "invalid request Content-Type, \"application/json\" expected",
    )
    ErrInvalidJSONStruct    = errors.New("invalid JSON structure")
    ErrInvalidSegmentsCount = errors.New("invalid segments count")
    ErrInvalidMethod        = errors.New("invalid method")
)

func NewHandler(loggerInstance logger.Logger) *Handler {
    return &Handler{
        serverService: server.NewServerService(
            repository.NewInMemStorage(),
        ),
        logger: loggerInstance,
    }
}

func handleErrors(err error, w http.ResponseWriter) {
    switch {
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

func (h *Handler) WithLogger(handler http.Handler) http.Handler {
    logFn := func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rd := &responseData{
            status: 0,
            size:   0,
        }

        lw := loggingResponseWriter{
            ResponseWriter: w,
            responseData:   rd,
        }

        handler.ServeHTTP(&lw, r)

        duration := time.Since(start)

        h.logger.Info(
            "Request",
            "uri", r.RequestURI,
            "method", r.Method,
            "status", rd.status,
            "duration", duration,
            "size", rd.size,
        )
    }

    return http.HandlerFunc(logFn)
}
