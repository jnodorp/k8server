package k8server

import "net/http"

type livezHandler struct{}

// Livez will always respond with HTTP status code 200 (OK).
func Livez() http.Handler {
	return livezHandler{}
}

// ServeHTTP will always respond with HTTP status code 200 (OK).
func (h livezHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}

type readyzHandler struct {
	ready bool
}

// Readyz will respond with HTTP status code 200 (OK) as long as the provided server is accepting requests. Once the
// provided server is shutdown, HTTP status code 503 (Service Unavailable) will be returned.
func Readyz(srv *http.Server) http.Handler {
	handler := &readyzHandler{
		ready: true,
	}

	srv.RegisterOnShutdown(func() {
		handler.ready = false
	})

	return handler
}

// ServeHTTP will respond with HTTP status code 200 (OK) as long as the provided server is accepting requests. Once the
// provided server is shutdown, HTTP status code 503 (Service Unavailable) will be returned.
func (h *readyzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.ready {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"title":"OK","detail":"The server is running and does accept new connections."}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/problem+json")
		_, _ = w.Write([]byte(`{"status":503,"title":"Service Unavailable","detail":"The server is shutting down and does not accept new connections."}`))
	}
}
