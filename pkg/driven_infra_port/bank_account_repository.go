package driven_infra_port

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	driven_infra_aws "app/pkg/driven_infra/aws"
	"context"
)

type bankAccountRepositoryDynamoDB struct {
	dynamo driven_infra_aws.DynamoDBClient
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
	dynamo driven_infra_aws.DynamoDBClient,
) repository.BankAccountRepository {
	return bankAccountRepositoryDynamoDB{
		dynamo: dynamo,
	}
}
