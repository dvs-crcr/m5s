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

func TestHandler_GetMetricsList(t *testing.T) {
    ctx := context.Background()

    serverStorage := memorystorage.NewMemStorage()
    serverService := server.NewServerService(serverStorage)

    mockMetrics := []*domain.Metric{
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
    }

    _ = serverStorage.UpdateMetrics(ctx, mockMetrics)

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Get("/", handler.GetMetricsList)

    testServer := httptest.NewServer(r)
    defer testServer.Close()

    testClient := http.Client{
        Timeout: time.Second * 1,
    }

    type want struct {
        code        int
        contentType string
        result      bool
    }

    tests := []struct {
        name   string
        method string
        want   want
    }{
        {
            name:   "negative__wrong_method",
            method: http.MethodPut,
            want: want{
                code: http.StatusMethodNotAllowed,
            },
        }, {
            name:   "positive__200",
            method: http.MethodGet,
            want: want{
                code:        http.StatusOK,
                contentType: "text/html",
                result:      true,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(
                tt.method,
                testServer.URL+"/",
                nil,
            )
            if err != nil {
                t.Fatal(err)
            }

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
            if tt.want.result {
                for _, metric := range mockMetrics {
                    assert.Contains(t, string(body), metric.String())
                }
            }

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}
