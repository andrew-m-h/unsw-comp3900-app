package guestbook

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ErrNotFound is returned when an entry does not exist.
var ErrNotFound = errors.New("guestbook: entry not found")

// GetEntry returns the guestbook entry by id, or ErrNotFound if it does not exist.
func (c *Client) GetEntry(ctx context.Context, id string) (*Entry, error) {
	out, err := c.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &c.table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var it Item
	if err := attributevalue.UnmarshalMap(out.Item, &it); err != nil {
		return nil, err
	}
	t, _ := time.Parse(time.RFC3339, it.CreatedAt)
	return &Entry{
		ID:        it.ID,
		Name:      it.Name,
		Message:   it.Message,
		CreatedAt: t,
	}, nil
}
