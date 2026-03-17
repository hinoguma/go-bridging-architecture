package model

import "time"

type DepositRequest struct {
	UseDBTransaction
	HasRequestAt
	BankAccountID BankAccountID
	Amount        Money
}

type DepositResult struct {
	BankAccount BankAccount
	Record      TransactionRecord
	TxID        DBTransactionID
}

func Deposit(account BankAccount, amount Money, t time.Time) (BankAccount, UpdateBankAccountRequest, TransactionRecord) {
	recordId := IssueTransactionRecordID()

	updateReq := NewUpdateBankAccountRequest(account.ID)
	updateReq.SetAmount(account.Amount).
		SetLastTransactionRecordID(recordId).
		SetLastTransactionTime(t)

	account = updateReq.UpdateItem(account)
	record := NewTransactionRecordDeposit(account, amount, t)
	return account, updateReq, record
}
