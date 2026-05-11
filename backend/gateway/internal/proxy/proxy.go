package proxy

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

func New(target string) *httputil.ReverseProxy {
	return newProxy(target, false)
}

func NewForMedia(target string) *httputil.ReverseProxy {
	return newProxy(target, true)
}

func newProxy(target string, overrideCORS bool) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		logrus.Fatalf("invalid upstream: %s", target)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 10 * time.Second,
	}

	if overrideCORS {
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Methods")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Allow-Credentials")
			resp.Header.Del("Access-Control-Expose-Headers")
			resp.Header.Del("Vary")

			origin := resp.Request.Header.Get("Origin")
			if origin != "" {
				resp.Header.Set("Access-Control-Allow-Origin", origin)
			} else {
				resp.Header.Set("Access-Control-Allow-Origin", "*")
			}
			resp.Header.Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			resp.Header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			resp.Header.Set("Vary", "Origin")

			return nil
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logrus.WithError(err).
			WithField("upstream", target).
			Error("proxy error")

		w.WriteHeader(http.StatusBadGateway)
		io.WriteString(w, "bad gateway")
	}

	return proxy
}