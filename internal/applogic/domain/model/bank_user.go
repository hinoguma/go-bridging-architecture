package model

import "app/internal/crosscutting"

type BankUserID string

func (value BankUserID) String() string {
	return string(value)
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
