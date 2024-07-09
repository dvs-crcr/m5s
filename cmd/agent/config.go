package main

import (
    "flag"
    "os"
)

type Config struct {
    Addr           string
    PollInterval   int64
    ReportInterval int64
}

func NewDefaultConfig() *Config {
    return &Config{
        Addr:           "localhost:80880",
        PollInterval:   2,
        ReportInterval: 10,
    }
}

// env -> arg -> default

func (c *Config) parseVariables() {
    c.parseAddress()
    c.parseReportInterval()
    c.parsePollInterval()
}

func (c *Config) parseAddress() {
    // ADDRESS
    if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
        return
    }

    flag.StringVar(
        &c.Addr, "a", c.Addr, "server endpoint address",
    )

    flag.Parse()
}

func (c *Config) parseReportInterval() {
    // REPORT_INTERVAL
    if addrEnv := os.Getenv("REPORT_INTERVAL"); addrEnv != "" {
        return
    }

    flag.Int64Var(
        &c.ReportInterval, "r", c.ReportInterval, "report interval (sec)",
    )

    flag.Parse()
}

func (c *Config) parsePollInterval() {
    // POLL_INTERVAL
    if addrEnv := os.Getenv("POLL_INTERVAL"); addrEnv != "" {
        return
    }

    flag.Int64Var(
        &c.PollInterval, "p", c.PollInterval, "poll interval (sec)",
    )

    flag.Parse()
}
