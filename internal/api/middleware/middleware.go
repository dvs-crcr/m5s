package middleware

import internalLogger "m5s/pkg/logger"

var logger = internalLogger.GetLogger()

type Middleware struct {
}

func NewMiddleware() *Middleware {
    logger = logger.With(
        "package", "middleware",
    )

    return &Middleware{}
}
