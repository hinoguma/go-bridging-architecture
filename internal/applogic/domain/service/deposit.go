package service

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
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
	transactionRecordID := model.IssueTransactionRecordID()
	if req.HasTransactionRecordID() {
		transactionRecordID = req.GetTransactionRecordID()
		existingRecord, err := serv.transactionRecordRepository.Get(
			ctx, transactionRecordID,
			model.WithDBTransactionID(txId),
		)
		isNotFound := errors.IsDataNotFoundError(err)
		if err != nil && !isNotFound {
			return result, errors.LiftWithCtx(err, ctx)
		}
		if !isNotFound {
			if existingRecord.BankAccountID != req.BankAccountID {
				return result, errors.NewWithCtx("transaction record does not match with bank account", ctx)
			}
			if existingRecord.Type != model.TransactionTypeDeposit {
				return result, errors.NewWithCtx("transaction record type is not identical", ctx)
			}
			result.Record = existingRecord
			return result, nil
		}
	}

	bankAccount, err := serv.bankAccountRepository.Get(
		ctx, req.BankAccountID,
		model.WithDBTransactionID(txId),
		model.WithSelectLock(),
	)
	if err != nil {
		return result, errors.LiftWithCtx(err, ctx)
	}

	depositRes := model.Deposit(
		transactionRecordID, bankAccount, req.Amount, req.GetRequestAt(),
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
