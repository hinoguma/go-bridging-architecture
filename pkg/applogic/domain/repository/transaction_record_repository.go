package repository

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type TransactionRecordRepository interface {
	Get(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) (model.TransactionRecord, error)
	Create(ctx context.Context, account model.TransactionRecord, optionaltFuncs ...model.DBOperationOptionalFunc) error
	Put(ctx context.Context, account model.TransactionRecord, optionaltFuncs ...model.DBOperationOptionalFunc) error
	Delete(ctx context.Context, id model.TransactionRecordID, optionaltFuncs ...model.DBOperationOptionalFunc) error
}
