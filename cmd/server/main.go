package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api/handlers"
    "m5s/internal/api/middleware"
    "m5s/internal/server"
    databasestorage "m5s/internal/storage/database_storage"
    filestorage "m5s/internal/storage/file_storage"
    memorystorage "m5s/internal/storage/memory_storage"
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

    serverStorage, err := selectServerStorage(ctx, logger, cfg)
    if err != nil {
        logger.Fatal(
            "select server storage",
            "error", err,
        )
    }

    serverService := server.NewServerService(
        serverStorage,
        server.WithLogger(logger),
    )

    apiHandler := handlers.NewHandler(
        serverService,
        handlers.WithLogger(logger),
    )

    apiMiddleware := middleware.NewMiddleware(logger)

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

    logger.Info(
        "starting server",
        "config", cfg,
    )
    return http.ListenAndServe(cfg.Addr, r)
}

func selectServerStorage(ctx context.Context, logger internalLogger.Logger, cfg *Config) (server.Storage, error) {
    var serverStorage server.Storage

    switch {
    case cfg.DatabaseDSN != "":
        var err error

        serverStorage, err = databasestorage.NewDBStorage(
            ctx,
            logger,
            cfg.DatabaseDSN,
            cfg.MigrationsPath,
            cfg.MigrationsVersion,
        )
        if err != nil {
            return nil, fmt.Errorf(
                "unable to create new db storage instance: %w", err,
            )
        }
    case cfg.FileStoragePath != "":
        var err error

        serverStorage, err = filestorage.NewFileStorage(
            ctx,
            logger,
            cfg.FileStoragePath,
            time.Duration(cfg.StoreInterval)*time.Second,
            cfg.Restore,
        )
        if err != nil {
            return nil, fmt.Errorf(
                "unable to create new file storage instance: %w", err,
            )
        }
    default:
        serverStorage = memorystorage.NewMemStorage(logger)
    }

    return serverStorage, nil
}
