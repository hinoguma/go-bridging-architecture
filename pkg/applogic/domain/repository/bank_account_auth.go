package repository

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type BankAccountAuthenticator interface {
	VerifyToken(ctx context.Context, req model.BankAccountVerifyTokenRequest) model.BankAccountVerifyTokenResult
}
