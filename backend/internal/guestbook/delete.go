package guestbook

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DeleteEntry deletes a guestbook entry by id. Returns ErrNotFound if entry does not exist.
func (c *Client) DeleteEntry(ctx context.Context, id string) error {
	// Use DynamoDB's ConditionExpression to only delete if the item exists.
	_, err := c.ddb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &c.table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: ptr("attribute_exists(id)"),
	})
	var conderr *types.ConditionalCheckFailedException
	if err != nil {
		if errors.As(err, &conderr) {
			// Entry not found.
			return ErrNotFound
		}
		return err
	}
	return nil
}
