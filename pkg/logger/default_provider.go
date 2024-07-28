package logger

import (
    "log"
    "os"
)

type DefaultProvider struct {
    defaultLogger *log.Logger
    logLevel      LogLevel
}

func NewDefaultProvider() *DefaultProvider {
    dp := &DefaultProvider{
        defaultLogger: &log.Logger{},
    }

    dp.defaultLogger.SetOutput(os.Stdout)
    dp.defaultLogger.SetPrefix("")
    dp.defaultLogger.SetFlags(log.LstdFlags)

    return dp
}

func (p *DefaultProvider) SetLogLevel(level LogLevel) {
    p.logLevel = level
}

func (p *DefaultProvider) Fatal(msg string, keysAndValues ...interface{}) {
    p.defaultLogger.SetPrefix("[FATAL] ")

    if p.logLevel >= FATAL {
        p.defaultLogger.Fatal(msg, keysAndValues)
    }
}

func (p *DefaultProvider) Error(msg string, keysAndValues ...interface{}) {
    p.defaultLogger.SetPrefix("[ERROR] ")

    if p.logLevel >= ERROR {
        p.defaultLogger.Println(msg, keysAndValues)
    }
}

func (p *DefaultProvider) Warn(msg string, keysAndValues ...interface{}) {
    p.defaultLogger.SetPrefix("[WARN] ")

    if p.logLevel >= WARN {
        p.defaultLogger.Println(msg, keysAndValues)
    }
}

func (p *DefaultProvider) Info(msg string, keysAndValues ...interface{}) {
    p.defaultLogger.SetPrefix("[INFO] ")

    if p.logLevel >= INFO {
        p.defaultLogger.Println(msg, keysAndValues)
    }
}

func (p *DefaultProvider) Debug(msg string, keysAndValues ...interface{}) {
    p.defaultLogger.SetPrefix("[DEBUG] ")

    if p.logLevel >= DEBUG {
        p.defaultLogger.Println(msg, keysAndValues)
    }
}
