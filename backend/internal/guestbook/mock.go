//go:build test

package guestbook

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// MockDDB is a DynamoDB API implementation for tests. Set the Out/Err fields to control behaviour.
// Use from this package or other packages that need a guestbook.Client with a fake DDB.
type MockDDB struct {
	QueryOut   *dynamodb.QueryOutput
	QueryErr   error
	GetItemOut *dynamodb.GetItemOutput
	GetItemErr error
	PutItemErr error
}

// Query returns MockDDB.QueryOut and MockDDB.QueryErr.
func (m *MockDDB) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return m.QueryOut, m.QueryErr
}

// GetItem returns MockDDB.GetItemOut and MockDDB.GetItemErr.
func (m *MockDDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.GetItemOut, m.GetItemErr
}

// PutItem returns an empty PutItemOutput and MockDDB.PutItemErr.
func (m *MockDDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, m.PutItemErr
}
