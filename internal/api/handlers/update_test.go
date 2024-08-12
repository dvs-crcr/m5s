package handlers

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "m5s/internal/repository"
    "m5s/internal/server"
)

func TestHandler_Update(t *testing.T) {
    serverRepository := repository.NewInMemStorage()
    serverService := server.NewServerService(
        serverRepository,
    )

    handler := NewHandler(
        serverService,
    )

    r := chi.NewRouter()
    r.Route("/update", func(r chi.Router) {
        r.Post("/", handler.UpdateJSON)
        r.Post("/{metricType}/{metricName}/{metricValue}", handler.Update)
    })

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

            req, err := http.NewRequest(tt.method, testServer.URL+tt.target, nil)
            require.NoError(t, err)

            res, err := testClient.Do(req)
            require.NoError(t, err)
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            _, err = io.ReadAll(res.Body)
            require.NoError(t, err)

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}
