package handlers

import (
    "m5s/internal/server"
    "m5s/pkg/logger"
)

type Handler struct {
    serverService *server.Service
    logger        logger.Logger
}

type Option func(handler *Handler)

func NewHandler(serverService *server.Service, options ...Option) *Handler {
    handler := &Handler{
        serverService: serverService,
    }

    for _, opt := range options {
        opt(handler)
    }

    return handler
}

func WithLogger(logger logger.Logger) Option {
    return func(handler *Handler) {
        handler.logger = logger
    }
}
