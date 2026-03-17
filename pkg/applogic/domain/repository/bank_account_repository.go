package repository

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type BankAccountRepository interface {
	Get(ctx context.Context, id model.BankAccountID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.BankAccount, error)
	Create(ctx context.Context, account model.BankAccount, optionaltFuncs ...model.DBOperationOptionalFunc) error
	Put(ctx context.Context, account model.BankAccount, optionaltFuncs ...model.DBOperationOptionalFunc) error
}
