//go:build integration

package integration

import (
	"context"
	"net/http/httptest"
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/server"
)

// setupServer starts a server and log capture for the test. Each test gets its own server and logs (no shared state).
// Call cleanup (e.g. defer cleanup()) when done. Tests can use t.Parallel(); each has its own server.
func setupServer(t testing.TB) (baseURL string, logs *LogCapture, cleanup func()) {
	t.Helper()

	logCapture, log := NewLogCapture()

	gbClient, err := guestbook.NewClient(context.Background())
	if err != nil {
		t.Fatalf("guestbook client: %v (is LocalStack up? run: make local-resources-up)", err)
	}

	handler := server.NewHandler(gbClient, log, "test")
	srv := httptest.NewServer(handler)

	cleanup = func() {
		srv.Close()
		logCapture.Close()
	}
	return srv.URL, logCapture, cleanup
}
