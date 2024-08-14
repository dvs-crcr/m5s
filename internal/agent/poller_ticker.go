package agent

import (
    "math/rand"
    "runtime"
    "time"

    "m5s/domain"
)

func (as *Service) StartPollTicker() {
    as.logger.Info(
        "Starting poll ticker",
        "pollInterval", as.config.pollInterval,
    )

    ticker := time.NewTicker(as.config.pollInterval)

    for range ticker.C {
        currentValues := getCurrentValues()

        metrics := make([]*domain.Metric, 0, len(currentValues))

        for name, value := range currentValues {
            metrics = append(metrics, domain.NewGauge(name, value))
        }

        metrics = append(metrics, domain.NewCounter("PollCount", 1))

        if err := as.storage.UpdateMetrics(metrics); err != nil {
            as.logger.Error("update agent metrics", "error", err)
        }
    }
}

func getCurrentValues() map[string]float64 {
    var mStat runtime.MemStats

    runtime.ReadMemStats(&mStat)

    return map[string]float64{
        "Alloc":         float64(mStat.Alloc),
        "BuckHashSys":   float64(mStat.BuckHashSys),
        "Frees":         float64(mStat.Frees),
        "GCCPUFraction": mStat.GCCPUFraction,
        "GCSys":         float64(mStat.GCSys),
        "HeapAlloc":     float64(mStat.HeapAlloc),
        "HeapIdle":      float64(mStat.HeapIdle),
        "HeapInuse":     float64(mStat.HeapInuse),
        "HeapObjects":   float64(mStat.HeapObjects),
        "HeapReleased":  float64(mStat.HeapReleased),
        "HeapSys":       float64(mStat.HeapSys),
        "LastGC":        float64(mStat.LastGC),
        "Lookups":       float64(mStat.Lookups),
        "MCacheInuse":   float64(mStat.MCacheInuse),
        "MCacheSys":     float64(mStat.MCacheSys),
        "MSpanInuse":    float64(mStat.MSpanInuse),
        "MSpanSys":      float64(mStat.MSpanSys),
        "Mallocs":       float64(mStat.Mallocs),
        "NextGC":        float64(mStat.NextGC),
        "NumForcedGC":   float64(mStat.NumForcedGC),
        "NumGC":         float64(mStat.NumGC),
        "OtherSys":      float64(mStat.OtherSys),
        "PauseTotalNs":  float64(mStat.PauseTotalNs),
        "StackInuse":    float64(mStat.StackInuse),
        "StackSys":      float64(mStat.StackSys),
        "Sys":           float64(mStat.Sys),
        "TotalAlloc":    float64(mStat.TotalAlloc),
        "RandomValue":   rand.Float64(),
    }
}
