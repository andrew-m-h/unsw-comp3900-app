package guestbook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// CreateEntry stores a new guestbook entry and returns it.
func (c *Client) CreateEntry(ctx context.Context, name, message string) (*Entry, error) {
	id, err := newID()
	if err != nil {
		return nil, fmt.Errorf("guestbook: generate id: %w", err)
	}
	now := time.Now().UTC()
	createdAtStr := now.Format(time.RFC3339)

	it := item{
		ID:        id,
		GSIPk:     gsiPkValue,
		CreatedAt: createdAtStr,
		Name:      name,
		Message:   message,
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return nil, fmt.Errorf("guestbook: marshal item: %w", err)
	}

	_, err = c.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &c.table,
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("guestbook: put item: %w", err)
	}

	return &Entry{
		ID:        id,
		Name:      name,
		Message:   message,
		CreatedAt: now,
	}, nil
}

func newID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
