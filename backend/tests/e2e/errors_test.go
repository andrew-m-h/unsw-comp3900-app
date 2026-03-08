//go:build e2e

package e2e

import (
	"bytes"
	"net/http"
	"testing"

	"bitbucket.org/atlassian/unsw-comp3900-app/tests/testclient"
	"github.com/stretchr/testify/require"
)

func TestErrorCodes_MalformedJSON_Returns400(t *testing.T) {
	client := testclient.NewClient(baseURL)
	req, err := http.NewRequest(http.MethodPost, client.BaseURL+testclient.APIGuestbook, bytes.NewReader([]byte("{")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestErrorCodes_NonExistentPath_Returns404(t *testing.T) {
	client := testclient.NewClient(baseURL)
	req, err := http.NewRequest(http.MethodGet, client.BaseURL+"/api/non-existant", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestErrorCodes_GuestbookPut_Returns405(t *testing.T) {
	client := testclient.NewClient(baseURL)
	req, err := http.NewRequest(http.MethodPut, client.BaseURL+testclient.APIGuestbook, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json") // pass RequireContentTypeJSONForBody so Chi returns 405

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestErrorCodes_InvalidAccept_Returns406(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, baseURL+testclient.APIHealth, nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "text/html")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
}

// TestErrorCodes_InvalidContentType_PostGuestbook_Returns415 requires middleware.RequireContentTypeJSONForBody.
// If you get 201 instead of 415, rebuild the backend: docker compose up -d --build backend
func TestErrorCodes_InvalidContentType_PostGuestbook_Returns415(t *testing.T) {
	body := []byte(`{"name":"E2E","message":"test"}`)
	req, err := http.NewRequest(http.MethodPost, baseURL+testclient.APIGuestbook, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}
