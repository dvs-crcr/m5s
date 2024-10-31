package logger

import (
    "go.uber.org/zap"
)

type Logger struct {
    zap.SugaredLogger
}

var isInit bool
var instance *Logger
var level zap.AtomicLevel
var globalLogLevel LogLevel = INFO

// NewLogger init new logger instance or return existed
func NewLogger() *Logger {
    if isInit {
        return instance
    }

    level = zap.NewAtomicLevel()
    level.UnmarshalText([]byte(globalLogLevel.String()))

    encoderConfig := zap.NewProductionEncoderConfig()

    zapConfig := zap.NewProductionConfig()
    zapConfig.Development = true
    zapConfig.EncoderConfig = encoderConfig
    zapConfig.Level = level
    zapConfig.Encoding = "console" // "json"
    zapConfig.OutputPaths = []string{"stdout"}
    zapConfig.ErrorOutputPaths = []string{"stderr"}

    zapLogger, _ := zapConfig.Build()
    defer zapLogger.Sync()

    instance = &Logger{
        *zapLogger.Sugar(),
    }

    isInit = true

    instance.Infow("init logger", "log_level", instance.SugaredLogger.Level().String())

    return instance
}

// With wrapper
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

    instance.Infow("change log level", "log_level", globalLogLevel)

    if err = level.UnmarshalText([]byte(globalLogLevel.String())); err != nil {
        return err
    }

    return nil
}
