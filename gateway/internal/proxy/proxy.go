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

	u, err := url.Parse(target)
	if err != nil {
		logrus.Fatalf("invalid upstream: %s", target)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 10 * time.Second,
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
