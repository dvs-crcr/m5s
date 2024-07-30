package agent

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "m5s/domain"
    "m5s/internal/models"
)

func (as *Service) StartReportTicker() {
    var client = &http.Client{
        Timeout: time.Second * 1,
        Transport: &http.Transport{
            DisableKeepAlives: false,
        },
    }

    as.logger.Info(
        "Starting report ticker",
        "reportInterval", as.config.reportInterval,
    )

    ticker := time.NewTicker(as.config.reportInterval)

    for range ticker.C {
        for _, metric := range as.repo.GetMetricsList() {
            if err := makeRequest(
                client,
                fmt.Sprintf("http://%s/update/", as.config.serverAddr),
                metric,
            ); err != nil {
                as.logger.Error("reporter request", "error", err)
            }
        }
    }
}

func makeRequest(
    client *http.Client,
    uri string,
    metric *domain.Metric,
) error {
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

    var buf bytes.Buffer

    zw := gzip.NewWriter(&buf)
    if _, err = zw.Write(bytesMetric); err != nil {
        return err
    }
    zw.Close()

    request, err := http.NewRequest(
        http.MethodPost,
        uri,
        &buf,
    )
    if err != nil {
        return fmt.Errorf("execute http request: %v", err)
    }

    request.Close = true

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Content-Encoding", "gzip")

    response, err := client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    return nil
}