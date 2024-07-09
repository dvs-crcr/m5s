package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/agent"
)

func main() {
    config := NewDefaultConfig()
    config.parseVariables()

    if err := execute(config); err != nil {
        log.Fatal(err)
    }
}

func execute(cfg *Config) error {
    fmt.Printf("%+v\n", cfg)
    as := agent.NewAgentService(
        time.Duration(cfg.PollInterval)*time.Second,
        time.Duration(cfg.ReportInterval)*time.Second,
        cfg.Addr,
    )

    go as.StartPoller()
    go as.StartReporter()

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalChannel

    return nil
}
