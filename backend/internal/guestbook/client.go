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

// Client talks to the Guestbook DynamoDB table (env: GUESTBOOK_TABLE_NAME).
type Client struct {
	ddb    *dynamodb.Client
	table  string
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

func tableName() string {
	if n := os.Getenv("GUESTBOOK_TABLE_NAME"); n != "" {
		return n
	}
	return "Guestbook"
}
