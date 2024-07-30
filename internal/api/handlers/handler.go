package handlers

import (
    "m5s/internal/server"
    "m5s/pkg/logger"
)

type Handler struct {
    serverService *server.Service
    logger        logger.Logger
}

func NewHandler(loggerInstance logger.Logger, service *server.Service) *Handler {
    return &Handler{
        serverService: service,
        logger:        loggerInstance,
    }
}
