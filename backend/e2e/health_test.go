package e2e

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	var healthResponse handlers.HealthResponse

	client := NewClient(baseURL)
	err := client.GetJSON(APIHealth, &healthResponse)
	require.NoError(t, err)
	require.Equal(t, handlers.HealthStatusOK, healthResponse.Status)
}
