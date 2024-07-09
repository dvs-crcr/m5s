package main

import (
    "flag"
    "os"
)

type Config struct {
    Addr string
}

func NewDefaultConfig() *Config {
    return &Config{
        Addr: "localhost:8080",
    }
}

func (c *Config) parseVariables() {
    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.Parse()

    if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
        c.Addr = addrEnv
    }
}
