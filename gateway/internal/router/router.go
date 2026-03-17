package router

import (
	"net/http"

	"github.com/4udiwe/coworking/gateway/internal/config"
	"github.com/4udiwe/coworking/gateway/internal/middleware"
	"github.com/4udiwe/coworking/gateway/internal/proxy"
)

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

	return handler
}
