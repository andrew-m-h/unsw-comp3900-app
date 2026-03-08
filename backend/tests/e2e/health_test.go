//go:build e2e

package e2e

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	var healthResponse handlers.HealthResponse

	client := testclient.NewClient(baseURL)
	err := client.GetJSON(testclient.APIHealth, &healthResponse)
	require.NoError(t, err)
	require.Equal(t, handlers.HealthStatusOK, healthResponse.Status)
}
