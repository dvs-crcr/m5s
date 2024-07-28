package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api/handlers"
    "m5s/internal/api/middleware"
    internalLogger "m5s/pkg/logger"
    "m5s/pkg/logger/providers"
)

func main() {
    config := NewDefaultConfig()
    config.parseVariables()

    if err := execute(config); err != nil {
        log.Fatal(err)
    }
}

func execute(cfg *Config) error {
    loggerProvider := providers.NewZapProvider()
    logger := internalLogger.NewLogger(
        internalLogger.WithProvider(loggerProvider),
        internalLogger.WithLogLevel(cfg.LogLevel),
    )

    apiHandler := handlers.NewHandler(logger)
    apiMiddleware := middleware.NewMiddleware(logger)

    r := chi.NewRouter()

    r.Use(apiMiddleware.WithLogger)
    r.Use(apiMiddleware.WithCompression)

    r.Get("/", apiHandler.GetMetricsList)
    r.Route("/update", func(r chi.Router) {
        r.Post("/", apiHandler.UpdateJSON)
        r.Post("/{metricType}/{metricName}/{metricValue}", apiHandler.Update)
    })
    r.Route("/value", func(r chi.Router) {
        r.Post("/", apiHandler.GetMetricJSON)
        r.Get("/{metricType}/{metricName}", apiHandler.GetMetric)
    })

    logger.Info(
        "Starting server",
        "addr", cfg.Addr,
    )
    return http.ListenAndServe(cfg.Addr, r)
}
