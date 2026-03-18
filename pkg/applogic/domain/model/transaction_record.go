package model

import (
	"app/pkg/crosscutting"
	"time"
)

type TransactionRecordID string

func (value TransactionRecordID) String() string {
	return string(value)
}

func IssueTransactionRecordID() TransactionRecordID {
	return TransactionRecordID(crosscutting.IssueRandomStrID())
}

type TransactionType string

func (value TransactionType) String() string {
	return string(value)
}

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

func NewTransactionRecordDeposit(
	accountAfter BankAccount, amount Money, t time.Time,
) TransactionRecord {
	record := TransactionRecord{
		ID:               IssueTransactionRecordID(),
		HasBankUserID:    HasBankUserID{BankUserID: accountAfter.BankUserID},
		HasBankAccountID: HasBankAccountID{BankAccountID: accountAfter.ID},
		Type:             TransactionTypeDeposit,
		DepositAmount:    amount,
		AfterAmount:      accountAfter.Amount,
	}
	record.Create(t)
	return record
}

func NewTransactionRecordWithdraw(
	accountAfter BankAccount, amount Money, t time.Time,
) TransactionRecord {
	record := TransactionRecord{
		ID:               IssueTransactionRecordID(),
		HasBankUserID:    HasBankUserID{BankUserID: accountAfter.BankUserID},
		HasBankAccountID: HasBankAccountID{BankAccountID: accountAfter.ID},
		Type:             TransactionTypeWithdrawal,
		WithdrawAmount:   amount,
		AfterAmount:      accountAfter.Amount,
	}
	record.Create(t)
	return record
}
