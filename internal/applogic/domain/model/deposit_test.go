package model

import (
	"reflect"
	"testing"
	"time"
)

func TestDeposit(t *testing.T) {
	baseTime := time.Date(2026, 3, 19, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		recordId       TransactionRecordID
		bankAccount    BankAccount
		amount         Money
		time           time.Time
		expectedResult DepositResult
	}{
		// Normal cases
		{
			name:     "deposit 100 USD to account with 1000 USD",
			recordId: "record-001",
			bankAccount: BankAccount{
				ID:            "account-001",
				HasBankUserID: HasBankUserID{BankUserID: "user-001"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 1000, Currency: USD},
			},
			amount: Money{Amount: 100, Currency: USD},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-001",
					HasBankUserID: HasBankUserID{BankUserID: "user-001"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 1000, Currency: USD},
					LastTransactionRecordID: "record-001",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-001",
					Amount:                  ptrMoney(Money{Amount: 1000, Currency: USD}),
					LastTransactionRecordID: ptrTxRecordID("record-001"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-001",
					HasBankUserID:    HasBankUserID{BankUserID: "user-001"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-001"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 100, Currency: USD},
					AfterAmount:      Money{Amount: 1000, Currency: USD},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
		{
			name:     "deposit 1 JPY to empty account",
			recordId: "record-002",
			bankAccount: BankAccount{
				ID:            "account-002",
				HasBankUserID: HasBankUserID{BankUserID: "user-002"},
				Type:          BankAccountTypeCurrent,
				Amount:        Money{Amount: 0, Currency: JPY},
			},
			amount: Money{Amount: 1, Currency: JPY},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-002",
					HasBankUserID: HasBankUserID{BankUserID: "user-002"},
					Type:          BankAccountTypeCurrent,
					Amount:        Money{Amount: 0, Currency: JPY},
					LastTransactionRecordID: "record-002",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-002",
					Amount:                  ptrMoney(Money{Amount: 0, Currency: JPY}),
					LastTransactionRecordID: ptrTxRecordID("record-002"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-002",
					HasBankUserID:    HasBankUserID{BankUserID: "user-002"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-002"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 1, Currency: JPY},
					AfterAmount:      Money{Amount: 0, Currency: JPY},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
		{
			name:     "deposit large amount",
			recordId: "record-003",
			bankAccount: BankAccount{
				ID:            "account-003",
				HasBankUserID: HasBankUserID{BankUserID: "user-003"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 999999999, Currency: EUR},
			},
			amount: Money{Amount: 1, Currency: EUR},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-003",
					HasBankUserID: HasBankUserID{BankUserID: "user-003"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 999999999, Currency: EUR},
					LastTransactionRecordID: "record-003",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-003",
					Amount:                  ptrMoney(Money{Amount: 999999999, Currency: EUR}),
					LastTransactionRecordID: ptrTxRecordID("record-003"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-003",
					HasBankUserID:    HasBankUserID{BankUserID: "user-003"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-003"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 1, Currency: EUR},
					AfterAmount:      Money{Amount: 999999999, Currency: EUR},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
		// Edge cases
		{
			name:     "deposit zero amount",
			recordId: "record-004",
			bankAccount: BankAccount{
				ID:            "account-004",
				HasBankUserID: HasBankUserID{BankUserID: "user-004"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 500, Currency: GBP},
			},
			amount: Money{Amount: 0, Currency: GBP},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-004",
					HasBankUserID: HasBankUserID{BankUserID: "user-004"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 500, Currency: GBP},
					LastTransactionRecordID: "record-004",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-004",
					Amount:                  ptrMoney(Money{Amount: 500, Currency: GBP}),
					LastTransactionRecordID: ptrTxRecordID("record-004"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-004",
					HasBankUserID:    HasBankUserID{BankUserID: "user-004"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-004"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 0, Currency: GBP},
					AfterAmount:      Money{Amount: 500, Currency: GBP},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
		{
			name:     "deposit with different currency than account",
			recordId: "record-005",
			bankAccount: BankAccount{
				ID:            "account-005",
				HasBankUserID: HasBankUserID{BankUserID: "user-005"},
				Type:          BankAccountTypeCurrent,
				Amount:        Money{Amount: 1000, Currency: USD},
			},
			amount: Money{Amount: 100, Currency: EUR},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-005",
					HasBankUserID: HasBankUserID{BankUserID: "user-005"},
					Type:          BankAccountTypeCurrent,
					Amount:        Money{Amount: 1000, Currency: USD},
					LastTransactionRecordID: "record-005",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-005",
					Amount:                  ptrMoney(Money{Amount: 1000, Currency: USD}),
					LastTransactionRecordID: ptrTxRecordID("record-005"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-005",
					HasBankUserID:    HasBankUserID{BankUserID: "user-005"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-005"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 100, Currency: EUR},
					AfterAmount:      Money{Amount: 1000, Currency: USD},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
		{
			name:     "deposit updates last transaction time",
			recordId: "record-006",
			bankAccount: BankAccount{
				ID:            "account-006",
				HasBankUserID: HasBankUserID{BankUserID: "user-006"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 100, Currency: USD},
				LastTransactionRecordID: "old-record",
				LastTransactionTime:     baseTime.Add(-24 * time.Hour),
			},
			amount: Money{Amount: 50, Currency: USD},
			time:   baseTime,
			expectedResult: DepositResult{
				BankAccount: BankAccount{
					ID:            "account-006",
					HasBankUserID: HasBankUserID{BankUserID: "user-006"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 100, Currency: USD},
					LastTransactionRecordID: "record-006",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-006",
					Amount:                  ptrMoney(Money{Amount: 100, Currency: USD}),
					LastTransactionRecordID: ptrTxRecordID("record-006"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-006",
					HasBankUserID:    HasBankUserID{BankUserID: "user-006"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-006"},
					Type:             TransactionTypeDeposit,
					DepositAmount:    Money{Amount: 50, Currency: USD},
					AfterAmount:      Money{Amount: 100, Currency: USD},
					DBItem:           DBItem{HasCreatedAt: HasCreatedAt{CreatedAt: baseTime}, HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Deposit(
				tc.recordId,
				tc.bankAccount,
				tc.amount,
				tc.time,
			)

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf(
					"Deposit() mismatch\ngot:  %+v\nwant: %+v",
					result,
					tc.expectedResult,
				)
			}
		})
	}
}

// Helper functions for pointers
func ptrMoney(m Money) *Money {
	return &m
}

func ptrTxRecordID(id TransactionRecordID) *TransactionRecordID {
	return &id
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
