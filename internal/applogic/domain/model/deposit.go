package model

import "time"

type DepositServiceRequest struct {
	UseDBTransaction
	HasRequestAt
	BankAccountID       BankAccountID
	TransactionRecordID *TransactionRecordID
	Amount              Money
}

func (model DepositServiceRequest) HasTransactionRecordID() bool {
	return model.TransactionRecordID != nil
}

func (model DepositServiceRequest) GetTransactionRecordID() TransactionRecordID {
	if model.TransactionRecordID == nil {
		return ""
	}
	return *model.TransactionRecordID
}

type DepositServiceResult struct {
	BankAccount BankAccount
	Record      TransactionRecord
	TxID        DBTransactionID
}

type DepositResult struct {
	BankAccount              BankAccount
	UpdateBankAccountRequest UpdateBankAccountRequest
	TransactionRecord        TransactionRecord
}

func NewDepositResult(
	bankAccount BankAccount,
	updateBankAccountRequest UpdateBankAccountRequest,
	transactionRecord TransactionRecord,
) DepositResult {
	return DepositResult{
		BankAccount:              bankAccount,
		UpdateBankAccountRequest: updateBankAccountRequest,
		TransactionRecord:        transactionRecord,
	}
}

func Deposit(recordId TransactionRecordID, account BankAccount, amount Money, t time.Time) DepositResult {

	updateReq := NewUpdateBankAccountRequest(account.ID)
	updateReq.SetAmount(account.Amount).
		SetLastTransactionRecordID(recordId).
		SetLastTransactionTime(t)

	account = updateReq.UpdateItem(account)
	record := NewTransactionRecordDeposit(recordId, account, amount, t)

	return NewDepositResult(account, updateReq, record)
}
