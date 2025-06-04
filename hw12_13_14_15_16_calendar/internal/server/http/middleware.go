package internalhttp

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		url := r.URL.String()
		userAgent := r.UserAgent()
		requestID := middleware.GetReqID(r.Context())

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log after the request is done
		fmt.Printf("request to %s completed in %v method: %s userAgent: %s requestID: %s",
			url,
			time.Since(start),
			method,
			userAgent,
			requestID,
		)
	}
	return http.HandlerFunc(fn)
}
