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

    "m5s/internal/server"
    "m5s/internal/storage"
)

func TestHandler_Update(t *testing.T) {
    serverRepository := storage.NewMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.Update)

    testServer := httptest.NewServer(r)
    defer testServer.Close()

    testClient := http.Client{
        Timeout: time.Second * 1,
    }

    type want struct {
        code        int
        contentType string
    }

    tests := []struct {
        name   string
        method string
        target string
        want   want
    }{
        {
            name:   "negative__wrong_method",
            method: http.MethodGet,
            target: "/update/gauge/heapSize/0.1111",
            want: want{
                code:        http.StatusMethodNotAllowed,
                contentType: "",
            },
        }, {
            name:   "negative__wrong_url",
            method: http.MethodPost,
            target: "/update/gauge/0.1111",
            want: want{
                code:        http.StatusNotFound,
                contentType: "text/plain",
            },
        }, {
            name:   "negative__invalid_metric_type",
            method: http.MethodPost,
            target: "/update/counter/someMetric/0.1111",
            want: want{
                code:        http.StatusBadRequest,
                contentType: "text/plain",
            },
        }, {
            name:   "negative__invalid_metric_name",
            method: http.MethodPost,
            target: "/update/gauge//0.1111",
            want: want{
                code:        http.StatusNotFound,
                contentType: "text/plain",
            },
        }, {
            name:   "negative__invalid_metric_value",
            method: http.MethodPost,
            target: "/update/gauge/someMetric/qwerty",
            want: want{
                code:        http.StatusBadRequest,
                contentType: "text/plain",
            },
        }, {
            name:   "positive__200",
            method: http.MethodPost,
            target: "/update/gauge/someMetric/0.000001",
            want: want{
                code:        http.StatusOK,
                contentType: "text/plain",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var err error

            req, err := http.NewRequest(
                tt.method,
                testServer.URL+tt.target,
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

            _, err = io.ReadAll(res.Body)
            if err != nil {
                t.Fatal(err)
            }

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}

func TestHandler_UpdateJSON(t *testing.T) {
    serverRepository := storage.NewMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Post("/update/", handler.UpdateJSON)

    testServer := httptest.NewServer(r)
    defer testServer.Close()

    testClient := http.Client{
        Timeout: time.Second * 1,
    }

    type want struct {
        code        int
        contentType string
    }

    tests := []struct {
        name        string
        method      string
        payload     string
        contentType string
        want        want
    }{
        {
            name:        "negative__wrong_method",
            method:      http.MethodGet,
            contentType: "application/json",
            payload: `
{
	"id": "someMetric",
	"type": "counter",
	"delta": 1
}`,
            want: want{
                code: http.StatusMethodNotAllowed,
            },
        },
        {
            name:        "negative__wrong_json_struct",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	"someData": 123
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        },
        {
            name:        "negative__non_json_payload",
            method:      http.MethodPost,
            contentType: "text/plain",
            payload:     `some non json data`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "negative__invalid_json_payload",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	someData: 123,
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "negative__invalid_metric_type",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	"id": "someMetric",
	"type": "conquer",
	"delta": 1
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "negative__invalid_metric_name",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	"id": "",
	"type": "counter",
	"delta": 1
}`,
            want: want{
                code: http.StatusNotFound,
            },
        }, {
            name:        "negative__invalid_metric_value",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	"id": "someMetric",
	"type": "gauge",
	"value": "someValue"
}`,
            want: want{
                code: http.StatusBadRequest,
            },
        }, {
            name:        "positive__200",
            method:      http.MethodPost,
            contentType: "application/json",
            payload: `
{
	"id": "someGaugeMetric",
	"type": "gauge",
	"value": 0.0001
}`,
            want: want{
                code: http.StatusOK,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var err error

            req, err := http.NewRequest(
                tt.method,
                testServer.URL+"/update/",
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

            _, err = io.ReadAll(res.Body)
            if err != nil {
                t.Fatal(err)
            }

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}
