package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api/handlers"
    "m5s/internal/api/middleware"
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
    ctx := context.Background()

    // Init logger
    loggerProvider := providers.NewZapProvider()
    logger := internalLogger.NewLogger(
        internalLogger.WithProvider(loggerProvider),
        internalLogger.WithLogLevel(cfg.LogLevel),
    )

    serverService := server.NewServerService(
        server.WithLogger(logger),
        server.WithStorage(
            ctx,
            cfg.Restore,
            cfg.FileStoragePath,
            time.Duration(cfg.StoreInterval)*time.Second,
            cfg.DatabaseDSN,
            cfg.MigrationsPath,
            cfg.MigrationsVersion,
        ),
    )

    apiHandler := handlers.NewHandler(
        serverService,
        handlers.WithLogger(logger),
    )

    apiMiddleware := middleware.NewMiddleware(logger)

    r := chi.NewRouter()

    // Middlewares
    r.Use(apiMiddleware.WithRequestLogger)
    r.Use(middleware.WithCompression)

    // Routes
    r.Route("/", func(r chi.Router) {
        r.Get("/", apiHandler.GetMetricsList)

        r.Get("/ping", apiHandler.Ping)

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
        "starting server",
        "config", cfg,
    )
    return http.ListenAndServe(cfg.Addr, r)
}
