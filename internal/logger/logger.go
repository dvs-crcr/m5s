package logger

import (
    "log"
)

type Logger interface {
    Fatal(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Debug(msg string, keysAndValues ...interface{})
    SetLogLevel(level LogLevel)
}

type InternalLogger struct {
    logLevel LogLevel
    provider Logger
}

type Option func(*InternalLogger)

func NewLogger(options ...Option) Logger {
    l := &InternalLogger{
        provider: NewDefaultProvider(),
    }

    for _, opt := range options {
        opt(l)
    }

    return l.provider
}

func WithProvider(loggerProvider Logger) Option {
    return func(l *InternalLogger) {
        l.provider = loggerProvider
        l.provider.SetLogLevel(l.logLevel)
    }
}

func WithLogLevel(logLevel string) Option {
    return func(l *InternalLogger) {
        var err error

        l.logLevel, err = parseLogLevel(logLevel)
        if err != nil {
            log.Println(err)
        }

        l.provider.SetLogLevel(l.logLevel)
    }
}
