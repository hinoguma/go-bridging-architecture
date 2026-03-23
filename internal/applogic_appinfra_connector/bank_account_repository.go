package applogic_appinfra_connector

import (
	driven_infra_aws "app/internal/appinfra/aws"
	"app/internal/appinfra/postgres"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/timer"
	"context"
	"database/sql"
)

type bankAccountRepositoryDynamoDBConnector struct {
	dynamo driven_infra_aws.DynamoDBClient
}

func (b bankAccountRepositoryDynamoDBConnector) Create(ctx context.Context, account model.BankAccount, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDBConnector) Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDBConnector) Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func NewBankAccountRepositoryDynamoDBConnector(
	dynamo driven_infra_aws.DynamoDBClient,
) repository.BankAccountRepository {
	return bankAccountRepositoryDynamoDBConnector{
		dynamo: dynamo,
	}
}

type bankAccountRepositoryPostgresConnector struct {
	client postgres.SQLClient
}

func (repo bankAccountRepositoryPostgresConnector) Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
	var row *sql.Row
	var err error

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := bankAccountDTO{}

	row, err = repo.client.GetRowByID(
		ctx,
		postgres.TableBankAccounts,
		postgres.SQLStrID(id.String()),
		sqlItem.Columns(),
		sqlOptions,
	)
	if err != nil {
		return model.BankAccount{}, errors.Lift(err)
	}

	err = sqlItem.SetBySQLRow(row)
	if err != nil {
		return model.BankAccount{}, errors.LiftWithCtx(err, ctx)
	}
	return sqlItem.ToModel(), nil
}

func (repo bankAccountRepositoryPostgresConnector) Create(
	ctx context.Context,
	item model.BankAccount,
	optionaltFuncs ...model.DBOperationOptionalFunc,
) (model.BankAccount, error) {

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := bankAccountDTO{}
	sqlItem.SetByModel(item)

	sqlRow, err := repo.client.CreateRow(ctx, postgres.TableBankAccounts, &sqlItem, sqlOptions)
	if err != nil {
		return item, errors.Lift(err)
	}
	err = sqlItem.SetBySQLRow(sqlRow)
	if err != nil {
		return model.BankAccount{}, errors.LiftWithCtx(err, ctx)
	}
	return sqlItem.ToModel(), nil
}

func (repo bankAccountRepositoryPostgresConnector) Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error {
	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	updateRequests := convertUpdateBankAccountRequestToUpdateRequests(request)
	_, err := repo.client.UpdateRowByStrID(
		ctx,
		postgres.TableBankAccounts,
		postgres.SQLStrID(request.ID.String()),
		updateRequests,
		sqlOptions,
	)
	if err != nil {
		return errors.Lift(err)
	}
	return nil
}

func NewBankAccountRepositoryPostgresConnector(client postgres.SQLClient) repository.BankAccountRepository {
	return bankAccountRepositoryPostgresConnector{
		client: client,
	}
}

type bankAccountDTO struct {
	postgres.BankAccount
}

func (item *bankAccountDTO) SetByModel(model model.BankAccount) {
	item.Id = model.ID.String()
	item.BankUserId = model.BankUserID.String()
	item.AccountType = model.Type.String()
	item.AmountNumber = model.Amount.Amount
	item.AmountCurrency = model.Amount.Currency.String()
	item.LastTransactionRecordId = model.LastTransactionRecordID.String()
	item.LastTransactionTime = model.LastTransactionTime.Unix()
	item.CreatedAt = model.CreatedAt.Unix()
	item.UpdatedAt = model.UpdatedAt.Unix()
}

func (item bankAccountDTO) ToModel() model.BankAccount {
	modelItem := model.BankAccount{
		ID: model.BankAccountID(item.Id),
		HasBankUserID: model.HasBankUserID{
			BankUserID: model.BankUserID(item.BankUserId),
		},
		Type: model.BankAccountType(item.AccountType),
		Amount: model.Money{
			Amount:   item.AmountNumber,
			Currency: model.Currency(item.AmountCurrency),
		},
		LastTransactionRecordID: model.TransactionRecordID(item.LastTransactionRecordId),
		LastTransactionTime:     timer.TimeFromInt64(item.LastTransactionTime),
	}
	modelItem.CreatedAt = timer.TimeFromInt64(item.CreatedAt)
	modelItem.UpdatedAt = timer.TimeFromInt64(item.UpdatedAt)
	return modelItem
}

func convertUpdateBankAccountRequestToUpdateRequests(request model.UpdateBankAccountRequest) postgres.UpdateFieldRequests {
	var updateRequests postgres.UpdateFieldRequests = make([]postgres.UpdateFieldRequest, 0)
	if request.Amount != nil {
		updateRequests.Append(postgres.ColumnAmountNumber, request.Amount.Amount)
		updateRequests.Append(postgres.ColumnAmountCurrency, request.Amount.Currency.String())
	}
	if request.LastTransactionRecordID != nil {
		updateRequests.Append(postgres.ColumnLastTransactionRecordID, request.LastTransactionRecordID.String())
	}
	if request.LastTransactionTime != nil {
		updateRequests.Append(postgres.ColumnLastTransactionTime, request.LastTransactionTime.Unix())
	}
	return updateRequests
}
