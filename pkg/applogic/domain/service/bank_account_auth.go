package service

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/crosscutting/errors"
	"context"
)

type BankAccountAuthService interface {
	Authenticate(ctx context.Context, req model.BankAccountAuthenticateRequest) model.BankAccountAuthenticateResult
}

func NewBankAccountAuthService(
	bankAccountRepository repository.BankAccountRepository,
	bankAccountAuthenticator repository.BankAccountAuthenticator,
) BankAccountAuthService {
	return bankAccountAuthenticateService{
		BankAccountRepository:    bankAccountRepository,
		BankAccountAuthenticator: bankAccountAuthenticator,
	}
}

type bankAccountAuthenticateService struct {
	BankAccountRepository    repository.BankAccountRepository
	BankAccountAuthenticator repository.BankAccountAuthenticator
}

func (service bankAccountAuthenticateService) Authenticate(
	ctx context.Context,
	req model.BankAccountAuthenticateRequest,
) model.BankAccountAuthenticateResult {
	// verify access token
	result := model.BankAccountAuthenticateResult{}
	verifyReq := model.BankAccountVerifyTokenRequest{
		Token: req.AccessToken,
	}
	verifyRes := service.BankAccountAuthenticator.VerifyToken(
		ctx, verifyReq,
	)
	if verifyRes.IsErr() {
		result.Err = errors.Lift(verifyRes.Err)
		return result
	}
	if !verifyRes.IsValid() {
		result.IsInValidToken = verifyRes.IsInValidToken
		result.IsTokenExpired = verifyRes.IsTokenExpired
		return result
	}

	// get user
	user, err := service.BankAccountRepository.GetByEmail(ctx, verifyRes.Email)
	if err != nil {
		result.Err = errors.Lift(err)
		return result
	}
	// todo
	result.BankAccount = user
	return result
}
