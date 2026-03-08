package guestbook

import (
	"context"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ListEntries returns guestbook entries newest first (query GSI by-created, then sort in app for consistent order).
func (c *Client) ListEntries(ctx context.Context) ([]Entry, error) {
	out, err := c.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              &c.table,
		IndexName:              ptr(gsiByNameCreated),
		KeyConditionExpression: ptr("gsiPk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: gsiPkValue},
		},
		ScanIndexForward: ptr(false), // descending by createdAt (DynamoDB); LocalStack may ignore, so we sort below
	})
	if err != nil {
		return nil, err
	}

	var items []Item
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
	// Sort newest first; DynamoDB/LocalStack may not respect ScanIndexForward for GSI.
	sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt.After(entries[j].CreatedAt) })
	return entries, nil
}

func ptr[T any](v T) *T { return &v }
