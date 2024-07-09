package domain

import (
    "math/rand"
    "runtime"
    "sync"
)

type Statistics struct {
    sync.RWMutex
    mStat         runtime.MemStats
    CurrentValues map[string]float64
    PollCount     int64
}

func NewStatistics() *Statistics {
    return &Statistics{
        CurrentValues: make(map[string]float64),
    }
}

func (s *Statistics) Refresh() {
    s.RLock()
    defer s.RUnlock()

    runtime.ReadMemStats(&s.mStat)

    s.CurrentValues["Alloc"] = float64(s.mStat.Alloc)
    s.CurrentValues["BuckHashSys"] = float64(s.mStat.BuckHashSys)
    s.CurrentValues["Frees"] = float64(s.mStat.Frees)
    s.CurrentValues["GCCPUFraction"] = s.mStat.GCCPUFraction
    s.CurrentValues["GCSys"] = float64(s.mStat.GCSys)
    s.CurrentValues["HeapAlloc"] = float64(s.mStat.HeapAlloc)
    s.CurrentValues["HeapIdle"] = float64(s.mStat.HeapIdle)
    s.CurrentValues["HeapInuse"] = float64(s.mStat.HeapInuse)
    s.CurrentValues["HeapObjects"] = float64(s.mStat.HeapObjects)
    s.CurrentValues["HeapReleased"] = float64(s.mStat.HeapReleased)
    s.CurrentValues["HeapSys"] = float64(s.mStat.HeapSys)
    s.CurrentValues["LastGC"] = float64(s.mStat.LastGC)
    s.CurrentValues["Lookups"] = float64(s.mStat.Lookups)
    s.CurrentValues["MCacheInuse"] = float64(s.mStat.MCacheInuse)
    s.CurrentValues["MCacheSys"] = float64(s.mStat.MCacheSys)
    s.CurrentValues["MSpanInuse"] = float64(s.mStat.MSpanInuse)
    s.CurrentValues["MSpanSys"] = float64(s.mStat.MSpanSys)
    s.CurrentValues["Mallocs"] = float64(s.mStat.Mallocs)
    s.CurrentValues["NextGC"] = float64(s.mStat.NextGC)
    s.CurrentValues["NumForcedGC"] = float64(s.mStat.NumForcedGC)
    s.CurrentValues["NumGC"] = float64(s.mStat.NumGC)
    s.CurrentValues["OtherSys"] = float64(s.mStat.OtherSys)
    s.CurrentValues["PauseTotalNs"] = float64(s.mStat.PauseTotalNs)
    s.CurrentValues["StackInuse"] = float64(s.mStat.StackInuse)
    s.CurrentValues["StackSys"] = float64(s.mStat.StackSys)
    s.CurrentValues["Sys"] = float64(s.mStat.Sys)
    s.CurrentValues["TotalAlloc"] = float64(s.mStat.TotalAlloc)
    s.CurrentValues["TotalAlloc"] = float64(s.mStat.TotalAlloc)
    s.CurrentValues["RandomValue"] = rand.Float64()

    s.PollCount++
}
