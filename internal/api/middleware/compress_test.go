package middleware

import (
    "compress/flate"
    "compress/gzip"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
)

func TestMiddleware_WithCompression(t *testing.T) {
    r := chi.NewRouter()

    r.Use(WithCompression)

    r.Get("/html", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte("<html><body>content</body></html>"))
    })

    r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"value\":0.111,\"id\":\"gaugeMetric\",\"type\":\"gauge\"}"))
    })

    r.Get("/javascript", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/javascript")
        w.Write([]byte("function f(){evil('alert(\"text\")');}f();"))
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
        name     string
        target   string
        encoding bool
        want     want
    }{
        {
            name:     "positive__with_encoding",
            target:   "/html",
            encoding: true,
            want: want{
                code: http.StatusOK,
            },
        }, {
            name:     "positive__without_encoding",
            target:   "/json",
            encoding: false,
            want: want{
                code: http.StatusOK,
            },
        }, {
            name:     "negative__compress_unsupported_type",
            target:   "/javascript",
            encoding: true,
            want: want{
                code: http.StatusBadRequest,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(
                http.MethodGet,
                testServer.URL+"/",
                nil,
            )
            if err != nil {
                t.Fatal(err)
            }

            if tt.encoding {
                req.Header.Set("Content-Encoding", "gzip")
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

            // Check content

            // Check Content-Type
            assert.Contains(t, res.Header.Get("Content-Type"), tt.want.contentType)
        })
    }
}

func decodeResponseBody(t *testing.T, resp *http.Response) string {
    var reader io.ReadCloser

    switch resp.Header.Get("Content-Encoding") {
    case "gzip":
        var err error

        reader, err = gzip.NewReader(resp.Body)
        if err != nil {
            t.Fatal(err)
        }
    case "deflate":
        reader = flate.NewReader(resp.Body)
    default:
        reader = resp.Body
    }
    defer reader.Close()

    respBody, err := io.ReadAll(reader)
    if err != nil {
        t.Fatal(err)
        return ""
    }

    return string(respBody)
}
