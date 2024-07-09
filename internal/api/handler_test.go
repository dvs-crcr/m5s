package api

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "m5s/internal/repository"
    "m5s/server"
)

//nolint:funlen
func TestHandler_Update(t *testing.T) {
    handler := &Handler{
        serverService: server.NewServerService(
            repository.NewInMemStorage(),
        ),
        Mux: http.NewServeMux(),
    }

    handler.Mux.HandleFunc("/update/", handler.Update)

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
            name:   "negative - wrong method",
            method: http.MethodGet,
            target: "/update/gauge/heapSize/0.1111",
            want: want{
                code:        http.StatusMethodNotAllowed,
                contentType: "text/plain",
            },
        }, {
            name:   "negative - wrong url",
            method: http.MethodPost,
            target: "/update/gauge/0.1111",
            want: want{
                code:        http.StatusNotFound,
                contentType: "text/plain",
            },
        }, {
            name:   "negative - invalid metric type",
            method: http.MethodPost,
            target: "/update/counter/someMetric/0.1111",
            want: want{
                code:        http.StatusBadRequest,
                contentType: "text/plain",
            },
        }, {
            name:   "negative - invalid metric name",
            method: http.MethodPost,
            target: "/update/gauge//0.1111",
            want: want{
                code:        http.StatusNotFound,
                contentType: "text/plain",
            },
        }, {
            name:   "negative - invalid metric value",
            method: http.MethodPost,
            target: "/update/gauge/someMetric/qwerty",
            want: want{
                code:        http.StatusBadRequest,
                contentType: "text/plain",
            },
        }, {
            name:   "positive - 200",
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
            r := httptest.NewRequest(tt.method, tt.target, nil)
            w := httptest.NewRecorder()

            handler.Update(w, r)

            res := w.Result()
            defer res.Body.Close()

            // Check StatusCode
            assert.Equal(t, tt.want.code, res.StatusCode)

            _, err := io.ReadAll(res.Body)
            require.NoError(t, err)

            // Check Content-Type
            assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
        })
    }
}
