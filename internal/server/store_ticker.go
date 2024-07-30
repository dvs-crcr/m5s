package server

import "time"

func (ss *Service) StartStoreTicker() {
    ss.logger.Info(
        "Starting store ticker", "duration", ss.storeInterval,
    )

    ticker := time.NewTicker(ss.storeInterval)

    for range ticker.C {
        // TODO: implement
    }
}
