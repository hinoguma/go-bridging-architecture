package model

import "app/pkg/crosscutting"

type BankUserID string

func (value BankUserID) String() string {
	return string(value)
}

func IssueBankUserID() BankUserID {
	return BankUserID(crosscutting.IssueRandomStrID())
}

type BankUser struct {
	ID    BankUserID
	Name  string
	Email crosscutting.Email
}

type HasBankUserID struct {
	BankUserID BankUserID
	DBItem
}
