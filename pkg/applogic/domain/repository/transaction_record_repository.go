package repository

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type TransactionRecordRepository interface {
	Get(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error)
	Create(ctx context.Context, account model.TransactionRecord, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error)
}
