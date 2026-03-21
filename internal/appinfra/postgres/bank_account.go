package postgres

import (
	"database/sql"
)

type BankAccount struct {
	Id                      string
	BankUserId              string
	AccountType             string
	AmountNumber            int64
	AmountCurrency          string
	LastTransactionRecordId string
	LastTransactionTime     int64
	CreatedAt               int64
	UpdatedAt               int64
}

func (item *BankAccount) SetBySQLRow(row *sql.Row) error {
	return row.Scan(
		&item.Id,
		&item.BankUserId,
		&item.AccountType,
		&item.AmountNumber,
		&item.AmountCurrency,
		&item.LastTransactionRecordId,
		&item.LastTransactionTime,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}

func (item BankAccount) ToMap() map[string]any {
	return map[string]any{
		ColumnID:                      item.Id,
		ColumnBankUserID:              item.BankUserId,
		ColumnType:                    item.AccountType,
		ColumnAmountNumber:            item.AmountNumber,
		ColumnAmountCurrency:          item.AmountCurrency,
		ColumnLastTransactionRecordID: item.LastTransactionRecordId,
		ColumnLastTransactionTime:     item.LastTransactionTime,
		ColumnCreatedAt:               item.CreatedAt,
		ColumnUpdatedAt:               item.UpdatedAt,
	}
}

func (item BankAccount) Columns() []string {
	return []string{
		ColumnID,
		ColumnBankUserID,
		ColumnType,
		ColumnAmountNumber,
		ColumnAmountCurrency,
		ColumnLastTransactionRecordID,
		ColumnLastTransactionTime,
		ColumnCreatedAt,
		ColumnUpdatedAt,
	}
}

func (item BankAccount) Values() []any {
	return []any{
		item.Id,
		item.BankUserId,
		item.AccountType,
		item.AmountNumber,
		item.AmountCurrency,
		item.LastTransactionRecordId,
		item.LastTransactionTime,
		item.CreatedAt,
		item.UpdatedAt,
	}
}
