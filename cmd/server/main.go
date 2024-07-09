package main

import (
    "log"
    "net/http"
    "net/url"

    "github.com/go-chi/chi/v5"

    "m5s/internal/api"
)

func main() {
    parseFlags()

    if err := execute(); err != nil {
        log.Fatal(err)
    }
}

func execute() error {
    apiHandler := api.NewHandler()

    r := chi.NewRouter()

    r.Get("/", apiHandler.GetMetricsList)
    r.Post("/update/{metricType}/{metricName}/{metricValue}", apiHandler.Update)
    r.Get("/value/{metricType}/{metricName}", apiHandler.GetMetric)

    if _, err := url.Parse(flagRunAddr); err != nil {
        return err
    }

    log.Printf("Running server on %s", flagRunAddr)
    return http.ListenAndServe(flagRunAddr, r)
}
