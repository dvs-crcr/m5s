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

var (
	DefaultAddress              = "localhost:8080"
	DefaultPollInterval   int64 = 2
	DefaultReportInterval int64 = 10
)

func NewDefaultConfig() *Config {
	return &Config{
		Addr:           DefaultAddress,
		PollInterval:   DefaultPollInterval,
		ReportInterval: DefaultReportInterval,
	}
}

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

	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
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
