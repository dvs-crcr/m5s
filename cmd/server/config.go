package main

import (
    "flag"
    "os"
)

type Config struct {
    LogLevel string
    Addr     string
}

var (
    DefaultLogLevel = "error"
    DefaultAddress  = "localhost:8080"
)

func NewDefaultConfig() *Config {
    return &Config{
        LogLevel: DefaultLogLevel,
        Addr:     DefaultAddress,
    }
}

func (c *Config) parseVariables() {
    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.StringVar(
        &c.LogLevel, "l", c.LogLevel, "logging level",
    )

    flag.Parse()

    if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
        c.Addr = addrEnv
    }

    if logLevelEnv := os.Getenv("LOG_LEVEL"); logLevelEnv != "" {
        c.LogLevel = logLevelEnv
    }
}
