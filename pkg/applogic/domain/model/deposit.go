package model

import "time"

type DepositRequest struct {
	UseDBTransaction
	HasRequestAt
	BankAccountID BankAccountID
	Amount        Money
}

type DepositResult struct {
	Record TransactionRecord
	TxID   DBTransactionID
}

func Deposit(account BankAccount, amount Money, t time.Time) (BankAccount, TransactionRecord) {
	recordId := IssueTransactionRecordID()
	account.LastTransactionRecordID = recordId
	account.LastTransactionTime = t
	record := NewTransactionRecordDeposit(account, amount, t)
	return account, record
}
