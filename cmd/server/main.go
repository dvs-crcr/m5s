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
)

var logger = internalLogger.NewLogger()

func main() {
    config := NewDefaultConfig()
    if err := config.parseVariables(); err != nil {
        log.Fatal(err)
    }

    if err := internalLogger.SetLogLevel(config.LogLevel); err != nil {
        log.Fatal(err)
    }

    if err := execute(config); err != nil {
        log.Fatal(err)
    }
}

func execute(cfg *Config) error {
    ctx := context.Background()

    logger.Infow(
        "starting server",
        "config", cfg,
    )

    serverStorage, err := selectServerStorage(ctx, cfg)
    if err != nil {
        logger.Fatal(err.Error())
    }

    serverService := server.NewServerService(
        serverStorage,
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

func selectServerStorage(ctx context.Context, cfg *Config) (server.Storage, error) {
    var serverStorage server.Storage

    switch {
    case cfg.DatabaseDSN != "":
        var err error

        serverStorage, err = databasestorage.NewDBStorage(
            ctx,
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
        serverStorage = memorystorage.NewMemStorage()
    }

    return serverStorage, nil
}
