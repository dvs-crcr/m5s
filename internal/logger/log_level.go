package logger

import (
    "errors"
    "strings"
)

var ErrWrongLogLevel = errors.New(
    "wrong log level (available: fatal, error, warn, info, debug)",
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
    strLevel = strings.ToLower(strLevel)
    switch strLevel {
    case "fatal":
        return FATAL, nil
    case "error":
        return ERROR, nil
    case "warn":
        return WARN, nil
    case "info":
        return INFO, nil
    case "debug":
        return DEBUG, nil
    default:
        return WARN, ErrWrongLogLevel
    }
}
