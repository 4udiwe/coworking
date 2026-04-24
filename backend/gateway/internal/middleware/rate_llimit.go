package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimit(rps int) func(http.Handler) http.Handler {

	limiter := rate.NewLimiter(rate.Limit(rps), rps)

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("rate limit exceeded"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}