package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchHost(t *testing.T) {
	tests := []struct {
		requestHost string
		expected    string
		want        bool
	}{
		{"myapp.localhost", "myapp.localhost", true},
		{"myapp.localhost:7890", "myapp.localhost", true},
		{"other.localhost", "myapp.localhost", false},
		{"other.localhost:7890", "myapp.localhost", false},
		{"localhost", "myapp.localhost", false},
		{"localhost:7890", "myapp.localhost", false},
		{"myapp.localhost:3000", "myapp.localhost", true},
		{"", "myapp.localhost", false},
	}

	for _, tt := range tests {
		t.Run(tt.requestHost, func(t *testing.T) {
			if got := MatchHost(tt.requestHost, tt.expected); got != tt.want {
				t.Errorf("MatchHost(%q, %q) = %v, want %v", tt.requestHost, tt.expected, got, tt.want)
			}
		})
	}
}

func TestHandler404(t *testing.T) {
	handler := NewHandler("myapp", 9999)

	req := httptest.NewRequest("GET", "http://wrong.localhost/", nil)
	req.Host = "wrong.localhost"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	body, _ := io.ReadAll(w.Body)
	if len(body) == 0 {
		t.Error("expected non-empty 404 body")
	}
}

func TestHandlerProxies(t *testing.T) {
	// Start a test upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello from upstream"))
	}))
	defer upstream.Close()

	// Extract port from upstream URL
	_, port, _ := splitHostPort(upstream.URL)

	handler := NewHandler("myapp", port)

	req := httptest.NewRequest("GET", "http://myapp.localhost/test", nil)
	req.Host = "myapp.localhost"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != "hello from upstream" {
		t.Errorf("body = %q, want %q", string(body), "hello from upstream")
	}
}

func TestHandlerForwardedHeaders(t *testing.T) {
	var gotHeaders http.Header
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	_, port, _ := splitHostPort(upstream.URL)

	handler := NewHandler("myapp", port)

	req := httptest.NewRequest("GET", "http://myapp.localhost:7890/", nil)
	req.Host = "myapp.localhost:7890"
	req.RemoteAddr = "192.168.1.10:54321"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := gotHeaders.Get("X-Forwarded-For"); got != "192.168.1.10" {
		t.Errorf("X-Forwarded-For = %q, want %q", got, "192.168.1.10")
	}
	if got := gotHeaders.Get("X-Forwarded-Host"); got != "myapp.localhost" {
		t.Errorf("X-Forwarded-Host = %q, want %q", got, "myapp.localhost")
	}
	if got := gotHeaders.Get("X-Forwarded-Proto"); got != "http" {
		t.Errorf("X-Forwarded-Proto = %q, want %q", got, "http")
	}
}

func TestHandler502(t *testing.T) {
	// Use a port that nothing is listening on
	handler := NewHandler("myapp", 1)

	req := httptest.NewRequest("GET", "http://myapp.localhost/", nil)
	req.Host = "myapp.localhost"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
}

// splitHostPort parses an http URL to extract the numeric port.
func splitHostPort(rawURL string) (string, int, error) {
	// rawURL is like "http://127.0.0.1:12345"
	// Strip scheme
	host := rawURL
	if idx := len("http://"); len(rawURL) > idx {
		host = rawURL[idx:]
	}
	var portStr string
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			portStr = host[i+1:]
			host = host[:i]
			break
		}
	}
	port := 0
	for _, c := range portStr {
		port = port*10 + int(c-'0')
	}
	return host, port, nil
}
