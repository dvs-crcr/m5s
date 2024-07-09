package main

import (
    "log"
    "net/url"
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
    if _, err := url.Parse(flagRunAddr); err != nil {
        return err
    }

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
