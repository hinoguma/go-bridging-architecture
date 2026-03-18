package model

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTransactionRecordDeposit(t *testing.T) {
	t1 := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		name         string
		id           TransactionRecordID
		accountAfter BankAccount
		amount       Money
		t            time.Time
		expect       TransactionRecord
	}{
		{
			name: "create deposit record",
			id:   "record-123",
			accountAfter: BankAccount{
				ID:     "account-456",
				Amount: Money{Amount: 1000, Currency: USD},
			},
			amount: Money{Amount: 1000, Currency: USD},
			t:      t1,
			expect: TransactionRecord{
				ID:               "record-123",
				Type:             TransactionTypeDeposit,
				HasBankAccountID: HasBankAccountID{BankAccountID: "account-456"},
				DepositAmount:    Money{Amount: 1000, Currency: USD},
				AfterAmount:      Money{Amount: 1000, Currency: USD},
				DBItem: DBItem{
					HasCreatedAt: HasCreatedAt{
						CreatedAt: t1,
					},
					HasUpdatedAt: HasUpdatedAt{
						UpdatedAt: t1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := NewTransactionRecordDeposit(tc.id, tc.accountAfter, tc.amount, tc.t)
			if !reflect.DeepEqual(record, tc.expect) {
				t.Errorf("unexpected record: got %v, want %v", record, tc.expect)
			}
		})
	}
}
