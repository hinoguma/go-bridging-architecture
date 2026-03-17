package driven_infra_port

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/crosscutting/errors"
	"app/pkg/crosscutting/timer"
	"app/pkg/driven_infra/aws"
	"app/pkg/driven_infra/rdb"
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
	client rdb.SQLClient
}

func NewTransactionRecordRepositorySQL(client rdb.SQLClient) repository.TransactionRecordRepository {
	return transactionRecordRepositorySQL{
		client: client,
	}
}

func (repo transactionRecordRepositorySQL) Get(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error) {
	var row *sql.Row
	var err error

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := transactionRecordSQLItem{}

	row, err = repo.client.GetRowByID(ctx, rdb.TableTransactionRecords, id.String(), sqlItem.Columns(), sqlOptions)
	if err != nil {
		return model.TransactionRecord{}, errors.Lift(err)
	}

	err = sqlItem.SetBySQLRow(row)
	if err != nil {
		return model.TransactionRecord{}, errors.Lift(err)
	}
	return sqlItem.ToModel(), nil
}

func (repo transactionRecordRepositorySQL) Create(
	ctx context.Context,
	item model.TransactionRecord,
	optionaltFuncs ...model.DBOperationOptionalFunc,
) (model.TransactionRecord, error) {

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := transactionRecordSQLItem{}
	sqlItem.SetByModel(item)

	sqlRow, err := repo.client.CreateRow(ctx, rdb.TableTransactionRecords, &sqlItem, sqlOptions)
	if err != nil {
		return item, errors.Lift(err)
	}
	err = sqlItem.SetBySQLRow(sqlRow)
	if err != nil {
		return model.TransactionRecord{}, errors.Lift(err)
	}
	return sqlItem.ToModel(), nil
}

type transactionRecordSQLItem struct {
	id                     string `db:"id"`
	bankAccountId          string `db:"bank_account_id"`
	bankUserId             string `db:"bank_user_id"`
	recordType             string `db:"type"`
	depositAmountNumber    int64  `db:"deposit_amount_number"`
	depositAmountCurrency  string `db:"deposit_amount_currency"`
	withdrawAmountNumber   int64  `db:"withdraw_amount_number"`
	withdrawAmountCurrency string `db:"withdraw_amount_currency"`
	afterAmountNumber      int64  `db:"after_amount_number"`
	afterAmountCurrency    string `db:"after_amount_currency"`
	createdAt              int64  `db:"created_at"`
	updatedAt              int64  `db:"updated_at"`
}

func (item *transactionRecordSQLItem) SetByModel(model model.TransactionRecord) {
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

func (item *transactionRecordSQLItem) SetBySQLRow(row *sql.Row) error {
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

func (item transactionRecordSQLItem) ToModel() model.TransactionRecord {
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

func (item transactionRecordSQLItem) ToMap() map[string]any {
	return map[string]any{
		"id":                       item.id,
		"bank_account_id":          item.bankAccountId,
		"bank_user_id":             item.bankUserId,
		"type":                     item.recordType,
		"deposit_amount_number":    item.depositAmountNumber,
		"deposit_amount_currency":  item.depositAmountCurrency,
		"withdraw_amount_number":   item.withdrawAmountNumber,
		"withdraw_amount_currency": item.withdrawAmountCurrency,
		"after_amount_number":      item.afterAmountNumber,
		"after_amount_currency":    item.afterAmountCurrency,
		"created_at":               item.createdAt,
		"updated_at":               item.updatedAt,
	}
}

func (item transactionRecordSQLItem) Columns() []string {
	return []string{
		"id",
		"bank_account_id",
		"bank_user_id",
		"type",
		"deposit_amount_number",
		"deposit_amount_currency",
		"withdraw_amount_number",
		"withdraw_amount_currency",
		"after_amount_number",
		"after_amount_currency",
		"created_at",
		"updated_at",
	}
}

func (item transactionRecordSQLItem) Values() []any {
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
