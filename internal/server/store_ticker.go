package server

import "time"

func (ss *Service) StartStoreTicker() {
    if ss.storage != nil && ss.config.storeInterval == 0 {
        return
    }

    ss.logger.Info(
        "Starting store ticker",
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
    }
}
