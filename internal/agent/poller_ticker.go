package agent

import (
    "time"

    "m5s/domain"
    "m5s/internal/models"
)

func (as *Service) StartPollTicker() {
    as.logger.Info(
        "Starting poll ticker",
        "pollInterval", as.config.pollInterval,
    )

    stat := models.NewStatistics()

    ticker := time.NewTicker(as.config.pollInterval)

    for range ticker.C {
        stat.Refresh()

        for name, value := range stat.CurrentValues {
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
