package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Route maps a hostname to an upstream port.
type Route struct {
	Host    string
	AppPort int
}

// NewMultiHandler returns an HTTP handler that routes requests to multiple
// upstream services based on the Host header.
func NewMultiHandler(routes []Route) http.Handler {
	handlers := make(map[string]http.Handler, len(routes))
	for _, r := range routes {
		handlers[r.Host+".localhost"] = newReverseProxy(r.Host, r.AppPort)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqHost := stripPort(r.Host)
		if h, ok := handlers[reqHost]; ok {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no route for host %q\n", r.Host)
	})
}

// NewHandler returns an HTTP handler that routes requests by Host header.
// Requests with Host matching "<host>.localhost" (with or without port) are
// proxied to http://127.0.0.1:<appPort>. All others get 404.
func NewHandler(host string, appPort int) http.Handler {
	rp := newReverseProxy(host, appPort)
	expectedHost := host + ".localhost"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if MatchHost(r.Host, expectedHost) {
			rp.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no route for host %q\n", r.Host)
	})
}

func newReverseProxy(_ string, appPort int) http.Handler {
	target := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", appPort),
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	defaultDirector := rp.Director
	rp.Director = func(req *http.Request) {
		origHost := req.Host
		defaultDirector(req)
		req.Header.Set("X-Forwarded-Host", stripPort(origHost))
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "upstream unreachable: %v\n", err)
	}

	return rp
}

// MatchHost compares the request host (which may include a port) against
// the expected hostname. The port is stripped before comparison.
func MatchHost(requestHost, expected string) bool {
	return stripPort(requestHost) == expected
}

func stripPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		// No port present
		return hostport
	}
	return host
}
