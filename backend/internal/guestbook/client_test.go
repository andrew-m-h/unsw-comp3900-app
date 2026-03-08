package guestbook

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEntries_Success(t *testing.T) {
	ctx := context.Background()
	created1 := "2025-01-01T12:00:00Z"
	created2 := "2025-01-02T12:00:00Z"
	items := []Item{
		{ID: "id1", GSIPk: gsiPkValue, CreatedAt: created1, Name: "Alice", Message: "Hello"},
		{ID: "id2", GSIPk: gsiPkValue, CreatedAt: created2, Name: "Bob", Message: "Hi"},
	}
	var itemMaps []map[string]types.AttributeValue
	for _, it := range items {
		m, err := attributevalue.MarshalMap(it)
		require.NoError(t, err)
		itemMaps = append(itemMaps, m)
	}
	mock := &MockDDB{
		QueryOut: &dynamodb.QueryOutput{Items: itemMaps},
	}
	client := NewClientWithDDB(mock, "Guestbook")

	entries, err := client.ListEntries(ctx)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	// Sorted newest first
	assert.Equal(t, "id2", entries[0].ID)
	assert.Equal(t, "Bob", entries[0].Name)
	assert.Equal(t, "Hi", entries[0].Message)
	assert.Equal(t, "id1", entries[1].ID)
	assert.Equal(t, "Alice", entries[1].Name)
	assert.Equal(t, "Hello", entries[1].Message)
}

func TestListEntries_QueryError(t *testing.T) {
	ctx := context.Background()
	queryErr := errors.New("dynamodb query failed")
	mock := &MockDDB{QueryErr: queryErr}
	client := NewClientWithDDB(mock, "Guestbook")

	entries, err := client.ListEntries(ctx)
	assert.Error(t, err)
	assert.Nil(t, entries)
	assert.Equal(t, queryErr, err)
}

func TestListEntries_Empty(t *testing.T) {
	ctx := context.Background()
	mock := &MockDDB{QueryOut: &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}}
	client := NewClientWithDDB(mock, "Guestbook")

	entries, err := client.ListEntries(ctx)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestGetEntry_Success(t *testing.T) {
	ctx := context.Background()
	createdAt := "2025-01-01T12:00:00Z"
	it := Item{ID: "id1", GSIPk: gsiPkValue, CreatedAt: createdAt, Name: "Alice", Message: "Hello"}
	av, err := attributevalue.MarshalMap(it)
	require.NoError(t, err)
	mock := &MockDDB{GetItemOut: &dynamodb.GetItemOutput{Item: av}}
	client := NewClientWithDDB(mock, "Guestbook")

	entry, err := client.GetEntry(ctx, "id1")
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, "id1", entry.ID)
	assert.Equal(t, "Alice", entry.Name)
	assert.Equal(t, "Hello", entry.Message)
	tm, _ := time.Parse(time.RFC3339, createdAt)
	assert.True(t, entry.CreatedAt.Equal(tm))
}

func TestGetEntry_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &MockDDB{GetItemOut: &dynamodb.GetItemOutput{Item: nil}}
	client := NewClientWithDDB(mock, "Guestbook")

	entry, err := client.GetEntry(ctx, "missing")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, entry)
}

func TestGetEntry_GetItemError(t *testing.T) {
	ctx := context.Background()
	getErr := errors.New("dynamodb get failed")
	mock := &MockDDB{GetItemErr: getErr}
	client := NewClientWithDDB(mock, "Guestbook")

	entry, err := client.GetEntry(ctx, "id1")
	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Equal(t, getErr, err)
}

func TestCreateEntry_Success(t *testing.T) {
	ctx := context.Background()
	mock := &MockDDB{}
	client := NewClientWithDDB(mock, "Guestbook")

	entry, err := client.CreateEntry(ctx, "Alice", "Hello world")
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, "Alice", entry.Name)
	assert.Equal(t, "Hello world", entry.Message)
	assert.False(t, entry.CreatedAt.IsZero())
}

func TestCreateEntry_PutItemError(t *testing.T) {
	ctx := context.Background()
	putErr := errors.New("dynamodb put failed")
	mock := &MockDDB{PutItemErr: putErr}
	client := NewClientWithDDB(mock, "Guestbook")

	entry, err := client.CreateEntry(ctx, "Alice", "Hello")
	assert.Error(t, err)
	assert.Nil(t, entry)
	assert.Contains(t, err.Error(), "put item")
	assert.ErrorIs(t, err, putErr)
}
