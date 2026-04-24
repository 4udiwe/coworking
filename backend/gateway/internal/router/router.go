package router

import (
	"net/http"

	"github.com/4udiwe/coworking/gateway/internal/config"
	"github.com/4udiwe/coworking/gateway/internal/middleware"
	"github.com/4udiwe/coworking/gateway/internal/proxy"
)

// CORSHandler оборачивает handler и отвечает на OPTIONS
func CORSHandler(h http.Handler, allowedOrigins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// Если preflight (OPTIONS), возвращаем 200 и ничего больше
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// иначе — вызываем основной handler
		h.ServeHTTP(w, r)
	})
}

func New(cfg *config.Config) http.Handler {

	mux := http.NewServeMux()

	for _, route := range cfg.Routes {

		p := proxy.New(route.Upstream)

		mux.Handle(route.Path+"/", p)
		mux.Handle(route.Path, p)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RateLimit(cfg.RateLimit.RequestsPerSecond)(
		middleware.RequestID(
			middleware.Logging(mux),
		),
	)

	return CORSHandler(handler, "*")
}
