package fileStorage

import (
    "context"
    "time"
)

// startStoreTicker uses for starting Auto-Backup ticker
func (ifs *FileStorage) startStoreTicker(ctx context.Context) {
    if ifs.storeInterval == 0 {
        return
    }

    ifs.logger.Info(
        "starting store ticker",
        "storeInterval", ifs.storeInterval,
    )

    ticker := time.NewTicker(ifs.storeInterval)

    for range ticker.C {
        if err := ifs.backupMetrics(ctx); err != nil {
            ifs.logger.Error(
                "backup metrics",
                "error", err,
            )
        }
    }
}
