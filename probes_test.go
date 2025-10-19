package k8server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jnodorp/k8server"

	"github.com/stretchr/testify/assert"
)

func TestLivez(t *testing.T) {
	livez := k8server.Livez()

	resp := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/livez", nil)

	livez.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}

func TestReadyz(t *testing.T) {
	srv := &http.Server{}

	readyz := k8server.Readyz(srv)

	resp := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/readyz", nil)

	readyz.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	body, err := io.ReadAll(resp.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"status":200,"title":"OK","detail":"The server is running and does accept new connections."}`, string(body))

	err = srv.Shutdown(t.Context())
	assert.NoError(t, err)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(t.Context(), "GET", "/readyz", nil)

		readyz.ServeHTTP(resp, req)

		assert.Equal(c, http.StatusServiceUnavailable, resp.Result().StatusCode)

		body, err := io.ReadAll(resp.Result().Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"status":503,"title":"Service Unavailable","detail":"The server is shutting down and does not accept new connections."}`, string(body))
	}, 1*time.Second, 100*time.Millisecond)
}
