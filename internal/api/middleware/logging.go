package middleware

import (
    "net/http"
    "time"
)

type responseData struct {
    status int
    size   int
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
}

func (m *Middleware) WithLogger(next http.Handler) http.Handler {
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

        next.ServeHTTP(&lw, r)

        duration := time.Since(start)

        m.logger.Info(
            "Request",
            "uri", r.RequestURI,
            "method", r.Method,
            "status", rd.status,
            "duration", duration,
            "size", rd.size,
        )
    }

    return http.HandlerFunc(logFn)
}
