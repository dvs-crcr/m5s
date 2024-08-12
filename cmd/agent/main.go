package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/internal/agent"
    "m5s/internal/repository"

    internalLogger "m5s/pkg/logger"
    "m5s/pkg/logger/providers"
)

func main() {
    config := NewDefaultConfig()
    if err := config.parseVariables(); err != nil {
        log.Fatal(err)
    }

    execute(config)
}

func execute(cfg *Config) {
    loggerProvider := providers.NewZapProvider()
    logger := internalLogger.NewLogger(
        internalLogger.WithProvider(loggerProvider),
        internalLogger.WithLogLevel(cfg.LogLevel),
    )

    logger.Info(
        "Starting agent",
        "config", cfg,
    )

    serverRepository := repository.NewInMemStorage()

    agentService := agent.NewAgentService(
        serverRepository,
        agent.WithLogger(logger),
        agent.WithAddress(cfg.Addr),
        agent.WithPollInterval(time.Duration(cfg.PollInterval)*time.Second),
        agent.WithReportInterval(time.Duration(cfg.ReportInterval)*time.Second),
    )

    go agentService.StartPollTicker()
    go agentService.StartReportTicker()

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalChannel
}
