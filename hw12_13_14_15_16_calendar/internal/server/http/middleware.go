package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewLoggingMiddleware(l Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			method := r.Method
			url := r.URL.String()
			userAgent := r.UserAgent()
			realIP := r.RemoteAddr
			httpVersion := r.Proto
			requestID := middleware.GetReqID(r.Context())

			next.ServeHTTP(ww, r)

			statusCode := ww.Status()

			logMessage := fmt.Sprintf("request to %s completed in %v method: %s userAgent: %s "+
				"realIP: %s httpVersion: %s status: %v requestID: %s",
				url,
				time.Since(start),
				method,
				userAgent,
				realIP,
				httpVersion,
				statusCode,
				requestID,
			)

			l.Info(logMessage)
		})
	}
}
