package repository

import (
	"app/internal/applogic/domain/model"
	"context"
)

type TransactionRecordRepository interface {
	Get(ctx context.Context, id model.TransactionRecordID, optionalFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error)
	Create(ctx context.Context, account model.TransactionRecord, optionalFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error)
}
