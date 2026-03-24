package applogic_appinfra_connector

import (
	"app/internal/appinfra/aws"
	"app/internal/appinfra/postgres"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/timer"
	"context"
	"database/sql"
)

type transactionRecordRepositoryDynamoDBConnector struct {
	dynamo driven_infra_aws.DynamoDBClient
}

func (repo transactionRecordRepositoryDynamoDBConnector) Get(ctx context.Context, id model.TransactionRecordID, optionalFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (repo transactionRecordRepositoryDynamoDBConnector) Create(ctx context.Context, account model.TransactionRecord, optionalFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	//TODO implement me
	panic("implement me")
}

func NewTransactionRecordRepositoryDynamoDBConnector(
	dynamo driven_infra_aws.DynamoDBClient,
) repository.TransactionRecordRepository {
	return transactionRecordRepositoryDynamoDBConnector{
		dynamo: dynamo,
	}
}

type transactionRecordRepositoryPostgresConnector struct {
	client postgres.SQLClient
}

func NewTransactionRecordRepositoryPostgresConnector(client postgres.SQLClient) repository.TransactionRecordRepository {
	return transactionRecordRepositoryPostgresConnector{
		client: client,
	}
}

func (repo transactionRecordRepositoryPostgresConnector) Get(ctx context.Context, id model.TransactionRecordID, optionalFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	var row *sql.Row
	var err error

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionalFuncs)

	sqlItem := transactionRecordDTO{}

	row, err = repo.client.GetRowByID(
		ctx,
		postgres.TableTransactionRecords,
		postgres.SQLStrID(id.String()),
		sqlItem.Columns(),
		sqlOptions,
	)
	if err != nil {
		return model.TransactionRecord{}, errors.Lift(err)
	}

	err = sqlItem.SetBySQLRow(row)
	if err != nil {
		return model.TransactionRecord{}, errors.LiftWithCtx(err, ctx)
	}
	return sqlItem.ToModel(), nil
}

func (repo transactionRecordRepositoryPostgresConnector) Create(
	ctx context.Context,
	item model.TransactionRecord,
	optionalFuncs ...model.DBOperationOptionalFunc,
) (model.TransactionRecord, error) {

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionalFuncs)

	sqlItem := transactionRecordDTO{}
	sqlItem.SetByModel(item)

	sqlRow, err := repo.client.CreateRow(ctx, postgres.TableTransactionRecords, &sqlItem, sqlOptions)
	if err != nil {
		return item, errors.Lift(err)
	}
	err = sqlItem.SetBySQLRow(sqlRow)
	if err != nil {
		return model.TransactionRecord{}, errors.LiftWithCtx(err, ctx)
	}
	return sqlItem.ToModel(), nil
}

type transactionRecordDTO struct {
	postgres.TransactionRecord
}

func (item *transactionRecordDTO) SetByModel(model model.TransactionRecord) {
	item.Id = model.ID.String()
	item.BankAccountId = model.BankAccountID.String()
	item.BankUserId = model.BankUserID.String()
	item.RecordType = model.Type.String()
	item.DepositAmountNumber = model.DepositAmount.Amount
	item.DepositAmountCurrency = model.DepositAmount.Currency.String()
	item.WithdrawAmountNumber = model.WithdrawAmount.Amount
	item.WithdrawAmountCurrency = model.WithdrawAmount.Currency.String()
	item.AfterAmountNumber = model.AfterAmount.Amount
	item.AfterAmountCurrency = model.AfterAmount.Currency.String()
	item.CreatedAt = model.CreatedAt.Unix()
	item.UpdatedAt = model.UpdatedAt.Unix()
}

func (item transactionRecordDTO) ToModel() model.TransactionRecord {
	modelItem := model.TransactionRecord{
		ID: model.TransactionRecordID(item.Id),
		HasBankAccountID: model.HasBankAccountID{
			BankAccountID: model.BankAccountID(item.BankAccountId),
		},
		HasBankUserID: model.HasBankUserID{
			BankUserID: model.BankUserID(item.BankUserId),
		},
		Type: model.TransactionType(item.RecordType),
		DepositAmount: model.Money{
			Amount:   item.DepositAmountNumber,
			Currency: model.Currency(item.DepositAmountCurrency),
		},
		WithdrawAmount: model.Money{
			Amount:   item.WithdrawAmountNumber,
			Currency: model.Currency(item.WithdrawAmountCurrency),
		},
		AfterAmount: model.Money{
			Amount:   item.AfterAmountNumber,
			Currency: model.Currency(item.AfterAmountCurrency),
		},
	}
	modelItem.CreatedAt = timer.TimeFromInt64(item.CreatedAt)
	modelItem.UpdatedAt = timer.TimeFromInt64(item.UpdatedAt)
	return modelItem
}
