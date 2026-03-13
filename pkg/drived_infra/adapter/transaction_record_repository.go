package driven_infra_adapter

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/drived_infra_port"
	"app/pkg/drived_infra/aws"
	"context"
)

type transactionRecordRepositoryDynamoDB struct {
	dynamo drived_infra_aws.DynamoDBClient
}

func (b transactionRecordRepositoryDynamoDB) Get(ctx context.Context, id model.TransactionRecordID) (model.TransactionRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (b transactionRecordRepositoryDynamoDB) Create(ctx context.Context, account model.TransactionRecord) error {
	//TODO implement me
	panic("implement me")
}

func (b transactionRecordRepositoryDynamoDB) Put(ctx context.Context, account model.TransactionRecord) error {
	//TODO implement me
	panic("implement me")
}

func NewTransactionRecordRepositoryDynamoDB(
	dynamo drived_infra_aws.DynamoDBClient,
) drived_infra_port.TransactionRecordRepository {
	return transactionRecordRepositoryDynamoDB{
		dynamo: dynamo,
	}
}
