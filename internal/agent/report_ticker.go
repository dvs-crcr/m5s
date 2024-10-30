package agent

import (
    "bytes"
    "compress/gzip"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "m5s/domain"
    "m5s/internal/models"
)

func (as *Service) StartReportTicker(ctx context.Context) {
    var client = &http.Client{
        Timeout: time.Second * 1,
        Transport: &http.Transport{
            DisableKeepAlives: false,
        },
    }

    logger.Infow(
        "starting report ticker",
        "reportInterval", as.config.reportInterval,
    )

    ticker := time.NewTicker(as.config.reportInterval)

    for range ticker.C {
        metricsList, err := as.storage.GetMetricsList(ctx)
        if err != nil {
            logger.Errorw(err.Error())
        }

        for _, metric := range metricsList {
            if err := makeRequest(
                client,
                fmt.Sprintf("http://%s/update/", as.config.serverAddr),
                metric,
            ); err != nil {
                logger.Errorw(err.Error())
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

    logger.Debugw(
        "prepare update request",
        "url", request.URL,
        "method", request.Method,
        "ip", request.RemoteAddr,
        "size", request.ContentLength,
        "headers", request.Header,
        "payload", string(bytesMetric),
    )

    response, err := client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    logger.Debugw(
        "update response",
        "status", response.Status,
        "size", response.ContentLength,
        "headers", response.Header,
    )

    return nil
}
