package drived_infra_port

import (
	"app/pkg/applogic/domain/model"
	"context"
)

type TransactionRecordRepository interface {
	Get(ctx context.Context, id model.TransactionRecordID) (model.TransactionRecord, error)
	Create(ctx context.Context, account model.TransactionRecord) error
	Put(ctx context.Context, account model.TransactionRecord) error
}
