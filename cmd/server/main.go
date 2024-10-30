package main

import (
    "context"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api/handlers"
    "m5s/internal/api/middleware"
    "m5s/internal/server"
    internalLogger "m5s/pkg/logger"
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

    logger, err := internalLogger.NewLogger("server", cfg.LogLevel)
    if err != nil {
        return err
    }

    logger.Infow(
        "starting server",
        "config", cfg,
    )

    serverService := server.NewServerService(
        ctx,
        &server.Config{
            Addr:              cfg.Addr,
            StoreInterval:     cfg.StoreInterval,
            FileStoragePath:   cfg.FileStoragePath,
            MigrationsPath:    cfg.MigrationsPath,
            MigrationsVersion: cfg.MigrationsVersion,
            DatabaseDSN:       cfg.DatabaseDSN,
            Restore:           cfg.Restore,
        },
    )

    apiHandler := handlers.NewHandler(
        serverService,
    )

    apiMiddleware := middleware.NewMiddleware()

    r := chi.NewRouter()
    r.Use(apiMiddleware.WithRequestLogger)
    r.Use(middleware.WithCompression)

    r.Route("/", func(r chi.Router) {
        r.Get("/", apiHandler.GetMetricsList)

        r.Get("/ping", apiHandler.Ping)

        r.Post("/updates/", apiHandler.UpdateBatch)

        r.Route("/update", func(r chi.Router) {
            r.Post("/", apiHandler.UpdateJSON)
            r.Post("/{metricType}/{metricName}/{metricValue}", apiHandler.Update)
        })

        r.Route("/value", func(r chi.Router) {
            r.Post("/", apiHandler.GetMetricJSON)
            r.Get("/{metricType}/{metricName}", apiHandler.GetMetric)
        })
    })

    return http.ListenAndServe(cfg.Addr, r)
}
