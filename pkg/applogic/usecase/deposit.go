package usecase

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/applogic/domain/service"
	"app/pkg/crosscutting/errors"
	"context"
)

type DepositUseCaseInput struct {
	BankAccountID model.BankAccountID
	Amount        model.Money
}

func (input DepositUseCaseInput) DepositRequest() model.DepositRequest {
	return model.DepositRequest{
		BankAccountID: input.BankAccountID,
		Amount:        input.Amount,
	}
}

type DepositUseCaseOutput struct {
	TransactionRecord model.TransactionRecord
}

type DepositUseCase interface {
	Do(ctx context.Context, input DepositUseCaseInput) (DepositUseCaseOutput, error)
}

func NewDepositUseCase(
	depositService service.DepositService,
	dbTransactionManager repository.DBTransactionManager,
) DepositUseCase {
	return &depositUseCase{
		depositService:       depositService,
		dbTransactionManager: dbTransactionManager,
	}
}

type depositUseCase struct {
	depositService       service.DepositService
	dbTransactionManager repository.DBTransactionManager
}

func (useCase depositUseCase) Do(ctx context.Context, input DepositUseCaseInput) (DepositUseCaseOutput, error) {
	depositReq := input.DepositRequest()

	txId, err := useCase.dbTransactionManager.Begin(ctx, model.DBTransactionBeginRequest{})
	if err != nil {
		return DepositUseCaseOutput{}, errors.Lift(err)
	}
	depositReq.SetTransactionID(txId)

	depositRes, err := func() (model.DepositResult, error) {
		depositRes, err := useCase.depositService.Do(ctx, depositReq)
		if err != nil {
			return model.DepositResult{}, errors.Lift(err)
		}

		err = useCase.dbTransactionManager.Commit(ctx, depositRes.TxID)
		if err != nil {
			return model.DepositResult{}, errors.Lift(err)
		}
		return depositRes, nil
	}()
	if err != nil {
		rollbackErr := useCase.dbTransactionManager.Rollback(ctx, txId)
		if rollbackErr != nil {
			return DepositUseCaseOutput{}, errors.Lift(rollbackErr)
		}
		return DepositUseCaseOutput{}, errors.Lift(err)
	}

	return DepositUseCaseOutput{
		TransactionRecord: depositRes.Record,
	}, nil

}
