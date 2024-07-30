package agent

import (
    "time"

    "m5s/domain"
)

func (as *Service) StartPollTicker() {
    as.logger.Info("Starting poll ticker", "duration", as.pollInterval)

    ticker := time.NewTicker(as.pollInterval)

    for range ticker.C {
        as.stat.Refresh()

        for name, value := range as.stat.CurrentValues {
            metric := domain.NewGauge(name, value)
            if err := as.repo.Update(metric); err != nil {
                as.logger.Error("update gauge", "error", err)
            }
        }

        pollCountMetric := domain.NewCounter("PollCount", 1)
        if err := as.repo.Update(pollCountMetric); err != nil {
            as.logger.Error("update counter", "error", err)
        }
    }
}
