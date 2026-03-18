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

func (client postgreSQLClient) QueryContext(
	ctx context.Context, query string, args []any, options SQLOperationOptions,
) (*sql.Rows, error) {

	var queryFunc ExecSQLQueryFunc = client.db.QueryContext
	if options.HasDBTransactionID() {
		transactionID := options.GetTransactionID()
		conn, ok := GlobalTxConnectionPool().Get(transactionID)
		if !ok {
			return nil, NewTransactionNotFoundError(transactionID)
		}
		queryFunc = conn.tx.QueryContext
	}

	rows, err := SQLQueryContext(ctx, queryFunc, query, args...)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return rows, nil
}

func (client postgreSQLClient) QueryRowContext(
	ctx context.Context, query string, args []any, options SQLOperationOptions,
) (*sql.Row, error) {
	var queryFunc ExecSQLQueryRowFunc = client.db.QueryRowContext
	if options.HasDBTransactionID() {
		transactionID := options.GetTransactionID()
		conn, ok := GlobalTxConnectionPool().Get(transactionID)
		if !ok {
			err := NewTransactionNotFoundError(transactionID)
			return nil, errors.LiftWithCtx(err, ctx)
		}
		queryFunc = conn.tx.QueryRowContext
	}

	row := SQLQueryRowContext(ctx, queryFunc, query, args...)
	return row, nil
}

func (client postgreSQLClient) GetRowByID(
	ctx context.Context, tableName string, id SQLRecordID, fields []string, options SQLOperationOptions,
) (*sql.Row, error) {
	query, args := BuildSelectQueryWithId(tableName, id, fields)
	row, err := client.QueryRowContext(ctx, query, args, options)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return row, nil
}

func (client postgreSQLClient) CreateRow(
	ctx context.Context,
	tableName string,
	item SQLDatabaseItem,
	options SQLOperationOptions,
) (*sql.Row, error) {
	query, args := BuildInsertQuery(tableName, item.ToMap())
	row, err := client.QueryRowContext(ctx, query, args, options)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return row, nil
}

func (client postgreSQLClient) UpdateRowByStrID(
	ctx context.Context,
	tableName string,
	id SQLRecordID,
	updateFields UpdateFieldRequests,
	options SQLOperationOptions,
) (*sql.Row, error) {
	query, args := BuildUpdateQueryWithId(tableName, id, updateFields)
	row, err := client.QueryRowContext(ctx, query, args, options)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return row, nil
}
