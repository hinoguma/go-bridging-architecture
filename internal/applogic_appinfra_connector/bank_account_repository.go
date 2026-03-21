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
	client postgres.SQLClient
}

func (repo bankAccountRepositorySQL) Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error) {
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

func (repo bankAccountRepositorySQL) Create(
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

func (repo bankAccountRepositorySQL) Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error {
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

func NewBankAccountRepositorySQL(client postgres.SQLClient) repository.BankAccountRepository {
	return bankAccountRepositorySQL{
		client: client,
	}
}

type bankAccountDTO struct {
	id                      string
	bankUserId              string
	accountType             string
	amountNumber            int64
	amountCurrency          string
	lastTransactionRecordId string
	lastTransactionTime     int64
	createdAt               int64
	updatedAt               int64
}

func (item *bankAccountDTO) SetByModel(model model.BankAccount) {
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

func (item *bankAccountDTO) SetBySQLRow(row *sql.Row) error {
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

func (item bankAccountDTO) ToModel() model.BankAccount {
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

func (item bankAccountDTO) ToMap() map[string]any {
	return map[string]any{
		postgres.ColumnID:                      item.id,
		postgres.ColumnBankUserID:              item.bankUserId,
		postgres.ColumnType:                    item.accountType,
		postgres.ColumnAmountNumber:            item.amountNumber,
		postgres.ColumnAmountCurrency:          item.amountCurrency,
		postgres.ColumnLastTransactionRecordID: item.lastTransactionRecordId,
		postgres.ColumnLastTransactionTime:     item.lastTransactionTime,
		postgres.ColumnCreatedAt:               item.createdAt,
		postgres.ColumnUpdatedAt:               item.updatedAt,
	}
}

func (item bankAccountDTO) Columns() []string {
	return []string{
		postgres.ColumnID,
		postgres.ColumnBankUserID,
		postgres.ColumnType,
		postgres.ColumnAmountNumber,
		postgres.ColumnAmountCurrency,
		postgres.ColumnLastTransactionRecordID,
		postgres.ColumnLastTransactionTime,
		postgres.ColumnCreatedAt,
		postgres.ColumnUpdatedAt,
	}
}

func (item bankAccountDTO) Values() []any {
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
