package middleware

import (
    "net/http"
    "strings"

    "m5s/internal/api"
    "m5s/pkg/compressor"
)

func isClientSupportCompression(
    accept string,
    contentType string,
    acceptEncoding string,
) bool {
    if !api.CheckContentType(
        contentType, "application/json", "text/html",
    ) && accept != "html/text" {
        return false
    }

    if !strings.Contains(acceptEncoding, "gzip") {
        return false
    }

    return true
}

func WithCompression(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ow := w

        if isClientSupportCompression(
            r.Header.Get("Accept"),
            r.Header.Get("Content-Type"),
            r.Header.Get("Accept-Encoding"),
        ) {
            cw := compressor.NewCompressWriter(w)
            ow = cw

            defer cw.Close()
        }

        contentEncoding := r.Header.Get("Content-Encoding")
        if strings.Contains(contentEncoding, "gzip") {
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
