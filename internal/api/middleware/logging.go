package middleware

import (
    "net/http"
    "time"
)

type responseData struct {
    status int
    size   int
    header http.Header
}

type loggingResponseWriter struct {
    http.ResponseWriter
    responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
    size, err := r.ResponseWriter.Write(b)
    r.responseData.size += size

    return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
    r.ResponseWriter.WriteHeader(statusCode)
    r.responseData.status = statusCode
    r.responseData.header = r.ResponseWriter.Header()
}

func (m *Middleware) WithRequestLogger(next http.Handler) http.Handler {
    logFn := func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rd := &responseData{
            status: 0,
            size:   0,
        }

        lw := loggingResponseWriter{
            ResponseWriter: w,
            responseData:   rd,
        }

        m.logger.Info(
            "request",
            "method", r.Method,
            "uri", r.RequestURI,
            "ip", r.RemoteAddr,
            "size", r.ContentLength,
            "headers", r.Header,
        )

        next.ServeHTTP(&lw, r)

        m.logger.Info(
            "response",
            "status", rd.status,
            "duration", time.Since(start),
            "size", rd.size,
            "headers", rd.header,
        )
    }

    return http.HandlerFunc(logFn)
}
