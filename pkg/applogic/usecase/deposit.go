package usecase

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"context"
)

type DepositUseCaseInput struct {
	BankAccountID model.BankAccountID
	Amount        model.Money
}

type DepositUseCaseOutput struct {
	TransactionID model.TransactionRecordID
}

type DepositUseCase interface {
	Do(ctx context.Context, input DepositUseCaseInput) (DepositUseCaseOutput, error)
}

func NewDepositUseCase(
	bankAccountRepository repository.BankAccountRepository,
	transactionRecordRepository repository.TransactionRecordRepository,
) DepositUseCase {
	return &depositUseCase{
		bankAccountRepository:       bankAccountRepository,
		transactionRecordRepository: transactionRecordRepository,
	}
}

type depositUseCase struct {
	bankAccountRepository       repository.BankAccountRepository
	transactionRecordRepository repository.TransactionRecordRepository
}

func (d depositUseCase) Do(ctx context.Context, input DepositUseCaseInput) (DepositUseCaseOutput, error) {
	//TODO implement me
	panic("implement me")
}
