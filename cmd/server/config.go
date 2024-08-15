package main

import (
    "flag"
    "os"
    "strconv"
)

var (
    DefaultLogLevel                      = "info"
    DefaultAddress                       = "localhost:8080"
    DefaultStoreInterval           int64 = 300
    DefaultFileStoragePath               = "tmp/file_storage"
    DefaultDatabaseDSN                   = ""
    DefaultMigrationsPath                = "migrations"
    DefaultMigrationsSchemaVersion       = ""
)

type Config struct {
    LogLevel          string
    Addr              string
    StoreInterval     int64
    FileStoragePath   string
    MigrationsPath    string
    MigrationsVersion string
    DatabaseDSN       string
    Restore           bool
}

func NewDefaultConfig() *Config {
    return &Config{
        LogLevel:          DefaultLogLevel,
        Addr:              DefaultAddress,
        StoreInterval:     DefaultStoreInterval,
        FileStoragePath:   DefaultFileStoragePath,
        MigrationsPath:    DefaultMigrationsPath,
        MigrationsVersion: DefaultMigrationsSchemaVersion,
        Restore:           true,
        DatabaseDSN:       DefaultDatabaseDSN,
    }
}

func (c *Config) parseVariables() error {
    var err error

    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.StringVar(
        &c.LogLevel, "l", c.LogLevel, "logging level",
    )

    flag.Int64Var(
        &c.StoreInterval, "i", c.StoreInterval, "store interval (sec)",
    )

    flag.StringVar(
        &c.FileStoragePath, "f", c.FileStoragePath, "file storage path",
    )

    flag.BoolVar(
        &c.Restore, "r", c.Restore, "restore file",
    )

    flag.StringVar(
        &c.DatabaseDSN, "d", c.DatabaseDSN, "database DSN",
    )

    flag.StringVar(
        &c.MigrationsPath, "m", c.MigrationsPath, "migrations folder",
    )

    flag.Parse()

    if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
        c.Addr = addrEnv
    }

    if logLevelEnv := os.Getenv("LOG_LEVEL"); logLevelEnv != "" {
        c.LogLevel = logLevelEnv
    }

    if storeEnv := os.Getenv("STORE_INTERVAL"); storeEnv != "" {
        if c.StoreInterval, err = strconv.ParseInt(
            storeEnv,
            10,
            64,
        ); err != nil {
            return err
        }
    }

    if fileStorageEnv := os.Getenv("FILE_STORAGE_PATH"); fileStorageEnv != "" {
        c.FileStoragePath = fileStorageEnv
    }

    if restoreEnv := os.Getenv("RESTORE"); restoreEnv != "" {
        if restoreEnv == "false" || restoreEnv == "FALSE" {
            c.Restore = false
        }
    }

    if dsnEnv := os.Getenv("DATABASE_DSN"); dsnEnv != "" {
        c.DatabaseDSN = dsnEnv
    }

    if migrationsPathEnv := os.Getenv("MIGRATIONS_PATH"); migrationsPathEnv != "" {
        c.MigrationsPath = migrationsPathEnv
    }

    if migrationsVersionEnv := os.Getenv("MIGRATIONS_VERSION"); migrationsVersionEnv != "" {
        c.MigrationsVersion = migrationsVersionEnv
    }

    return nil
}
