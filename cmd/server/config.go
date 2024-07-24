package main

import (
    "flag"
    "os"
)

type Config struct {
    Addr string
}

var (
    DefaultAddress = "localhost:8080"
)

func NewDefaultConfig() *Config {
    return &Config{
        Addr: DefaultAddress,
    }
}

func (c *Config) parseVariables() {
    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.Parse()

    if addrEnv := os.Getenv("METRICS_SERVER_ADDRESS"); addrEnv != "" {
        c.Addr = addrEnv
    }
}
