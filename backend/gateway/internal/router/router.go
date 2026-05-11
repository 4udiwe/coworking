package router

import (
	"net/http"

	"github.com/4udiwe/coworking/gateway/internal/config"
	"github.com/4udiwe/coworking/gateway/internal/middleware"
	"github.com/4udiwe/coworking/gateway/internal/proxy"
)

// маршруты, которые проксируются напрямую в MinIO
var mediaRoutes = map[string]bool{
	"/media": true,
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func New(cfg *config.Config) http.Handler {
	mainMux := http.NewServeMux()
	mediaMux := http.NewServeMux()

	for _, route := range cfg.Routes {
		if mediaRoutes[route.Path] {
			// Медиа: CORS управляется только через ModifyResponse в прокси,
			// чтобы избежать дублирования заголовков от MinIO и gateway
			p := proxy.NewForMedia(route.Upstream)
			mediaMux.Handle(route.Path+"/", p)
			mediaMux.Handle(route.Path, p)
		} else {
			p := proxy.New(route.Upstream)
			mainMux.Handle(route.Path+"/", p)
			mainMux.Handle(route.Path, p)
		}
	}

	mainMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mainHandler := middleware.RateLimit(cfg.RateLimit.RequestsPerSecond)(
		middleware.RequestID(
			middleware.Logging(mainMux),
		),
	)

	// Корневой mux: медиа идёт без corsHandler (у него свой CORS в ModifyResponse),
	// все остальные маршруты оборачиваются в corsHandler
	root := http.NewServeMux()
	root.Handle("/media/", mediaMux)
	root.Handle("/media", mediaMux)
	root.Handle("/", corsHandler(mainHandler))

	return root
}