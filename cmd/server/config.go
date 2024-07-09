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
        Addr: "localhost:80880",
    }
}

func (c *Config) parseVariables() {
    c.parseAddress()
}

func (c *Config) parseAddress() {
    if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
        return
    }

    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.Parse()
}
