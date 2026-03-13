package model

import "app/pkg/crosscutting"

type TransactionRecordID string

func IssueTransactionRecordID() TransactionRecordID {
	return TransactionRecordID(crosscutting.IssueRandomStrID())
}

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

type TransactionRecord struct {
	ID TransactionRecordID
	HasBankUserID
	HasBankAccountID
	Type           TransactionType
	DepositAmount  Money
	WithdrawAmount Money
	AfterAmount    Money
	DBItem
}
