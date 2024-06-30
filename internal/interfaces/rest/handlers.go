package rest

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "m5s/internal/application"
)

type MetricsServer struct {
    storage *application.MemStorage
    server  http.Server
}

// NewMetricsServer returns new MetricsServer instance
func NewMetricsServer(host string, port int) *MetricsServer {
    return &MetricsServer{
        storage: application.NewMemStorage(),
        server: http.Server{
            Addr:              fmt.Sprintf("%s:%d", host, port),
            ReadTimeout:       1 * time.Second,
            ReadHeaderTimeout: 2 * time.Second,
            WriteTimeout:      1 * time.Second,
            IdleTimeout:       30 * time.Second,
        },
    }
}

// Start uses to launch server application
func (s *MetricsServer) Start() {
    mux := http.NewServeMux()

    s.server.Handler = mux

    // Register interfaces
    mux.HandleFunc("/update/", s.Update)

    log.Fatal(s.server.ListenAndServe())
}
