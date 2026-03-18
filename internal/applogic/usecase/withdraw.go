package usecase

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/applogic/domain/service"
	"app/internal/crosscutting/errors"
	"context"
)

type WithdrawUseCaseInput struct {
	BankAccountID model.BankAccountID
	TransactionID *model.TransactionRecordID
	Amount        model.Money
}

func (input WithdrawUseCaseInput) WithdrawServiceRequest() model.WithdrawServiceRequest {
	return model.WithdrawServiceRequest{
		BankAccountID: input.BankAccountID,
		Amount:        input.Amount,
	}
}

type WithdrawUseCaseOutput struct {
	TransactionRecord model.TransactionRecord
	NotEnoughBalance  bool
}

type WithdrawUseCase interface {
	Do(ctx context.Context, input WithdrawUseCaseInput) (WithdrawUseCaseOutput, error)
}

func NewWithdrawUseCase(
	withdrawService service.WithdrawService,
	dbTransactionManager repository.DBTransactionManager,
) WithdrawUseCase {
	return &withdrawUseCase{
		withdrawService:      withdrawService,
		dbTransactionManager: dbTransactionManager,
	}
}

type withdrawUseCase struct {
	withdrawService      service.WithdrawService
	dbTransactionManager repository.DBTransactionManager
}

func (useCase withdrawUseCase) Do(ctx context.Context, input WithdrawUseCaseInput) (WithdrawUseCaseOutput, error) {
	withdrawReq := input.WithdrawServiceRequest()

	txId, err := useCase.dbTransactionManager.Begin(ctx, model.DBTransactionBeginRequest{})
	if err != nil {
		return WithdrawUseCaseOutput{}, errors.Lift(err)
	}
	withdrawReq.SetTransactionID(txId)

	withdrawRes, err := func() (model.WithdrawServiceResult, error) {
		withdrawRes, err := useCase.withdrawService.Do(ctx, withdrawReq)
		if err != nil {
			return model.WithdrawServiceResult{}, errors.Lift(err)
		}

		err = useCase.dbTransactionManager.Commit(ctx, withdrawRes.TxID)
		if err != nil {
			return model.WithdrawServiceResult{}, errors.Lift(err)
		}
		return withdrawRes, nil
	}()

	if err != nil {
		rollbackErr := useCase.dbTransactionManager.Rollback(ctx, txId)
		if rollbackErr != nil {
			err = errors.AddSubErr(err, rollbackErr)
		}
		return WithdrawUseCaseOutput{}, err
	}

	if withdrawRes.NotEnoughBalance {
		return WithdrawUseCaseOutput{
			NotEnoughBalance: true,
		}, nil
	}

	return WithdrawUseCaseOutput{
		TransactionRecord: withdrawRes.Record,
	}, nil

}
