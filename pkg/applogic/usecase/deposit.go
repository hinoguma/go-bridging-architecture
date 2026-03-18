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
	TransactionID *model.TransactionRecordID
	Amount        model.Money
}

func (input DepositUseCaseInput) DepositServiceRequest() model.DepositServiceRequest {
	return model.DepositServiceRequest{
		BankAccountID: input.BankAccountID,
		Amount:        input.Amount,
	}
}

type DepositUseCaseOutput struct {
	TransactionRecord model.TransactionRecord
	NotEnoughBalance  bool
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
	depositReq := input.DepositServiceRequest()

	txId, err := useCase.dbTransactionManager.Begin(ctx, model.DBTransactionBeginRequest{})
	if err != nil {
		return DepositUseCaseOutput{}, errors.Lift(err)
	}
	depositReq.SetTransactionID(txId)

	depositRes, err := func() (model.DepositServiceResult, error) {
		depositRes, err := useCase.depositService.Do(ctx, depositReq)
		if err != nil {
			return model.DepositServiceResult{}, errors.Lift(err)
		}

		err = useCase.dbTransactionManager.Commit(ctx, depositRes.TxID)
		if err != nil {
			return model.DepositServiceResult{}, errors.Lift(err)
		}
		return depositRes, nil
	}()

	if err != nil {
		rollbackErr := useCase.dbTransactionManager.Rollback(ctx, txId)
		if rollbackErr != nil {
			err = errors.AddSubErr(err, rollbackErr)
		}
		return DepositUseCaseOutput{}, err
	}

	if depositRes.NotEnoughBalance {
		return DepositUseCaseOutput{
			NotEnoughBalance: true,
		}, nil
	}

	return DepositUseCaseOutput{
		TransactionRecord: depositRes.Record,
	}, nil

}
