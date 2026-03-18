package model

import (
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"fmt"
)

type BankAccountPassword string

type BankAccountPasswordValidateResult struct {
	HasLower       bool
	HasUpper       bool
	HasNumber      bool
	HasSpecialChar bool
	Has12MoreChars bool
}

func (model BankAccountPasswordValidateResult) IsValid() bool {
	return model.HasLower &&
		model.HasUpper &&
		model.HasNumber &&
		model.HasSpecialChar &&
		model.Has12MoreChars
}

func ValidateBankAccountPassword(pw string) BankAccountPasswordValidateResult {
	return BankAccountPasswordValidateResult{}
}

type SignUpMethod string

const (
	SignUpMethodEmail  SignUpMethod = "EMAIL"
	SignUpMethodGoogle SignUpMethod = "GOOGLE"
)

// auth service id
// - cognito: username
type BankAccountAuthID string

func (value BankAccountAuthID) String() string {
	return string(value)
}

func NewBankAccountAuthIDForSignUpEmail(email crosscutting.Email) BankAccountAuthID {
	return BankAccountAuthID(
		fmt.Sprintf(`%s_%s`, SignUpMethodEmail, email),
	)
}

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
