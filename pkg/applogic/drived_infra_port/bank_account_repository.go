package drived_infra_port

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type BankAccountRepository interface {
	Get(ctx context.Context, id model.BankAccountID) (model.BankAccount, error)
	Create(ctx context.Context, account model.BankAccount) error
	Put(ctx context.Context, account model.BankAccount) error
}