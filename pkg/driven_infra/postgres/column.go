package postgres

/**
 * Shared column names for all tables.
 */

const (
	ColumnID             string = "id"
	ColumnType           string = "type"
	ColumnBankAccountID  string = "bank_account_id"
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

/**
 * transaction_records table column names.
 */

const (
	ColumnDepositAmountNumber    string = "deposit_amount_number"
	ColumnDepositAmountCurrency  string = "deposit_amount_currency"
	ColumnWithdrawAmountNumber   string = "withdraw_amount_number"
	ColumnWithdrawAmountCurrency string = "withdraw_amount_currency"
	ColumnAfterAmountNumber      string = "after_amount_number"
	ColumnAfterAmountCurrency    string = "after_amount_currency"
)
