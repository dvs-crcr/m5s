package middleware

import (
    "net/http"
    "strings"

    "m5s/internal/api"
    "m5s/pkg/compressor"
)

func isClientSupportCompression(contentType string, acceptEncoding string) bool {
    if !api.CheckContentType(
        contentType, "application/json", "text/html",
    ) {
        return false
    }

    if isGzip := strings.Contains(acceptEncoding, "gzip"); !isGzip {
        return false
    }

    return true
}

func (m *Middleware) WithCompression(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ow := w

        if isClientSupportCompression(
            r.Header.Get("Content-Type"),
            r.Header.Get("Accept-Encoding"),
        ) {
            cw := compressor.NewCompressWriter(w)
            ow = cw

            defer cw.Close()
        }

        contentEncoding := r.Header.Get("Content-Encoding")
        sendsGzip := strings.Contains(contentEncoding, "gzip")
        if sendsGzip {
            cr, err := compressor.NewCompressReader(r.Body)
            if err != nil {
                api.HandleErrors(api.ErrInternal, w)

                return
            }

            r.Body = cr
            defer cr.Close()
        }

        next.ServeHTTP(ow, r)
    })
}
