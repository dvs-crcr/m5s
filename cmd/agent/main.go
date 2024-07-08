package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "m5s/agent"
)

func main() {
    pollInterval := 2 * time.Second
    reportInterval := 10 * time.Second
    port := 8080
    host := "localhost"

    as := agent.NewAgentService(
        pollInterval, reportInterval, fmt.Sprintf("%s:%d", host, port),
    )

    go as.StartPoller()
    go as.StartReporter()

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalChannel
}
