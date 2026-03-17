package rdb

import (
	"app/pkg/crosscutting/errors"
	"context"
	"database/sql"
)

type postgreSQLClient struct {
	db *sql.DB
}

func NewPostgreSQLClient(db *sql.DB) SQLClient {
	return &postgreSQLClient{db: db}
}

func (client postgreSQLClient) GetDB() *sql.DB {
	return client.db
}

func (client postgreSQLClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := SQLQueryContext(ctx, client.db.QueryContext, query, args...)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return rows, nil
}

func (client postgreSQLClient) TxQueryContext(ctx context.Context, conn TransactionConnection, query string, args ...any) (*sql.Rows, error) {
	rows, err := SQLQueryContext(ctx, conn.tx.QueryContext, query, args...)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return rows, nil
}

func (client postgreSQLClient) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	rows, err := SQLQueryRowContext(ctx, client.db.QueryRowContext, query, args...)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return rows, nil
}

func (client postgreSQLClient) TxQueryRowContext(ctx context.Context, conn TransactionConnection, query string, args ...any) (*sql.Row, error) {
	rows, err := SQLQueryRowContext(ctx, conn.tx.QueryRowContext, query, args...)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return rows, nil
}
