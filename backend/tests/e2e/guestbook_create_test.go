//go:build e2e

package e2e

import (
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/handlers"
	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestGuestbookCreate(t *testing.T) {
	var entry guestbook.Entry
	body := handlers.CreateGuestbookRequest{Name: "E2E User", Message: "Hello from e2e"}

	client := testclient.NewClient(baseURL)
	err := client.PostJSONExpectCreated(testclient.APIGuestbook, body, &entry)
	require.NoError(t, err)
	require.NotEmpty(t, entry.ID)
	require.Equal(t, body.Name, entry.Name)
	require.Equal(t, body.Message, entry.Message)
	require.False(t, entry.CreatedAt.IsZero())
}
