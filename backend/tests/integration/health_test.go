//go:build integration

package integration

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	t.Parallel()
	baseURL, _, cleanup := setupServer(t)
	defer cleanup()

	client := testclient.NewClient(baseURL)
	var health handlers.HealthResponse
	err := client.GetJSON(testclient.APIHealth, &health)
	require.NoError(t, err)
	require.NotEmpty(t, health.Version)
}
