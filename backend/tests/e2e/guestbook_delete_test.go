//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestGuestbookDelete_ExistingEntry_Returns204(t *testing.T) {
	var created guestbook.Entry
	body := handlers.CreateGuestbookRequest{Name: "Delete E2E", Message: "Will be deleted"}

	client := testclient.NewClient(baseURL)
	err := client.PostJSONExpectCreated(testclient.APIGuestbook, body, &created)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	req, err := http.NewRequest(http.MethodDelete, client.BaseURL+testclient.APIGuestbook+created.ID, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify entry is gone: GET should return 404
	getReq, err := http.NewRequest(http.MethodGet, client.BaseURL+testclient.APIGuestbook+created.ID, nil)
	require.NoError(t, err)
	getResp, err := client.Do(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestGuestbookDelete_NonExistentId_Returns404(t *testing.T) {
	client := testclient.NewClient(baseURL)
	req, err := http.NewRequest(http.MethodDelete, client.BaseURL+testclient.APIGuestbook+"nonexistent-id", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
