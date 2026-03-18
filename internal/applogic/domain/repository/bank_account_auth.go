package repository

import (
	"app/internal/applogic/domain/model"
	"context"
)

type BankAccountAuthenticator interface {
	VerifyToken(ctx context.Context, req model.BankAccountVerifyTokenRequest) model.BankAccountVerifyTokenResult
}
