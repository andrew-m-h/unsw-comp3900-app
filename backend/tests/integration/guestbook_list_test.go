//go:build integration

package integration

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestGuestbookList(t *testing.T) {
	t.Parallel()
	baseURL, _, cleanup := setupServer(t)
	defer cleanup()

	var list []guestbook.Entry

	client := testclient.NewClient(baseURL)
	err := client.GetJSON(testclient.APIGuestbook, &list)
	require.NoError(t, err)
	require.NotNil(t, list)
}

func TestGuestbookList_IncludesCreatedEntry(t *testing.T) {
	t.Parallel()
	baseURL, _, cleanup := setupServer(t)
	defer cleanup()

	var created guestbook.Entry
	body := handlers.CreateGuestbookRequest{Name: "List Integration", Message: "Should appear in list"}

	client := testclient.NewClient(baseURL)
	err := client.PostJSONExpectCreated(testclient.APIGuestbook, body, &created)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	var list []guestbook.Entry
	err = client.GetJSON(testclient.APIGuestbook, &list)
	require.NoError(t, err)

	var found bool
	for _, e := range list {
		if e.ID == created.ID {
			found = true
			require.Equal(t, created.Name, e.Name)
			require.Equal(t, created.Message, e.Message)
			break
		}
	}
	require.True(t, found, "created entry %q should be in list (len=%d)", created.ID, len(list))
}
