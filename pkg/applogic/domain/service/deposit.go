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
			return result, errors.LiftWithCtx(err, ctx)
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
			return result, errors.LiftWithCtx(err, ctx)
		}
		if existingRecord.BankAccountID != req.BankAccountID {
			return result, errors.NewWithCtx("transaction record does not match with bank account", ctx)
		}
		if existingRecord.Type != model.TransactionTypeDeposit {
			return result, errors.NewWithCtx("transaction record type is not identical", ctx)
		}
		result.Record = existingRecord
		return result, nil
	}

	bancAccount, err := serv.bankAccountRepository.Get(
		ctx, req.BankAccountID, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.LiftWithCtx(err, ctx)
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
		return result, errors.LiftWithCtx(err, ctx)
	}

	err = serv.bankAccountRepository.Update(
		ctx, updateReq, model.WithDBTransactionID(txId),
	)
	if err != nil {
		return result, errors.LiftWithCtx(err, ctx)
	}

	result.BankAccount = updatedBankAccount
	result.Record = transactionRecord
	return result, nil
}
