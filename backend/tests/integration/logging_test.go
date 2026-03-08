//go:build integration

package integration

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestServerLogsRequest(t *testing.T) {
	t.Parallel()
	baseURL, logs, cleanup := setupServer(t)
	defer cleanup()

	client := testclient.NewClient(baseURL)
	err := client.GetJSON(testclient.APIHealth, &map[string]any{})
	require.NoError(t, err)

	entries := logs.Entries()
	require.NotEmpty(t, entries, "expected at least one log entry after request")

	var found bool
	for _, e := range entries {
		if e.Message != "request" {
			continue
		}
		method, _ := e.Attrs["method"].(string)
		path, _ := e.Attrs["path"].(string)
		status, _ := e.Attrs["status"].(int64)
		if method == "GET" && path == "/api/health" && status == 200 {
			found = true
			break
		}
	}
	require.True(t, found, "expected a request log entry with method=GET, path=/api/health, status=200; got entries: %+v", entries)
}
