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

func TestGuestbookGet_ExistingEntry_Returns200(t *testing.T) {
	var created guestbook.Entry
	body := handlers.CreateGuestbookRequest{Name: "Get E2E", Message: "Fetch me by id"}

	client := testclient.NewClient(baseURL)
	err := client.PostJSONExpectCreated(testclient.APIGuestbook, body, &created)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	var got guestbook.Entry
	err = client.GetJSON(testclient.APIGuestbook+created.ID, &got)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Name, got.Name)
	require.Equal(t, created.Message, got.Message)
	require.Equal(t, created.CreatedAt.Unix(), got.CreatedAt.Unix())
}

func TestGuestbookGet_NonExistentId_Returns404(t *testing.T) {
	client := testclient.NewClient(baseURL)
	req, err := http.NewRequest(http.MethodGet, client.BaseURL+testclient.APIGuestbook+"nonexistent-id", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
