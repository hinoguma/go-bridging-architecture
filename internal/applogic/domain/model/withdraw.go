package model

import "time"

const MaxWithdrawAmountOneTime = 1_000_000

type WithdrawServiceRequest struct {
	UseDBTransaction
	HasRequestAt
	BankAccountID       BankAccountID
	TransactionRecordID *TransactionRecordID
	Amount              Money
}

func (model WithdrawServiceRequest) HasTransactionRecordID() bool {
	return model.TransactionRecordID != nil
}

func (model WithdrawServiceRequest) GetTransactionRecordID() TransactionRecordID {
	if model.TransactionRecordID == nil {
		return ""
	}
	return *model.TransactionRecordID
}

type WithdrawServiceResult struct {
	BankAccount      BankAccount
	Record           TransactionRecord
	TxID             DBTransactionID
	NotEnoughBalance bool
}

type WithdrawResult struct {
	BankAccount              BankAccount
	UpdateBankAccountRequest UpdateBankAccountRequest
	TransactionRecord        TransactionRecord
	NotEnoughBalance         bool
}

func NewWithdrawResultNotEnoughBalance() WithdrawResult {
	return WithdrawResult{NotEnoughBalance: true}
}

func NewWithdrawResult(
	bankAccount BankAccount,
	updateBankAccountRequest UpdateBankAccountRequest,
	transactionRecord TransactionRecord,
) WithdrawResult {
	return WithdrawResult{
		BankAccount:              bankAccount,
		UpdateBankAccountRequest: updateBankAccountRequest,
		TransactionRecord:        transactionRecord,
		NotEnoughBalance:         false,
	}
}

func Withdraw(recordId TransactionRecordID, account BankAccount, amount Money, t time.Time) WithdrawResult {

	if !account.Amount.GreaterThan(amount) {
		return NewWithdrawResultNotEnoughBalance()
	}

	updateReq := NewUpdateBankAccountRequest(account.ID)
	updateReq.SetAmount(account.Amount).
		SetLastTransactionRecordID(recordId).
		SetLastTransactionTime(t)

	account = updateReq.UpdateItem(account)
	record := NewTransactionRecordWithdraw(recordId, account, amount, t)

	return NewWithdrawResult(account, updateReq, record)
}
