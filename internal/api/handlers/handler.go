package handlers

import (
    "m5s/internal/repository"
    "m5s/internal/server"
    "m5s/pkg/logger"
)

type Handler struct {
    serverService *server.Service
    logger        logger.Logger
}

func NewHandler(loggerInstance logger.Logger) *Handler {
    return &Handler{
        serverService: server.NewServerService(
            repository.NewInMemStorage(),
        ),
        logger: loggerInstance,
    }
}
