package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/internal/agent"
    memorystorage "m5s/internal/storage/memory_storage"
    internalLogger "m5s/pkg/logger"
)

var logger = internalLogger.NewLogger()

func main() {
    config := NewDefaultConfig()
    if err := config.parseVariables(); err != nil {
        log.Fatal(err)
    }

    if err := internalLogger.SetLogLevel(config.LogLevel); err != nil {
        log.Fatal(err)
    }

    execute(config)
}

func execute(cfg *Config) {
    ctx := context.Background()

    logger.Infow(
        "starting agent",
        "config", cfg,
    )

    agentStorage := memorystorage.NewMemStorage()

    agentService := agent.NewAgentService(
        agentStorage,
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
