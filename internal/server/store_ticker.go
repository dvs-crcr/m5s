package server

import (
    "time"

    "m5s/internal/storage"
)

// StartStoreTicker uses for starting Auto-Backup ticker
func (ss *Service) StartStoreTicker() {
    if ss.storage == nil {
        return
    }

    if ss.storage.MyType() != storage.TypeFile ||
        ss.config.storeInterval == 0 {
        return
    }

    ss.logger.Info(
        "starting store ticker",
        "storeInterval", ss.config.storeInterval,
    )

    ticker := time.NewTicker(ss.config.storeInterval)

    for range ticker.C {
        if err := ss.BackupMetrics(); err != nil {
            ss.logger.Error(
                "backup metrics",
                "error", err,
            )
        }

        ss.logger.Info(
            "metrics have been successfully backed up to a file",
        )
    }
}
