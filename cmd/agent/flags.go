package main

import (
    "flag"
)

var (
    flagRunAddr           string
    flagRunPollInterval   int64
    flagRunReportInterval int64
)

func parseFlags() {
    flag.StringVar(
        &flagRunAddr, "a", "localhost:8080", "server endpoint address",
    )
    flag.Int64Var(
        &flagRunPollInterval, "p", 2, "poll interval (sec)",
    )
    flag.Int64Var(
        &flagRunReportInterval, "r", 10, "report interval (sec)",
    )

    flag.Parse()
}
