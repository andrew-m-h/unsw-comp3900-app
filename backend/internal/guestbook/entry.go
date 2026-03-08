package guestbook

import "time"

// Entry is a single guestbook entry.
type Entry struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Name      string    `json:"name" dynamodbav:"name"`
	Message   string    `json:"message" dynamodbav:"message"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

// item is the raw DynamoDB item (includes GSI keys for by-created index).
type item struct {
	ID        string `dynamodbav:"id"`
	GSIPk     string `dynamodbav:"gsiPk"`
	CreatedAt string `dynamodbav:"createdAt"`
	Name      string `dynamodbav:"name"`
	Message   string `dynamodbav:"message"`
}
