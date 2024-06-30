package main

import (
    "m5s/internal/interfaces/rest"
)

func main() {
    port := 8080
    host := "localhost"

    ms := rest.NewMetricsServer(host, port)

    ms.Start()
}
