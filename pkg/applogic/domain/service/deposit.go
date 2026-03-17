package service

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/crosscutting/errors"
	"context"
)

type DepositService interface {
	Do(ctx context.Context, req model.DepositRequest) (model.DepositResult, error)
}

type depositService struct {
	dbTransactionManager        repository.DBTransactionManager
	bankAccountRepository       repository.BankAccountRepository
	transactionRecordRepository repository.TransactionRecordRepository
}

func NewDepositService(
	dbTransactionManager repository.DBTransactionManager,
	bankAccountRepository repository.BankAccountRepository,
	transactionRecordRepository repository.TransactionRecordRepository,
) DepositService {
	return &depositService{
		dbTransactionManager:        dbTransactionManager,
		bankAccountRepository:       bankAccountRepository,
		transactionRecordRepository: transactionRecordRepository,
	}
}

func (serv depositService) Do(ctx context.Context, req model.DepositRequest) (model.DepositResult, error) {
	var result model.DepositResult
	txId := req.GetTransactionID()
	if !req.HasTransactionID() {
		newTxId, err := serv.dbTransactionManager.Begin(ctx, model.DBTransactionBeginRequest{})
		if err != nil {
			return result, errors.Lift(err)
		}
		txId = newTxId
	}
	result.TxID = txId

	bancAccount, err := serv.bankAccountRepository.Get(
		ctx, req.BankAccountID, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.Lift(err)
	}

	bankAccount, transactionRecord := model.Deposit(
		bancAccount, req.Amount, req.GetRequestAt(),
	)

	err = serv.transactionRecordRepository.Create(
		ctx, transactionRecord, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.Lift(err)
	}

	err = serv.bankAccountRepository.Put(ctx, bankAccount, model.WithDBTransactionID(txId))
	if err != nil {
		err = errors.Lift(err)
		subErr := serv.transactionRecordRepository.Delete(ctx, transactionRecord.ID, model.WithDBTransactionID(txId))
		return result, errors.AddSubErr(err, subErr)
	}

	result.Record = transactionRecord
	return result, nil
}
