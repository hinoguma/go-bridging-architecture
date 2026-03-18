package model

import (
	"reflect"
	"testing"
	"time"
)

func TestWithdraw(t *testing.T) {
	baseTime := time.Date(2026, 3, 19, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		recordId       TransactionRecordID
		bankAccount    BankAccount
		amount         Money
		time           time.Time
		expectedResult WithdrawResult
	}{
		// Normal cases
		{
			name:     "withdraw 100 USD from account with 1000 USD",
			recordId: "record-001",
			bankAccount: BankAccount{
				ID:            "account-001",
				HasBankUserID: HasBankUserID{BankUserID: "user-001"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 1000, Currency: USD},
			},
			amount: Money{Amount: 100, Currency: USD},
			time:   baseTime,
			expectedResult: WithdrawResult{
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
					Type:             TransactionTypeWithdrawal,
					WithdrawAmount:   Money{Amount: 100, Currency: USD},
					AfterAmount:      Money{Amount: 1000, Currency: USD},
					DBItem: DBItem{
						HasCreatedAt: HasCreatedAt{CreatedAt: baseTime},
						HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime},
					},
				},
				NotEnoughBalance: false,
			},
		},
		{
			name:     "withdraw 1 JPY from account with 100 JPY",
			recordId: "record-002",
			bankAccount: BankAccount{
				ID:            "account-002",
				HasBankUserID: HasBankUserID{BankUserID: "user-002"},
				Type:          BankAccountTypeCurrent,
				Amount:        Money{Amount: 100, Currency: JPY},
			},
			amount: Money{Amount: 1, Currency: JPY},
			time:   baseTime,
			expectedResult: WithdrawResult{
				BankAccount: BankAccount{
					ID:            "account-002",
					HasBankUserID: HasBankUserID{BankUserID: "user-002"},
					Type:          BankAccountTypeCurrent,
					Amount:        Money{Amount: 100, Currency: JPY},
					LastTransactionRecordID: "record-002",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-002",
					Amount:                  ptrMoney(Money{Amount: 100, Currency: JPY}),
					LastTransactionRecordID: ptrTxRecordID("record-002"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-002",
					HasBankUserID:    HasBankUserID{BankUserID: "user-002"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-002"},
					Type:             TransactionTypeWithdrawal,
					WithdrawAmount:   Money{Amount: 1, Currency: JPY},
					AfterAmount:      Money{Amount: 100, Currency: JPY},
					DBItem: DBItem{
						HasCreatedAt: HasCreatedAt{CreatedAt: baseTime},
						HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime},
					},
				},
				NotEnoughBalance: false,
			},
		},
		{
			name:     "withdraw updates last transaction time",
			recordId: "record-003",
			bankAccount: BankAccount{
				ID:            "account-003",
				HasBankUserID: HasBankUserID{BankUserID: "user-003"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 500, Currency: EUR},
				LastTransactionRecordID: "old-record",
				LastTransactionTime:     baseTime.Add(-24 * time.Hour),
			},
			amount: Money{Amount: 100, Currency: EUR},
			time:   baseTime,
			expectedResult: WithdrawResult{
				BankAccount: BankAccount{
					ID:            "account-003",
					HasBankUserID: HasBankUserID{BankUserID: "user-003"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 500, Currency: EUR},
					LastTransactionRecordID: "record-003",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-003",
					Amount:                  ptrMoney(Money{Amount: 500, Currency: EUR}),
					LastTransactionRecordID: ptrTxRecordID("record-003"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-003",
					HasBankUserID:    HasBankUserID{BankUserID: "user-003"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-003"},
					Type:             TransactionTypeWithdrawal,
					WithdrawAmount:   Money{Amount: 100, Currency: EUR},
					AfterAmount:      Money{Amount: 500, Currency: EUR},
					DBItem: DBItem{
						HasCreatedAt: HasCreatedAt{CreatedAt: baseTime},
						HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime},
					},
				},
				NotEnoughBalance: false,
			},
		},
		// Edge cases - Not enough balance
		{
			name:     "withdraw from empty account",
			recordId: "record-004",
			bankAccount: BankAccount{
				ID:            "account-004",
				HasBankUserID: HasBankUserID{BankUserID: "user-004"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 0, Currency: USD},
			},
			amount:         Money{Amount: 100, Currency: USD},
			time:           baseTime,
			expectedResult: WithdrawResult{NotEnoughBalance: true},
		},
		{
			name:     "withdraw more than balance",
			recordId: "record-005",
			bankAccount: BankAccount{
				ID:            "account-005",
				HasBankUserID: HasBankUserID{BankUserID: "user-005"},
				Type:          BankAccountTypeCurrent,
				Amount:        Money{Amount: 100, Currency: GBP},
			},
			amount:         Money{Amount: 200, Currency: GBP},
			time:           baseTime,
			expectedResult: WithdrawResult{NotEnoughBalance: true},
		},
		{
			name:     "withdraw exact balance (not allowed)",
			recordId: "record-006",
			bankAccount: BankAccount{
				ID:            "account-006",
				HasBankUserID: HasBankUserID{BankUserID: "user-006"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 500, Currency: USD},
			},
			amount:         Money{Amount: 500, Currency: USD},
			time:           baseTime,
			expectedResult: WithdrawResult{NotEnoughBalance: true},
		},
		// Edge cases - Zero amount
		{
			name:     "withdraw zero amount from empty account",
			recordId: "record-007",
			bankAccount: BankAccount{
				ID:            "account-007",
				HasBankUserID: HasBankUserID{BankUserID: "user-007"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 0, Currency: USD},
			},
			amount:         Money{Amount: 0, Currency: USD},
			time:           baseTime,
			expectedResult: WithdrawResult{NotEnoughBalance: true},
		},
		{
			name:     "withdraw zero amount from account with balance",
			recordId: "record-008",
			bankAccount: BankAccount{
				ID:            "account-008",
				HasBankUserID: HasBankUserID{BankUserID: "user-008"},
				Type:          BankAccountTypeSavings,
				Amount:        Money{Amount: 100, Currency: USD},
			},
			amount: Money{Amount: 0, Currency: USD},
			time:   baseTime,
			expectedResult: WithdrawResult{
				BankAccount: BankAccount{
					ID:            "account-008",
					HasBankUserID: HasBankUserID{BankUserID: "user-008"},
					Type:          BankAccountTypeSavings,
					Amount:        Money{Amount: 100, Currency: USD},
					LastTransactionRecordID: "record-008",
					LastTransactionTime:     baseTime,
				},
				UpdateBankAccountRequest: UpdateBankAccountRequest{
					ID:                      "account-008",
					Amount:                  ptrMoney(Money{Amount: 100, Currency: USD}),
					LastTransactionRecordID: ptrTxRecordID("record-008"),
					LastTransactionTime:     ptrTime(baseTime),
				},
				TransactionRecord: TransactionRecord{
					ID:               "record-008",
					HasBankUserID:    HasBankUserID{BankUserID: "user-008"},
					HasBankAccountID: HasBankAccountID{BankAccountID: "account-008"},
					Type:             TransactionTypeWithdrawal,
					WithdrawAmount:   Money{Amount: 0, Currency: USD},
					AfterAmount:      Money{Amount: 100, Currency: USD},
					DBItem: DBItem{
						HasCreatedAt: HasCreatedAt{CreatedAt: baseTime},
						HasUpdatedAt: HasUpdatedAt{UpdatedAt: baseTime},
					},
				},
				NotEnoughBalance: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Withdraw(
				tc.recordId,
				tc.bankAccount,
				tc.amount,
				tc.time,
			)

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf(
					"Withdraw() mismatch\ngot:  %+v\nwant: %+v",
					result,
					tc.expectedResult,
				)
			}
		})
	}
}
