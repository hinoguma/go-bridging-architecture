package service

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"context"
)

type WithdrawService interface {
	Do(ctx context.Context, req model.WithdrawServiceRequest) (model.WithdrawServiceResult, error)
}

type withdrawService struct {
	dbTransactionManager        repository.DBTransactionManager
	bankAccountRepository       repository.BankAccountRepository
	transactionRecordRepository repository.TransactionRecordRepository
}

func NewWithdrawService(
	dbTransactionManager repository.DBTransactionManager,
	bankAccountRepository repository.BankAccountRepository,
	transactionRecordRepository repository.TransactionRecordRepository,
) WithdrawService {
	return &withdrawService{
		dbTransactionManager:        dbTransactionManager,
		bankAccountRepository:       bankAccountRepository,
		transactionRecordRepository: transactionRecordRepository,
	}
}

func (serv withdrawService) Do(ctx context.Context, req model.WithdrawServiceRequest) (model.WithdrawServiceResult, error) {
	var result model.WithdrawServiceResult
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
			ctx, transactionRecordID, model.WithDBTransactionID(txId),
		)
		isNotFound := errors.IsDataNotFoundError(err)
		if err != nil && !isNotFound {
			return result, errors.LiftWithCtx(err, ctx)
		}
		if !isNotFound {
			if existingRecord.BankAccountID != req.BankAccountID {
				return result, errors.NewWithCtx("transaction record does not match with bank account", ctx)
			}
			if existingRecord.Type != model.TransactionTypeWithdrawal {
				return result, errors.NewWithCtx("transaction record type is not identical", ctx)
			}
			result.Record = existingRecord
			return result, nil
		}
	}

	bancAccount, err := serv.bankAccountRepository.Get(
		ctx, req.BankAccountID,
		model.WithDBTransactionID(txId),
		model.WithSelectLock(),
	)
	if err != nil {
		return result, errors.LiftWithCtx(err, ctx)
	}

	withdrawRes := model.Withdraw(
		transactionRecordID, bancAccount, req.Amount, req.GetRequestAt(),
	)
	if withdrawRes.NotEnoughBalance {
		result.NotEnoughBalance = true
		return result, nil
	}

	transactionRecord := withdrawRes.TransactionRecord
	updateReq := withdrawRes.UpdateBankAccountRequest
	updatedBankAccount := withdrawRes.BankAccount

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
