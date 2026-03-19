package applogic_driven_infra_bridge

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/timer"
	"app/internal/appinfra/aws"
	"app/internal/appinfra/postgres"
	"context"
	"database/sql"
)

type transactionRecordRepositoryDynamoDB struct {
	dynamo driven_infra_aws.DynamoDBClient
}

func (repo transactionRecordRepositoryDynamoDB) Get(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (repo transactionRecordRepositoryDynamoDB) Create(ctx context.Context, account model.TransactionRecord, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
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

type transactionRecordRepositorySQL struct {
	client postgres.SQLClient
}

func NewTransactionRecordRepositorySQL(client postgres.SQLClient) repository.TransactionRecordRepository {
	return transactionRecordRepositorySQL{
		client: client,
	}
}

func (repo transactionRecordRepositorySQL) Get(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	var row *sql.Row
	var err error

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

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

func (repo transactionRecordRepositorySQL) Create(
	ctx context.Context,
	item model.TransactionRecord,
	optionaltFuncs ...model.DBOperationOptionalFunc,
) (model.TransactionRecord, error) {

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

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
	id                     string
	bankAccountId          string
	bankUserId             string
	recordType             string
	depositAmountNumber    int64
	depositAmountCurrency  string
	withdrawAmountNumber   int64
	withdrawAmountCurrency string
	afterAmountNumber      int64
	afterAmountCurrency    string
	createdAt              int64
	updatedAt              int64
}

func (item *transactionRecordDTO) SetByModel(model model.TransactionRecord) {
	item.id = model.ID.String()
	item.bankAccountId = model.BankAccountID.String()
	item.bankUserId = model.BankUserID.String()
	item.recordType = model.Type.String()
	item.depositAmountNumber = model.DepositAmount.Amount
	item.depositAmountCurrency = model.DepositAmount.Currency.String()
	item.withdrawAmountNumber = model.WithdrawAmount.Amount
	item.withdrawAmountCurrency = model.WithdrawAmount.Currency.String()
	item.afterAmountNumber = model.AfterAmount.Amount
	item.afterAmountCurrency = model.AfterAmount.Currency.String()
	item.createdAt = model.CreatedAt.Unix()
	item.updatedAt = model.UpdatedAt.Unix()
}

func (item *transactionRecordDTO) SetBySQLRow(row *sql.Row) error {
	return row.Scan(
		&item.id,
		&item.bankAccountId,
		&item.bankUserId,
		&item.recordType,
		&item.depositAmountNumber,
		&item.depositAmountCurrency,
		&item.withdrawAmountNumber,
		&item.withdrawAmountCurrency,
		&item.afterAmountNumber,
		&item.afterAmountCurrency,
		&item.createdAt,
		&item.updatedAt,
	)
}

func (item transactionRecordDTO) ToModel() model.TransactionRecord {
	modelItem := model.TransactionRecord{
		ID: model.TransactionRecordID(item.id),
		HasBankAccountID: model.HasBankAccountID{
			BankAccountID: model.BankAccountID(item.bankAccountId),
		},
		HasBankUserID: model.HasBankUserID{
			BankUserID: model.BankUserID(item.bankUserId),
		},
		Type: model.TransactionType(item.recordType),
		DepositAmount: model.Money{
			Amount:   item.depositAmountNumber,
			Currency: model.Currency(item.depositAmountCurrency),
		},
		WithdrawAmount: model.Money{
			Amount:   item.withdrawAmountNumber,
			Currency: model.Currency(item.withdrawAmountCurrency),
		},
		AfterAmount: model.Money{
			Amount:   item.afterAmountNumber,
			Currency: model.Currency(item.afterAmountCurrency),
		},
	}
	modelItem.CreatedAt = timer.TimeFromInt64(item.createdAt)
	modelItem.UpdatedAt = timer.TimeFromInt64(item.updatedAt)
	return modelItem
}

func (item transactionRecordDTO) ToMap() map[string]any {
	return map[string]any{
		postgres.ColumnID:                     item.id,
		postgres.ColumnBankAccountID:          item.bankAccountId,
		postgres.ColumnBankUserID:             item.bankUserId,
		postgres.ColumnType:                   item.recordType,
		postgres.ColumnDepositAmountNumber:    item.depositAmountNumber,
		postgres.ColumnDepositAmountCurrency:  item.depositAmountCurrency,
		postgres.ColumnWithdrawAmountNumber:   item.withdrawAmountNumber,
		postgres.ColumnWithdrawAmountCurrency: item.withdrawAmountCurrency,
		postgres.ColumnAfterAmountNumber:      item.afterAmountNumber,
		postgres.ColumnAfterAmountCurrency:    item.afterAmountCurrency,
		postgres.ColumnCreatedAt:              item.createdAt,
		postgres.ColumnUpdatedAt:              item.updatedAt,
	}
}

func (item transactionRecordDTO) Columns() []string {
	return []string{
		postgres.ColumnID,
		postgres.ColumnBankAccountID,
		postgres.ColumnBankUserID,
		postgres.ColumnType,
		postgres.ColumnDepositAmountNumber,
		postgres.ColumnDepositAmountCurrency,
		postgres.ColumnWithdrawAmountNumber,
		postgres.ColumnWithdrawAmountCurrency,
		postgres.ColumnAfterAmountNumber,
		postgres.ColumnAfterAmountCurrency,
		postgres.ColumnCreatedAt,
		postgres.ColumnUpdatedAt,
	}
}

func (item transactionRecordDTO) Values() []any {
	return []any{
		item.id,
		item.bankAccountId,
		item.bankUserId,
		item.recordType,
		item.depositAmountNumber,
		item.depositAmountCurrency,
		item.withdrawAmountNumber,
		item.withdrawAmountCurrency,
		item.afterAmountNumber,
		item.afterAmountCurrency,
		item.createdAt,
		item.updatedAt,
	}
}
