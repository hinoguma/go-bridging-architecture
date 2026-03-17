package usecase

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/service"
	"app/pkg/crosscutting/errors"
	"context"
)

type BankAccountAuthMWRequest struct {
	Token string
}

func (req BankAccountAuthMWRequest) AuthServiceRequest() model.BankAccountAuthenticateRequest {
	return model.BankAccountAuthenticateRequest{AccessToken: req.Token}
}

type BankAccountAuthMWResult struct {
	IsTokenExpired bool
	IsInValidToken bool
	BankAccount    model.BankAccount
	errors.HasError
}

func (res BankAccountAuthMWResult) SetByAuthServiceResult(authRes model.BankAccountAuthenticateResult) {
	res.IsTokenExpired = authRes.IsTokenExpired
	res.IsInValidToken = authRes.IsInValidToken
	res.BankAccount = authRes.BankAccount
	res.Err = authRes.Err
}

type BankAccountAuthMWUseCase interface {
	Do(ctx context.Context, req BankAccountAuthMWRequest) BankAccountAuthMWResult
}

type bankAccountAuthMWUseCase struct {
	BankAccountAuthService service.BankAccountAuthService
}

func NewBankAccountAuthMWUseCase(
	bankAccountAuthService service.BankAccountAuthService,
) BankAccountAuthMWUseCase {
	return bankAccountAuthMWUseCase{
		BankAccountAuthService: bankAccountAuthService,
	}
}

func (uc bankAccountAuthMWUseCase) Do(ctx context.Context, req BankAccountAuthMWRequest) BankAccountAuthMWResult {
	result := BankAccountAuthMWResult{}
	authRes := uc.BankAccountAuthService.Authenticate(
		ctx, req.AuthServiceRequest(),
	)
	if authRes.IsErr() {
		result.Err = errors.LiftWithCtx(authRes.Err, ctx)
		return result
	}
	result.SetByAuthServiceResult(authRes)
	return result
}
