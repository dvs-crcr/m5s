package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/internal/agent"
    memoryStorage "m5s/internal/storage/memory_storage"

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
    ctx := context.Background()

    loggerProvider := providers.NewZapProvider()
    logger := internalLogger.NewLogger(
        internalLogger.WithProvider(loggerProvider),
        internalLogger.WithLogLevel(cfg.LogLevel),
    )

    logger.Info(
        "starting agent",
        "config", cfg,
    )

    agentStorage := memoryStorage.NewMemStorage()

    agentService := agent.NewAgentService(
        agentStorage,
        agent.WithLogger(logger),
        agent.WithAddress(cfg.Addr),
        agent.WithPollInterval(time.Duration(cfg.PollInterval)*time.Second),
        agent.WithReportInterval(time.Duration(cfg.ReportInterval)*time.Second),
    )

    go agentService.StartPollTicker(ctx)
    go agentService.StartReportTicker(ctx)

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalChannel
}
