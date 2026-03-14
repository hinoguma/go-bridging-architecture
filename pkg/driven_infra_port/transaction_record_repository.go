package driven_infra_port

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	driven_infra_aws "app/pkg/driven_infra/aws"
	"context"
)

type transactionRecordRepositoryDynamoDB struct {
	dynamo driven_infra_aws.DynamoDBClient
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
	dynamo driven_infra_aws.DynamoDBClient,
) repository.TransactionRecordRepository {
	return transactionRecordRepositoryDynamoDB{
		dynamo: dynamo,
	}
}
