//go:build integration

package integration

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestGuestbookCreate(t *testing.T) {
	t.Parallel()
	baseURL, _, cleanup := setupServer(t)
	defer cleanup()

	var entry guestbook.Entry
	body := handlers.CreateGuestbookRequest{Name: "Integration User", Message: "Hello from integration"}

	client := testclient.NewClient(baseURL)
	err := client.PostJSONExpectCreated(testclient.APIGuestbook, body, &entry)
	require.NoError(t, err)
	require.NotEmpty(t, entry.ID)
	require.Equal(t, body.Name, entry.Name)
	require.Equal(t, body.Message, entry.Message)
	require.False(t, entry.CreatedAt.IsZero())
}
