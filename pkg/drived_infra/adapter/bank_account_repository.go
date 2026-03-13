package driven_infra_adapter

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/drived_infra_port"
	"app/pkg/drived_infra/aws"
	"context"
)

type bankAccountRepositoryDynamoDB struct {
	dynamo drived_infra_aws.DynamoDBClient
}

func (b bankAccountRepositoryDynamoDB) Get(ctx context.Context, id model.BankAccountID) (model.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDB) Create(ctx context.Context, account model.BankAccount) error {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDB) Put(ctx context.Context, account model.BankAccount) error {
	//TODO implement me
	panic("implement me")
}

func NewBankAccountRepositoryDynamoDB(
	dynamo drived_infra_aws.DynamoDBClient,
) drived_infra_port.BankAccountRepository {
	return bankAccountRepositoryDynamoDB{
		dynamo: dynamo,
	}
}
