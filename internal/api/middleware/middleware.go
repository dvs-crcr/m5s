package middleware

import (
    "m5s/pkg/logger"
)

type Middleware struct {
    logger logger.Logger
}

func NewMiddleware(loggerInstance logger.Logger) *Middleware {
    return &Middleware{
        logger: loggerInstance,
    }
}
