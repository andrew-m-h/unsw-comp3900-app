package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const guestbookTableName = "Guestbook"

// NewGuestbookTable creates the DynamoDB table for the wedding guestbook.
// Partition key: id (string). GSI "by-created": gsiPk (string), createdAt (string) for listing entries by time.
func NewGuestbookTable(scope constructs.Construct, id *string) awsdynamodb.Table {
	table := awsdynamodb.NewTable(scope, id, &awsdynamodb.TableProps{
		TableName:     jsii.String(guestbookTableName),
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("id"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName:    jsii.String("by-created"),
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("gsiPk"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:      &awsdynamodb.Attribute{Name: jsii.String("createdAt"), Type: awsdynamodb.AttributeType_STRING},
	})
	return table
}
