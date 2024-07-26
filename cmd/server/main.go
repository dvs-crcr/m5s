package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api"
    internalLogger "m5s/internal/logger"
    "m5s/internal/logger/providers"
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

    apiHandler := api.NewHandler()

    r := chi.NewRouter()

    r.Get("/", apiHandler.GetMetricsList)
    r.Post("/update/{metricType}/{metricName}/{metricValue}", apiHandler.Update)
    r.Get("/value/{metricType}/{metricName}", apiHandler.GetMetric)

    logger.Info(
        "Starting server",
        "addr", cfg.Addr,
    )
    return http.ListenAndServe(cfg.Addr, r)
}
