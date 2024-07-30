package agent

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "m5s/domain"
    "m5s/internal/models"
)

func (as *Service) StartReportTicker() {
    var client = &http.Client{
        Timeout:   time.Second * 1,
        Transport: &http.Transport{},
    }

    as.logger.Info("Starting report ticker", "duration", as.reportInterval)

    ticker := time.NewTicker(as.reportInterval)

    for range ticker.C {
        for _, metric := range as.repo.GetMetricsList() {
            if err := as.makeRequest(client, metric); err != nil {
                as.logger.Error("make reporter request", "error", err)
            }
        }
    }
}

func (as *Service) makeRequest(client *http.Client, metric *domain.Metric) error {
    modelMetric := &models.Metrics{
        ID:    metric.Name,
        MType: metric.Type.String(),
        Delta: &metric.IntValue,
        Value: &metric.FloatValue,
    }

    bytesMetric, err := json.Marshal(modelMetric)
    if err != nil {
        return err
    }

    //var buf bytes.Buffer
    //
    //zw := gzip.NewWriter(&buf)
    //if _, err = zw.Write(bytesMetric); err != nil {
    //    return err
    //}
    //zw.Close()

    request, err := http.NewRequest(
        http.MethodPost,
        fmt.Sprintf("http://%s/update/", as.serverAddr),
        bytes.NewBuffer(bytesMetric),
    )
    if err != nil {
        return fmt.Errorf("execute http request: %v", err)
    }

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Content-Encoding", "gzip")

    response, err := client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    //if _, err := io.Copy(io.Discard, response.Body); err != nil {
    //    return err
    //}

    return nil
}
