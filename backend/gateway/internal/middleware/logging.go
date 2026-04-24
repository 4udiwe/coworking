package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type recorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *recorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rec := &recorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)
		fields := logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rec.status,
			"size":     rec.size,
			"duration": time.Since(start).String(),
		}

		if rec.status >= 400 {
			logrus.WithFields(fields).Warn("request")

		}
		logrus.WithFields(fields).Info("request")
	})
}
