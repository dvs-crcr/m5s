package models

import (
    "math/rand"
    "runtime"
    "sync"
)

type Statistics struct {
    sync.RWMutex
    mStat         runtime.MemStats
    CurrentValues map[string]float64
}

func NewStatistics() *Statistics {
    return &Statistics{
        CurrentValues: make(map[string]float64),
    }
}

func (s *Statistics) Refresh() {
    s.Lock()
    defer s.Unlock()

    runtime.ReadMemStats(&s.mStat)

    s.CurrentValues = map[string]float64{
        "Alloc":         float64(s.mStat.Alloc),
        "BuckHashSys":   float64(s.mStat.BuckHashSys),
        "Frees":         float64(s.mStat.Frees),
        "GCCPUFraction": s.mStat.GCCPUFraction,
        "GCSys":         float64(s.mStat.GCSys),
        "HeapAlloc":     float64(s.mStat.HeapAlloc),
        "HeapIdle":      float64(s.mStat.HeapIdle),
        "HeapInuse":     float64(s.mStat.HeapInuse),
        "HeapObjects":   float64(s.mStat.HeapObjects),
        "HeapReleased":  float64(s.mStat.HeapReleased),
        "HeapSys":       float64(s.mStat.HeapSys),
        "LastGC":        float64(s.mStat.LastGC),
        "Lookups":       float64(s.mStat.Lookups),
        "MCacheInuse":   float64(s.mStat.MCacheInuse),
        "MCacheSys":     float64(s.mStat.MCacheSys),
        "MSpanInuse":    float64(s.mStat.MSpanInuse),
        "MSpanSys":      float64(s.mStat.MSpanSys),
        "Mallocs":       float64(s.mStat.Mallocs),
        "NextGC":        float64(s.mStat.NextGC),
        "NumForcedGC":   float64(s.mStat.NumForcedGC),
        "NumGC":         float64(s.mStat.NumGC),
        "OtherSys":      float64(s.mStat.OtherSys),
        "PauseTotalNs":  float64(s.mStat.PauseTotalNs),
        "StackInuse":    float64(s.mStat.StackInuse),
        "StackSys":      float64(s.mStat.StackSys),
        "Sys":           float64(s.mStat.Sys),
        "TotalAlloc":    float64(s.mStat.TotalAlloc),
        "RandomValue":   rand.Float64(),
    }
}
