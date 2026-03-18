package service

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/domain/repository"
	"app/pkg/crosscutting/errors"
	"context"
)

type DepositService interface {
	Do(ctx context.Context, req model.DepositServiceRequest) (model.DepositServiceResult, error)
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

func (serv depositService) Do(ctx context.Context, req model.DepositServiceRequest) (model.DepositServiceResult, error) {
	var result model.DepositServiceResult
	txId := req.GetTransactionID()
	if !req.HasDBTransactionID() {
		newTxId, err := serv.dbTransactionManager.Begin(
			ctx, model.DBTransactionBeginRequest{},
		)
		if err != nil {
			return result, errors.Lift(err)
		}
		txId = newTxId
	}
	result.TxID = txId

	// For idempotency
	if req.HasTransactionRecordID() {
		existingRecord, err := serv.transactionRecordRepository.Get(
			ctx, req.GetTransactionRecordID(), model.WithDBTransactionID(txId),
		)
		if err != nil {
			return result, errors.Lift(err)
		}
		if existingRecord.BankAccountID != req.BankAccountID {
			return result, errors.New("transaction record does not match with bank account")
		}
		if existingRecord.Type != model.TransactionTypeDeposit {
			return result, errors.New("transaction record type is not identical")
		}
		result.Record = existingRecord
		return result, nil
	}

	bancAccount, err := serv.bankAccountRepository.Get(
		ctx, req.BankAccountID, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.Lift(err)
	}

	depositRes := model.Deposit(
		bancAccount, req.Amount, req.GetRequestAt(),
	)

	transactionRecord := depositRes.TransactionRecord
	updateReq := depositRes.UpdateBankAccountRequest
	updatedBankAccount := depositRes.BankAccount

	transactionRecord, err = serv.transactionRecordRepository.Create(
		ctx, transactionRecord, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.Lift(err)
	}

	err = serv.bankAccountRepository.Update(
		ctx, updateReq, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.Lift(err)
	}

	result.BankAccount = updatedBankAccount
	result.Record = transactionRecord
	return result, nil
}
