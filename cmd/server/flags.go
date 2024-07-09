package main

import "flag"

var flagRunAddr string

func parseFlags() {
    flag.StringVar(
        &flagRunAddr, "a", "localhost:8080", "server endpoint address",
    )
    flag.Parse()
}
