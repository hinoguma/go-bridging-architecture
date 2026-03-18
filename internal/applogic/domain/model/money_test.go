package model

import "testing"

func TestCurrency_IsValid(t *testing.T) {
	testCases := []struct {
		name     string
		value    Currency
		expected bool
	}{
		{"valid USD", USD, true},
		{"valid EUR", EUR, true},
		{"valid GBP", GBP, true},
		{"valid JPY", JPY, true},
		{"invalid currency", Currency("ABC"), false},
		{"empty string", Currency(""), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.value.IsValid()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid USD", "USD", true},
		{"valid EUR", "EUR", true},
		{"valid GBP", "GBP", true},
		{"valid JPY", "JPY", true},
		{"invalid currency", "ABC", false},
		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateCurrency(tc.value)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	testCases := []struct {
		name   string
		money1 Money
		money2 Money
		expect bool
	}{
		{
			name:   "same currency, greater amount",
			money1: Money{Amount: 100, Currency: USD},
			money2: Money{Amount: 50, Currency: USD},
			expect: true,
		},
		{
			name:   "same currency, equal amount",
			money1: Money{Amount: 100, Currency: USD},
			money2: Money{Amount: 100, Currency: USD},
			expect: false,
		},
		{
			name:   "same currency, smaller amount",
			money1: Money{Amount: 50, Currency: USD},
			money2: Money{Amount: 100, Currency: USD},
			expect: false,
		},
		{
			name:   "different currency",
			money1: Money{Amount: 100, Currency: USD},
			money2: Money{Amount: 50, Currency: EUR},
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.money1.GreaterThan(tc.money2)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}
