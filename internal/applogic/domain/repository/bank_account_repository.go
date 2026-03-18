package repository

import (
	"app/internal/applogic/domain/model"
	"context"
)

type BankAccountRepository interface {
	Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error)
	Create(ctx context.Context, account model.BankAccount, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error)
	Update(ctx context.Context, request model.UpdateBankAccountRequest, optionaltFuncs ...model.DBOperationOptionalFunc) error
}
