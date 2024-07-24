package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Addr           string
	PollInterval   int64
	ReportInterval int64
}

func NewDefaultConfig() *Config {
	return &Config{
		Addr:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
}

// env -> arg -> default

func (c *Config) parseVariables() error {
	var err error

	flag.StringVar(
		&c.Addr, "a", c.Addr, "server endpoint address",
	)
	flag.Int64Var(
		&c.ReportInterval, "r", c.ReportInterval, "report interval (sec)",
	)
	flag.Int64Var(
		&c.PollInterval, "p", c.PollInterval, "poll interval (sec)",
	)

	flag.Parse()

	if addrEnv := os.Getenv("METRICS_AGENT_ADDRESS"); addrEnv != "" {
		c.Addr = addrEnv
	}

	if reportEnv := os.Getenv("REPORT_INTERVAL"); reportEnv != "" {
		if c.ReportInterval, err = strconv.ParseInt(
			reportEnv,
			10,
			64,
		); err != nil {
			return err
		}
	}

	if pollEnv := os.Getenv("POLL_INTERVAL"); pollEnv != "" {
		if c.PollInterval, err = strconv.ParseInt(
			pollEnv,
			10,
			64,
		); err != nil {
			return err
		}
	}

	return nil
}
