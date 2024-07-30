package server

import "time"

func (ss *Service) StartStoreTicker() {
    ss.logger.Info(
        "Starting store ticker",
        "storeInterval", ss.config.storeInterval,
        "storagePath", ss.config.fileStoragePath,
    )

    ticker := time.NewTicker(ss.config.storeInterval)

    for range ticker.C {

        // TODO: implement
    }
}
