package guestbook

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	gsiByNameCreated = "by-created"
	gsiPkValue       = "ENTRY"
)

// DDBAPI is the subset of DynamoDB used by the guestbook client (for testing and injection).
type DDBAPI interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// Client talks to the Guestbook DynamoDB table (env: GUESTBOOK_TABLE_NAME).
type Client struct {
	ddb   DDBAPI
	table string
}

// NewClient builds a guestbook client using default AWS config (env: AWS_ENDPOINT_URL, etc.).
func NewClient(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	table := tableName()
	return &Client{
		ddb:   dynamodb.NewFromConfig(cfg),
		table: table,
	}, nil
}

// NewClientWithDDB builds a client with the given DDB implementation (for tests and other packages).
func NewClientWithDDB(api DDBAPI, table string) *Client {
	return &Client{ddb: api, table: table}
}

func tableName() string {
	if n := os.Getenv("GUESTBOOK_TABLE_NAME"); n != "" {
		return n
	}
	return "Guestbook"
}
