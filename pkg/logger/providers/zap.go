package providers

import (
    "log"

    "go.uber.org/zap"

    "m5s/pkg/logger"
)

type ZapProvider struct {
    zapLogger zap.SugaredLogger
}

var _ logger.Logger = (*ZapProvider)(nil)

func NewZapProvider() *ZapProvider {
    zapLogger, err := zap.NewDevelopment()
    if err != nil {
        panic(err)
    }
    defer zapLogger.Sync()

    zapSugarLogger := *zapLogger.Sugar()

    return &ZapProvider{
        zapLogger: zapSugarLogger,
    }
}

func (p *ZapProvider) SetLogLevel(level logger.LogLevel) {
    zapLevel := p.zapLogger.Level()
    if err := zapLevel.Set(level.String()); err != nil {
        log.Println(err)
    }
}

func (p *ZapProvider) Fatal(msg string, keysAndValues ...interface{}) {
    p.zapLogger.Fatalw(msg, keysAndValues...)
}

func (p *ZapProvider) Error(msg string, keysAndValues ...interface{}) {
    p.zapLogger.Errorw(msg, keysAndValues...)
}

func (p *ZapProvider) Warn(msg string, keysAndValues ...interface{}) {
    p.zapLogger.Warnw(msg, keysAndValues...)
}

func (p *ZapProvider) Info(msg string, keysAndValues ...interface{}) {
    p.zapLogger.Infow(msg, keysAndValues...)
}

func (p *ZapProvider) Debug(msg string, keysAndValues ...interface{}) {
    p.zapLogger.Debugw(msg, keysAndValues...)
}
