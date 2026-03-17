package model

import (
	"app/pkg/crosscutting"
	"time"
)

type BankAccountID string

func (value BankAccountID) String() string {
	return string(value)
}

func IssueBankAccountID() BankAccountID {
	return BankAccountID(crosscutting.IssueRandomStrID())
}

type BankAccountType string

func (value BankAccountType) String() string {
	return string(value)
}

const (
	BankAccountTypeCurrent BankAccountType = "current"
	BankAccountTypeSavings BankAccountType = "savings"
)

type BankAccount struct {
	ID BankAccountID
	HasBankUserID
	Type                    BankAccountType
	Amount                  Money
	LastTransactionRecordID TransactionRecordID
	LastTransactionTime     time.Time
	DBItem
}

type HasBankAccountID struct {
	BankAccountID BankAccountID
}
