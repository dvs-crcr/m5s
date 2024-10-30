package handlers

import (
    "m5s/internal/server"
    internalLogger "m5s/pkg/logger"
)

var logger = internalLogger.GetLogger()

type Handler struct {
    serverService *server.Service
}

type Option func(handler *Handler)

func NewHandler(serverService *server.Service, options ...Option) *Handler {
    logger = logger.With(
        "package", "handlers",
    )

    handler := &Handler{
        serverService: serverService,
    }

    for _, opt := range options {
        opt(handler)
    }

    return handler
}
