package handlers

import (
    "context"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"

    "m5s/domain"
    "m5s/internal/server"
    memorystorage "m5s/internal/storage/memory_storage"
)

func TestHandler_GetMetric(t *testing.T) {
    ctx := context.Background()

    serverStorage := memorystorage.NewMemStorage()
    serverService := server.NewServerService(serverStorage)

    _ = serverStorage.UpdateMetrics(
        ctx,
        []*domain.Metric{
            {
                Name:     "counterMetric",
                Type:     domain.MetricTypeCounter,
                IntValue: 111,
            },
            {
                Name:       "gaugeMetric",
                Type:       domain.MetricTypeGauge,
                FloatValue: 0.111,
            },
        },
    )

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Get("/value/{metricType}/{metricName}", handler.GetMetric)

    testServer := httptest.NewServer(r)
    defer testServer.Close()

    testClient := http.Client{
        Timeout: time.Second * 1,
    }

    type want struct {
        code        int
        contentType string
        value       string
    }

    tests := []struct {
        name   string
        method string
        target string
        want   want
    }{
        {
            name:   "negative__wrong_method",
            method: http.MethodPost,
            target: "/value/gauge/gaugeMetric",
            want: want{
                code: http.StatusMethodNotAllowed,
            },
        }, {
            name:   "negative__invalid_metric_type",
            method: http.MethodGet,
            target: "/value/counters/counterMetric",
            want: want{
                code:        http.StatusBadRequest,
                contentType: "text/plain",
            },
        }, {
            name:   "negative__unknown_metric_name",
            method: http.MethodGet,
            target: "/value/gauge/gaugeShvalue",
            want: want{
                code:        http.StatusNotFound,
                contentType: "text/plain",
            },
        }, {
            name:   "positive__200_gauge",
            method: http.MethodGet,
            target: "/value/gauge/gaugeMetric",
            want: want{
                code:        http.StatusOK,
                contentType: "text/plain",
                value:       "0.111",
            },
        }, {
            name:   "positive__200_counter",
            method: http.MethodGet,
            target: "/value/counter/counterMetric",
            want: want{
                code:        http.StatusOK,
                contentType: "text/plain",
                value:       "111",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(tt.method, testServer.URL+tt.target, nil)
            if err != nil {
                t.Fatal(err)
            }

            req.Header.Set("Content-Type", "text/plain")

            res, err := testClient.Do(req)
            if err != nil {
                t.Fatal(err)
            }
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            body, err := io.ReadAll(res.Body)
            if err != nil {
                t.Fatal(err)
            }

            // Check content
            assert.Equal(t, tt.want.value, string(body))

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}