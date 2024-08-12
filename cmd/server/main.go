package main

import (
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api/handlers"
    "m5s/internal/api/middleware"
    "m5s/internal/repository"
    "m5s/internal/server"
    internalLogger "m5s/pkg/logger"
    "m5s/pkg/logger/providers"
)

func main() {
    config := NewDefaultConfig()
    if err := config.parseVariables(); err != nil {
        log.Fatal(err)
    }

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

    serverRepository := repository.NewInMemStorage()

    serverService := server.NewServerService(
        serverRepository,
        server.WithLogger(logger),
        server.WithStoreInterval(time.Duration(cfg.StoreInterval)*time.Second),
        server.WithStorage(cfg.FileStoragePath),
        server.WithRestore(cfg.Restore),
    )

    apiHandler := handlers.NewHandler(
        serverService,
        handlers.WithLogger(logger),
    )

    apiMiddleware := middleware.NewMiddleware(logger)

    r := chi.NewRouter()

    // Middlewares
    r.Use(apiMiddleware.WithLogger)
    r.Use(apiMiddleware.WithCompression)

    r.Route("/", func(r chi.Router) {
        r.Get("/", apiHandler.GetMetricsList)

        r.Route("/update", func(r chi.Router) {
            r.Post("/", apiHandler.UpdateJSON)
            r.Post("/{metricType}/{metricName}/{metricValue}", apiHandler.Update)
        })

        r.Route("/value", func(r chi.Router) {
            r.Post("/", apiHandler.GetMetricJSON)
            r.Get("/{metricType}/{metricName}", apiHandler.GetMetric)
        })
    })

    logger.Info(
        "Starting server",
        "config", cfg,
    )
    return http.ListenAndServe(cfg.Addr, r)
}
