package model

import (
	"app/internal/crosscutting"
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

type UpdateBankAccountRequest struct {
	ID                      BankAccountID
	Amount                  *Money
	LastTransactionRecordID *TransactionRecordID
	LastTransactionTime     *time.Time
}

func NewUpdateBankAccountRequest(id BankAccountID) UpdateBankAccountRequest {
	return UpdateBankAccountRequest{
		ID: id,
	}
}

func (model *UpdateBankAccountRequest) SetAmount(amount Money) *UpdateBankAccountRequest {
	model.Amount = &amount
	return model
}

func (model *UpdateBankAccountRequest) SetLastTransactionRecordID(recordId TransactionRecordID) *UpdateBankAccountRequest {
	model.LastTransactionRecordID = &recordId
	return model
}

func (model *UpdateBankAccountRequest) SetLastTransactionTime(t time.Time) *UpdateBankAccountRequest {
	model.LastTransactionTime = &t
	return model
}

func (model UpdateBankAccountRequest) UpdateItem(item BankAccount) BankAccount {
	if item.ID != model.ID {
		return item
	}
	if model.Amount != nil {
		item.Amount = *model.Amount
	}
	if model.LastTransactionRecordID != nil {
		item.LastTransactionRecordID = *model.LastTransactionRecordID
	}
	if model.LastTransactionTime != nil {
		item.LastTransactionTime = *model.LastTransactionTime
	}
	return item
}
