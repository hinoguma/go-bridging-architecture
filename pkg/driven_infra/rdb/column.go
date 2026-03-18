package rdb

/**
 * Shared column names for all tables.
 */

const (
	ColumnID             string = "id"
	ColumnType           string = "type"
	ColumnBankUserID     string = "bank_user_id"
	ColumnAmountNumber   string = "amount_number"
	ColumnAmountCurrency string = "amount_currency"
	ColumnCreatedAt      string = "created_at"
	ColumnUpdatedAt      string = "updated_at"
)

/**
 * bank_accounts table column names.
 */

const (
	ColumnLastTransactionRecordID string = "last_transaction_record_id"
	ColumnLastTransactionTime     string = "last_transaction_time"
)
