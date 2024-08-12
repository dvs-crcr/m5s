package handlers

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "m5s/domain"
    "m5s/internal/repository"
    "m5s/internal/server"
)

func TestHandler_GetMetric(t *testing.T) {
    serverRepository := repository.NewInMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

    _ = serverRepository.UpdateMetrics([]*domain.Metric{
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
    })

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
            require.NoError(t, err)

            req.Header.Set("Content-Type", "text/plain")

            res, err := testClient.Do(req)
            require.NoError(t, err)
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            body, err := io.ReadAll(res.Body)
            require.NoError(t, err)

            // Check content
            assert.Equal(t, tt.want.value, string(body))

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}

func TestHandler_GetMetricJSON(t *testing.T) {
    serverRepository := repository.NewInMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

    _ = serverRepository.Update(&domain.Metric{
        Name:       "gaugeMetric",
        Type:       domain.MetricTypeGauge,
        FloatValue: 0.111,
    })

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
            require.NoError(t, err)

            req.Header.Set("Content-Type", tt.contentType)

            res, err := testClient.Do(req)
            require.NoError(t, err)
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            body, err := io.ReadAll(res.Body)
            require.NoError(t, err)

            // Check content
            if tt.want.value != "" {
                assert.JSONEq(t, tt.want.value, string(body))
            }

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}

func TestHandler_GetMetricsList(t *testing.T) {
    serverRepository := repository.NewInMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

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

    _ = serverRepository.UpdateMetrics(mockMetrics)

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
            require.NoError(t, err)

            res, err := testClient.Do(req)
            require.NoError(t, err)
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            body, err := io.ReadAll(res.Body)
            require.NoError(t, err)

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
