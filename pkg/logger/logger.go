package logger

import (
    "go.uber.org/zap"
)

type Logger struct {
    zap.SugaredLogger
}

var instance Logger
var globalLogLevel LogLevel = ERROR

// NewLogger init new logger instance with LogLevel and appName
func NewLogger(appName string, logLevelStr string) (*Logger, error) {
    var err error

    zapLogger, err := zap.NewDevelopment()
    if err != nil {
        panic(err)
    }
    defer zapLogger.Sync()

    zapSugarLogger := *zapLogger.Sugar()

    zapSugarLogger = *zapSugarLogger.With("cmd", appName)

    if err := SetLogLevel(logLevelStr); err != nil {
        return nil, err
    }

    zapLevel := zapSugarLogger.Level()
    zapLevel.Set(globalLogLevel.String())

    instance = Logger{
        zapSugarLogger,
    }

    return &instance, nil
}

// GetLogger uses to get Logger instance
func GetLogger() *Logger {
    return &instance
}

func (l *Logger) With(args ...interface{}) *Logger {
    return &Logger{
        *instance.SugaredLogger.With(args...),
    }
}

// SetLogLevel uses to define log level
func SetLogLevel(logLevelStr string) error {
    var err error

    if globalLogLevel, err = parseLogLevel(logLevelStr); err != nil {
        return err
    }

    return nil
}
