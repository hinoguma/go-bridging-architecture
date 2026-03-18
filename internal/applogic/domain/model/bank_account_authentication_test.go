package model

import "testing"

func TestBankAccountVerifyTokenResult_IsValid(t *testing.T) {
	testCases := []struct {
		name   string
		result BankAccountVerifyTokenResult
		expect bool
	}{
		{
			name: "Valid token",
			result: BankAccountVerifyTokenResult{
				IsInValidToken: false,
				IsTokenExpired: false,
			},
			expect: true,
		},
		{
			name: "Invalid token",
			result: BankAccountVerifyTokenResult{
				IsInValidToken: true,
				IsTokenExpired: false,
			},
			expect: false,
		},
		{
			name: "Expired token",
			result: BankAccountVerifyTokenResult{
				IsInValidToken: false,
				IsTokenExpired: true,
			},
			expect: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.result.IsValid()
			if result != testCase.expect {
				t.Errorf("Expected IsValid to be %v, but got %v", testCase.expect, result)
			}
		})
	}
}
