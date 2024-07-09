package main

import (
    "fmt"
    "log"
    "net/http"

    "m5s/internal/api"
)

func main() {
    port := 8080
    host := "localhost"

    apiHandler := api.NewHandler()
    log.Fatal(
        http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), apiHandler.Mux),
    )
}
