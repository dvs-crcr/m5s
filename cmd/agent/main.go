package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/agent"
)

func main() {
    parseFlags()

    if err := execute(); err != nil {
        log.Fatal(err)
    }
}

func execute() error {
    as := agent.NewAgentService(
        time.Duration(flagRunPollInterval)*time.Second,
        time.Duration(flagRunReportInterval)*time.Second,
        flagRunAddr,
    )

    go as.StartPoller()
    go as.StartReporter()

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalChannel

    return nil
}
