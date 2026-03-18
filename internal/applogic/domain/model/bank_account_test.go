package model

import (
	"app/internal/crosscutting"
	"reflect"
	"testing"
	"time"
)

func TestUpdateBankAccountRequest_UpdateItem(t *testing.T) {
	time1 := time.Unix(1773846000, 0) // 2026-03-19 00:00:00
	time2 := time.Unix(1773849600, 0) // 2026-03-19 01:00:00
	testCases := []struct {
		name   string
		model  UpdateBankAccountRequest
		item   BankAccount
		expect BankAccount
	}{
		{
			name: "Update amount",
			model: UpdateBankAccountRequest{
				ID: "account-123",
				Amount: crosscutting.Ptr(Money{
					Amount:   100,
					Currency: "USD",
				}),
			},
			item: BankAccount{
				ID:     "account-123",
				Amount: Money{Amount: 50, Currency: "USD"},
			},
			expect: BankAccount{
				ID:     "account-123",
				Amount: Money{Amount: 100, Currency: "USD"},
			},
		},
		{
			name: "Update last transaction record ID and time",
			model: UpdateBankAccountRequest{
				ID:                      "account-123",
				LastTransactionRecordID: crosscutting.Ptr(TransactionRecordID("record-456")),
				LastTransactionTime:     crosscutting.Ptr(time1),
			},
			item: BankAccount{
				ID:                      "account-123",
				LastTransactionRecordID: TransactionRecordID("record-123"),
				LastTransactionTime:     time2,
			},
			expect: BankAccount{
				ID:                      "account-123",
				LastTransactionRecordID: TransactionRecordID("record-456"),
				LastTransactionTime:     time1,
			},
		},

		{
			name: "id does not match, no update",
			model: UpdateBankAccountRequest{
				ID: "account-123",
				Amount: crosscutting.Ptr(Money{
					Amount:   100,
					Currency: "USD",
				}),
			},
			item: BankAccount{
				ID: "account-456",
			},
			expect: BankAccount{
				ID: "account-456",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.model.UpdateItem(testCase.item)
			if !reflect.DeepEqual(result, testCase.expect) {
				t.Errorf("Expected %v, but got %v", testCase.expect, result)
			}
		})
	}
}
