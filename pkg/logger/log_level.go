package logger

import (
    "errors"
)

var ErrWrongLogLevel = errors.New(
    "wrong log level (available levels: fatal, error, warn, info, debug)",
)

type LogLevel int

const (
    FATAL LogLevel = iota
    ERROR
    WARN
    INFO
    DEBUG
)

func (ll LogLevel) String() string {
    return [...]string{"fatal", "error", "warn", "info", "debug"}[ll]
}

func parseLogLevel(strLevel string) (LogLevel, error) {
    switch strLevel {
    case "fatal", "FATAL":
        return FATAL, nil
    case "error", "ERROR":
        return ERROR, nil
    case "warn", "WARN":
        return WARN, nil
    case "info", "INFO":
        return INFO, nil
    case "debug", "DEBUG":
        return DEBUG, nil
    default:
        return WARN, ErrWrongLogLevel
    }
}
