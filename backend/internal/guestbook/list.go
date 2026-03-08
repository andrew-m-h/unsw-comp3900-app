package guestbook

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ListEntries returns guestbook entries newest first (query GSI by-created).
func (c *Client) ListEntries(ctx context.Context) ([]Entry, error) {
	out, err := c.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              &c.table,
		IndexName:              ptr(gsiByNameCreated),
		KeyConditionExpression: ptr("gsiPk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: gsiPkValue},
		},
		ScanIndexForward: ptr(false), // descending by createdAt
	})
	if err != nil {
		return nil, err
	}

	var items []item
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(items))
	for _, it := range items {
		t, _ := time.Parse(time.RFC3339, it.CreatedAt)
		entries = append(entries, Entry{
			ID:        it.ID,
			Name:      it.Name,
			Message:   it.Message,
			CreatedAt: t,
		})
	}
	return entries, nil
}

func ptr[T any](v T) *T { return &v }
