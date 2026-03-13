package drived_infra_aws

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TableSequentialID string = "sequential_ids"

type sequentialIDItem struct {
	ID        string `dynamodbav:"id"`
	CurrentID int64  `dynamodbav:"current_id"`
	CreatedAt int64  `dynamodbav:"created_at,unixtime"`
	UpdatedAt int64  `dynamodbav:"updated_at,unixtime"`
}

func (item *sequentialIDItem) SetByDynamoDBAttrs(attrs map[string]types.AttributeValue) error {
	if attrs == nil {
		return nil
	}
	return attributevalue.UnmarshalMap(attrs, item)
}
