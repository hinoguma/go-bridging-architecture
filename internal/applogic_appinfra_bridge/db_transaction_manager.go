package applogic_appinfra_bridge

import (
	"app/internal/appinfra/postgres"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"context"
)

func NewDBTransactionManager(transactionManager postgres.TransactionManagerIF) repository.DBTransactionManager {
	return dbTransactionManager{
		transactionManager: transactionManager,
	}
}

type dbTransactionManager struct {
	transactionManager postgres.TransactionManagerIF
}

func (manager dbTransactionManager) Begin(ctx context.Context, req model.DBTransactionBeginRequest) (model.DBTransactionID, error) {
	txId := req.IssueTransactionIDIfNotExists()
	_, err := manager.transactionManager.Begin(ctx, txId.String(), nil)
	if err != nil {
		return "", errors.Lift(err)
	}
	return txId, nil
}

func (manager dbTransactionManager) Commit(ctx context.Context, id model.DBTransactionID) error {
	err := manager.transactionManager.Commit(ctx, id.String())
	if err != nil {
		return errors.Lift(err)
	}
	return nil
}

func (manager dbTransactionManager) Rollback(ctx context.Context, id model.DBTransactionID) error {
	err := manager.transactionManager.Rollback(ctx, id.String())
	if err != nil {
		return errors.Lift(err)
	}
	return nil
}
