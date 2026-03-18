package model

import (
	"app/internal/crosscutting/errors"
)

type BankAccountVerifyTokenRequest struct {
	Token string
	HasRequestAt
}

type BankAccountVerifyTokenResult struct {
	errors.HasError
	IsTokenExpired bool
	IsInValidToken bool
	BankAccountID  BankAccountID
}

func (res BankAccountVerifyTokenResult) IsValid() bool {
	return !res.IsTokenExpired && !res.IsInValidToken
}

type BankAccountAuthenticateRequest struct {
	AccessToken string
}

type BankAccountAuthenticateResult struct {
	errors.HasError
	IsTokenExpired bool
	IsInValidToken bool
	BankAccount    BankAccount
}
