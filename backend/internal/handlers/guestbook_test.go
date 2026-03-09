//go:build test

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpErrors "bitbucket.org/atlassian/unsw-comp3900-app/internal/errors"
	"bitbucket.org/atlassian/unsw-comp3900-app/internal/guestbook"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuestbook_Create_Success(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{}, "Guestbook")
	create, _ := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Post("/", create)

	body := []byte(`{"name":"Alice","message":"Hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var entry guestbook.Entry
	err := json.NewDecoder(rr.Body).Decode(&entry)
	require.NoError(t, err)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, "Alice", entry.Name)
	assert.Equal(t, "Hello", entry.Message)
	assert.False(t, entry.CreatedAt.IsZero())
}

func TestGuestbook_Create_BadJSON(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{}, "Guestbook")
	create, _ := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Post("/", create)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGuestbook_Create_EmptyNameOrMessage(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{}, "Guestbook")
	create, _ := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Post("/", create)

	for _, body := range []string{`{"name":"","message":"x"}`, `{"name":"x","message":""}`} {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code, "body: %s", body)
	}
}

func TestGuestbook_Create_ClientError(t *testing.T) {
	putErr := errors.New("dynamo failed")
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{PutItemErr: putErr}, "Guestbook")
	create, _ := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Post("/", create)

	body := []byte(`{"name":"Alice","message":"Hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGuestbook_List_Success(t *testing.T) {
	createdAt := "2025-01-01T12:00:00Z"
	it := guestbook.Item{ID: "id1", GSIPk: "ENTRY", CreatedAt: createdAt, Name: "Alice", Message: "Hi"}
	av, err := attributevalue.MarshalMap(it)
	require.NoError(t, err)
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{
		QueryOut: &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{av}},
	}, "Guestbook")
	_, list := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/", list)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var entries []guestbook.Entry
	err = json.NewDecoder(rr.Body).Decode(&entries)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "id1", entries[0].ID)
	assert.Equal(t, "Alice", entries[0].Name)
	assert.Equal(t, "Hi", entries[0].Message)
}

func TestGuestbook_List_Empty(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{
		QueryOut: &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}},
	}, "Guestbook")
	_, list := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/", list)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var entries []guestbook.Entry
	err := json.NewDecoder(rr.Body).Decode(&entries)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestGuestbook_List_ClientError(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{QueryErr: errors.New("dynamo failed")}, "Guestbook")
	_, list := Guestbook(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/", list)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGuestbookGet_Success(t *testing.T) {
	createdAt := "2025-01-01T12:00:00Z"
	it := guestbook.Item{ID: "id1", GSIPk: "ENTRY", CreatedAt: createdAt, Name: "Alice", Message: "Hi"}
	av, err := attributevalue.MarshalMap(it)
	require.NoError(t, err)
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{
		GetItemOut: &dynamodb.GetItemOutput{Item: av},
	}, "Guestbook")
	get := GuestbookGet(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/{id}", get)

	req := httptest.NewRequest(http.MethodGet, "/id1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var entry guestbook.Entry
	err = json.NewDecoder(rr.Body).Decode(&entry)
	require.NoError(t, err)
	assert.Equal(t, "id1", entry.ID)
	assert.Equal(t, "Alice", entry.Name)
	assert.Equal(t, "Hi", entry.Message)
}

func TestGuestbookGet_NotFound(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{
		GetItemOut: &dynamodb.GetItemOutput{Item: nil},
	}, "Guestbook")
	get := GuestbookGet(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/{id}", get)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGuestbookGet_ClientError(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{GetItemErr: errors.New("dynamo failed")}, "Guestbook")
	get := GuestbookGet(client)

	r := chi.NewRouter()
	r.Use(httpErrors.HandleHTTPError)
	r.Get("/{id}", get)

	req := httptest.NewRequest(http.MethodGet, "/id1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGuestbookGet_EmptyID(t *testing.T) {
	client := guestbook.NewClientWithDDB(&guestbook.MockDDB{}, "Guestbook")
	get := GuestbookGet(client)

	// Chi's "/{id}" never matches with empty id, so inject route context so the handler sees id=="" and returns 400.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	httpErrors.HandleHTTPError(get).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}