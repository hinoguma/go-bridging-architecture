package postgres

import (
	"database/sql"
)

type TransactionRecord struct {
	Id                     string
	BankAccountId          string
	BankUserId             string
	RecordType             string
	DepositAmountNumber    float64
	DepositAmountCurrency  string
	WithdrawAmountNumber   float64
	WithdrawAmountCurrency string
	AfterAmountNumber      float64
	AfterAmountCurrency    string
	CreatedAt              int64
	UpdatedAt              int64
}

func (item *TransactionRecord) SetBySQLRow(row *sql.Row) error {
	return row.Scan(
		&item.Id,
		&item.BankAccountId,
		&item.BankUserId,
		&item.RecordType,
		&item.DepositAmountNumber,
		&item.DepositAmountCurrency,
		&item.WithdrawAmountNumber,
		&item.WithdrawAmountCurrency,
		&item.AfterAmountNumber,
		&item.AfterAmountCurrency,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}

func (item TransactionRecord) ToMap() map[string]any {
	return map[string]any{
		ColumnID:                     item.Id,
		ColumnBankAccountID:          item.BankAccountId,
		ColumnBankUserID:             item.BankUserId,
		ColumnType:                   item.RecordType,
		ColumnDepositAmountNumber:    item.DepositAmountNumber,
		ColumnDepositAmountCurrency:  item.DepositAmountCurrency,
		ColumnWithdrawAmountNumber:   item.WithdrawAmountNumber,
		ColumnWithdrawAmountCurrency: item.WithdrawAmountCurrency,
		ColumnAfterAmountNumber:      item.AfterAmountNumber,
		ColumnAfterAmountCurrency:    item.AfterAmountCurrency,
		ColumnCreatedAt:              item.CreatedAt,
		ColumnUpdatedAt:              item.UpdatedAt,
	}
}

func (item TransactionRecord) Columns() []string {
	return []string{
		ColumnID,
		ColumnBankAccountID,
		ColumnBankUserID,
		ColumnType,
		ColumnDepositAmountNumber,
		ColumnDepositAmountCurrency,
		ColumnWithdrawAmountNumber,
		ColumnWithdrawAmountCurrency,
		ColumnAfterAmountNumber,
		ColumnAfterAmountCurrency,
		ColumnCreatedAt,
		ColumnUpdatedAt,
	}
}

func (item TransactionRecord) Values() []any {
	return []any{
		item.Id,
		item.BankAccountId,
		item.BankUserId,
		item.RecordType,
		item.DepositAmountNumber,
		item.DepositAmountCurrency,
		item.WithdrawAmountNumber,
		item.WithdrawAmountCurrency,
		item.AfterAmountNumber,
		item.AfterAmountCurrency,
		item.CreatedAt,
		item.UpdatedAt,
	}
}
