package main

import (
	"log"
	"m5s/internal/agent"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := NewDefaultConfig()
	if err := config.parseVariables(); err != nil {
		log.Fatal(err)
	}

	execute(config)
}

func execute(cfg *Config) {
	agentService := agent.NewAgentService(
		time.Duration(cfg.PollInterval)*time.Second,
		time.Duration(cfg.ReportInterval)*time.Second,
		cfg.Addr,
	)

	go agentService.StartPoller()
	go agentService.StartReporter()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
}
