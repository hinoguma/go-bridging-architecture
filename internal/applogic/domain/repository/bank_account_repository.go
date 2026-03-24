package repository

import (
	"app/internal/applogic/domain/model"
	"context"
)

type BankAccountRepository interface {
	Get(ctx context.Context, id model.BankAccountID, optionalFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error)
	Create(ctx context.Context, account model.BankAccount, optionalFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error)
	Update(ctx context.Context, request model.UpdateBankAccountRequest, optionalFuncs ...model.DBOperationOptionalFunc) error
}
