package handlers

import (
    "context"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"

    "m5s/domain"
    "m5s/internal/server"
    memorystorage "m5s/internal/storage/memory_storage"
)

func TestHandler_GetMetricJSON(t *testing.T) {
    ctx := context.Background()

    serverStorage := memorystorage.NewMemStorage()
    serverService := server.NewServerService(serverStorage)

    _ = serverStorage.Update(
        ctx,
        &domain.Metric{
            Name:       "gaugeMetric",
            Type:       domain.MetricTypeGauge,
            FloatValue: 0.111,
        },
    )

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Post("/value", handler.GetMetricJSON)

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
        name        string
        method      string
        contentType string
        payload     string
        want        want
    }{
        {
            name:        "negative__wrong_method",
            method:      http.MethodGet,
            contentType: "application/json",
            payload: `{
	"id": "gaugeMetric",
	"type": "gauge"
}`,
            want: want{
                code: http.StatusMethodNotAllowed,
            },
        }, {
            name:        "negative__wrong_content_type",
            method:      http.MethodPost,
            contentType: "text/plain",
            payload: `{
	"id": "gaugeMetric",
	"type": "gauge"
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "negative__invalid_metric_type",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `{
	"id": "gaugeMetric",
	"type": "gauges"
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "negative__unknown_metric_name",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `{
	"id": "gaugeMetrics",
	"type": "gauge"
}`,
            want: want{
                code: http.StatusNotFound,
            },
        }, {
            name:        "positive__200",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `{
	"id": "gaugeMetric",
	"type": "gauge"
}`,
            want: want{
                code:        http.StatusOK,
                contentType: "application/json",
                value:       "{\"value\":0.111,\"id\":\"gaugeMetric\",\"type\":\"gauge\"}",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(
                tt.method,
                testServer.URL+"/value",
                strings.NewReader(tt.payload),
            )
            if err != nil {
                t.Fatal(err)
            }

            req.Header.Set("Content-Type", tt.contentType)

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
            if tt.want.value != "" {
                assert.JSONEq(t, tt.want.value, string(body))
            }

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}
