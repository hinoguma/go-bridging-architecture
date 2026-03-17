package driven_infra_port

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/crosscutting/errors"
	"app/pkg/crosscutting/timer"
	driven_infra_aws "app/pkg/driven_infra/aws"
	"app/pkg/driven_infra/rdb"
	"context"
	"database/sql"
)

type bankAccountRepositoryDynamoDB struct {
	dynamo driven_infra_aws.DynamoDBClient
}

func (b bankAccountRepositoryDynamoDB) Create(ctx context.Context, account model.BankAccount, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDB) Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error {
	//TODO implement me
	panic("implement me")
}

func (b bankAccountRepositoryDynamoDB) Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
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

type bankAccountRepositorySQL struct {
	client rdb.SQLClient
}

func (repo bankAccountRepositorySQL) Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
	var row *sql.Row
	var err error

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := bankAccountSQLItem{}

	row, err = repo.client.GetRowByID(ctx, rdb.TableBankAccounts, id.String(), sqlItem.Columns(), sqlOptions)
	if err != nil {
		return model.BankAccount{}, errors.Lift(err)
	}

	err = sqlItem.SetBySQLRow(row)
	if err != nil {
		return model.BankAccount{}, errors.Lift(err)
	}
	return sqlItem.ToModel(), nil
}

func (repo bankAccountRepositorySQL) Create(
	ctx context.Context,
	item model.BankAccount,
	optionaltFuncs ...model.DBOperationOptionalFunc,
) (model.BankAccount, error) {

	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	sqlItem := bankAccountSQLItem{}
	sqlItem.SetByModel(item)

	sqlRow, err := repo.client.CreateRow(ctx, rdb.TableBankAccounts, &sqlItem, sqlOptions)
	if err != nil {
		return item, errors.Lift(err)
	}
	err = sqlItem.SetBySQLRow(sqlRow)
	if err != nil {
		return model.BankAccount{}, errors.Lift(err)
	}
	return sqlItem.ToModel(), nil
}

func (repo bankAccountRepositorySQL) Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error {
	sqlOptions := convertOptionalFuncsToSQLOperationOptions(optionaltFuncs)

	updateRequests := convertUpdateBankAccountRequestToUpdateRequests(request)
	_, err := repo.client.UpdateRowByStrID(ctx, rdb.TableBankAccounts, request.ID.String(), updateRequests, sqlOptions)
	if err != nil {
		return errors.Lift(err)
	}
	return nil
}

func NewBankAccountRepositorySQL(client rdb.SQLClient) repository.BankAccountRepository {
	return bankAccountRepositorySQL{
		client: client,
	}
}

type bankAccountSQLItem struct {
	id                      string `db:"id"`
	bankUserId              string `db:"bank_user_id"`
	accountType             string `db:"type"`
	amountNumber            int64  `db:"amount_number"`
	amountCurrency          string `db:"amount_currency"`
	lastTransactionRecordId string `db:"last_transaction_record_id"`
	lastTransactionTime     int64  `db:"last_transaction_time"`
	createdAt               int64  `db:"created_at"`
	updatedAt               int64  `db:"updated_at"`
}

func (item *bankAccountSQLItem) SetByModel(model model.BankAccount) {
	item.id = model.ID.String()
	item.bankUserId = model.BankUserID.String()
	item.accountType = model.Type.String()
	item.amountNumber = model.Amount.Amount
	item.amountCurrency = model.Amount.Currency.String()
	item.lastTransactionRecordId = model.LastTransactionRecordID.String()
	item.lastTransactionTime = model.LastTransactionTime.Unix()
	item.createdAt = model.CreatedAt.Unix()
	item.updatedAt = model.UpdatedAt.Unix()
}

func (item *bankAccountSQLItem) SetBySQLRow(row *sql.Row) error {
	return row.Scan(
		&item.id,
		&item.bankUserId,
		&item.accountType,
		&item.amountNumber,
		&item.amountCurrency,
		&item.lastTransactionRecordId,
		&item.lastTransactionTime,
		&item.createdAt,
		&item.updatedAt,
	)
}

func (item bankAccountSQLItem) ToModel() model.BankAccount {
	modelItem := model.BankAccount{
		ID: model.BankAccountID(item.id),
		HasBankUserID: model.HasBankUserID{
			BankUserID: model.BankUserID(item.bankUserId),
		},
		Type: model.BankAccountType(item.accountType),
		Amount: model.Money{
			Amount:   item.amountNumber,
			Currency: model.Currency(item.amountCurrency),
		},
		LastTransactionRecordID: model.TransactionRecordID(item.lastTransactionRecordId),
		LastTransactionTime:     timer.TimeFromInt64(item.lastTransactionTime),
	}
	modelItem.CreatedAt = timer.TimeFromInt64(item.createdAt)
	modelItem.UpdatedAt = timer.TimeFromInt64(item.updatedAt)
	return modelItem
}

func (item bankAccountSQLItem) ToMap() map[string]any {
	return map[string]any{
		"id":                         item.id,
		"bank_user_id":               item.bankUserId,
		"type":                       item.accountType,
		"amount_number":              item.amountNumber,
		"amount_currency":            item.amountCurrency,
		"last_transaction_record_id": item.lastTransactionRecordId,
		"last_transaction_time":      item.lastTransactionTime,
		"created_at":                 item.createdAt,
		"updated_at":                 item.updatedAt,
	}
}

func (item bankAccountSQLItem) Columns() []string {
	return []string{
		"id",
		"bank_user_id",
		"type",
		"amount_number",
		"amount_currency",
		"last_transaction_record_id",
		"last_transaction_time",
		"created_at",
		"updated_at",
	}
}

func (item bankAccountSQLItem) Values() []any {
	return []any{
		item.id,
		item.bankUserId,
		item.accountType,
		item.amountNumber,
		item.amountCurrency,
		item.lastTransactionRecordId,
		item.lastTransactionTime,
		item.createdAt,
		item.updatedAt,
	}
}

func convertUpdateBankAccountRequestToUpdateRequests(request model.UpdateBankAccountRequest) rdb.UpdateFieldRequests {
	var updateRequests rdb.UpdateFieldRequests = make([]rdb.UpdateFieldRequest, 0)
	if request.Amount != nil {
		updateRequests.Append("amount_number", request.Amount.Amount)
		updateRequests.Append("amount_currency", request.Amount.Currency.String())
	}
	if request.LastTransactionRecordID != nil {
		updateRequests.Append("last_transaction_record_id", request.LastTransactionRecordID.String())
	}
	if request.LastTransactionTime != nil {
		updateRequests.Append("last_transaction_time", request.LastTransactionTime.Unix())
	}
	return updateRequests
}
