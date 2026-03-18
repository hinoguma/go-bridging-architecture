package repository

import (
	"app/internal/applogic/domain/model"
	"context"
)

type DBTransactionManager interface {
	Begin(ctx context.Context, req model.DBTransactionBeginRequest) (model.DBTransactionID, error)
	Commit(ctx context.Context, id model.DBTransactionID) error
	Rollback(ctx context.Context, id model.DBTransactionID) error
}
