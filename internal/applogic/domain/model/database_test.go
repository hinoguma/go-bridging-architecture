package model

import (
	"app/internal/crosscutting"
	"reflect"
	"testing"
)

func TestWithDBTransactionID(t *testing.T) {
	id1 := DBTransactionID("12345")
	emptyId := DBTransactionID("")
	testCases := []struct {
		name          string
		transactionID DBTransactionID
		options       *DBOperationOptions
		expect        *DBOperationOptions
	}{
		{
			name:          "With transaction ID",
			transactionID: id1,
			options:       crosscutting.Ptr(DBOperationOptions{}),
			expect: crosscutting.Ptr(
				DBOperationOptions{TransactionID: &id1},
			),
		},
		{
			name:          "With transaction ID to nil options",
			transactionID: "12345",
			options:       nil,
			expect:        nil,
		},
		{
			name:          "With empty transaction ID",
			transactionID: emptyId,
			options:       crosscutting.Ptr(DBOperationOptions{}),
			expect: crosscutting.Ptr(
				DBOperationOptions{TransactionID: &emptyId},
			),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := testCase.options
			WithDBTransactionID(testCase.transactionID)(options)
			if testCase.expect == nil && options != nil {
				t.Errorf("Expected options to be nil, but got %v", options)
			}
			if !reflect.DeepEqual(options, testCase.expect) {
				t.Errorf("Expected options to be %v, but got %v", testCase.expect, options)
			}
		})
	}
}

func TestWithSelectLock(t *testing.T) {
	testCases := []struct {
		name    string
		options *DBOperationOptions
		expect  *DBOperationOptions
	}{
		{
			name:    "With select lock",
			options: crosscutting.Ptr(DBOperationOptions{}),
			expect:  crosscutting.Ptr(DBOperationOptions{selectLock: crosscutting.Ptr(true)}),
		},
		{
			name:    "With select lock to nil options",
			options: nil,
			expect:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := testCase.options
			WithSelectLock()(options)
			if testCase.expect == nil && options != nil {
				t.Errorf("Expected options to be nil, but got %v", options)
			}
			if !reflect.DeepEqual(options, testCase.expect) {
				t.Errorf("Expected options to be %v, but got %v", testCase.expect, options)
			}
		})
	}
}

func TestApplyDBOperationOptionalFuncs(t *testing.T) {
	id1 := DBTransactionID("12345")
	id2 := DBTransactionID("67890")
	testCases := []struct {
		name   string
		optFns []DBOperationOptionalFunc
		expect DBOperationOptions
	}{
		{
			name: "Apply WithDBTransactionID and WithSelectLock",
			optFns: []DBOperationOptionalFunc{
				WithDBTransactionID(id1),
				WithSelectLock(),
			},
			expect: DBOperationOptions{
				TransactionID: &id1,
				selectLock:    crosscutting.Ptr(true),
			},
		},
		{
			name:   "Apply no options",
			optFns: []DBOperationOptionalFunc{},
			expect: DBOperationOptions{},
		},
		{
			name: "Apply Dupulicate options",
			optFns: []DBOperationOptionalFunc{
				WithDBTransactionID(id1),
				WithSelectLock(),
				WithDBTransactionID(id2),
			},
			expect: DBOperationOptions{
				TransactionID: &id2,
				selectLock:    crosscutting.Ptr(true),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := ApplyDBOperationOptionalFuncs(testCase.optFns...)
			if !reflect.DeepEqual(options, testCase.expect) {
				t.Errorf("Expected options to be %v, but got %v", testCase.expect, options)
			}
		})
	}
}
