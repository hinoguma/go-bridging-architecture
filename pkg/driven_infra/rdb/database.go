package rdb

import (
	"app/pkg/crosscutting/errors"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type SQLDatabaseItemID interface {
	string | int64
}

type SQLDatabaseItem interface {
	ToMap() map[string]interface{}
	SetBySQLRow(row *sql.Row) error
}

type SQLOperationOptions struct {
	transactionID *string
}

func (options SQLOperationOptions) HasTransactionID() bool {
	return options.transactionID != nil
}

func (options SQLOperationOptions) GetTransactionID() string {
	if options.transactionID == nil {
		return ""
	}
	return *options.transactionID
}

func (options *SQLOperationOptions) SetTransactionID(id string) {
	options.transactionID = &id
}

type SQLOperationOptionalFunc func(options *SQLOperationOptions)

func WithTransactionID(id string) SQLOperationOptionalFunc {
	return func(options *SQLOperationOptions) {
		options.transactionID = &id
	}
}

func ApplySQLOperationOptionalFunc(optionalFuncs ...SQLOperationOptionalFunc) SQLOperationOptions {
	options := SQLOperationOptions{}
	for _, optionalFunc := range optionalFuncs {
		optionalFunc(&options)
	}
	return options
}

type SQLClient interface {
	GetDB() *sql.DB

	QueryContext(
		ctx context.Context,
		query string,
		args []any,
		options SQLOperationOptions,
	) (*sql.Rows, error)

	QueryRowContext(
		ctx context.Context,
		query string,
		args []any,
		options SQLOperationOptions,
	) (*sql.Row, error)

	GetRowByID(
		ctx context.Context,
		tableName string,
		id string,
		fields []string,
		options SQLOperationOptions,
	) (*sql.Row, error)

	CreateRow(
		ctx context.Context,
		tableName string,
		item SQLDatabaseItem,
		options SQLOperationOptions,
	) (*sql.Row, error)

	UpdateRowByStrID(
		ctx context.Context,
		tableName string,
		id string,
		updateFields UpdateFieldRequests,
		options SQLOperationOptions,
	) (*sql.Row, error)
}

type ExecSQLQueryFunc func(ctx context.Context, query string, args ...any) (*sql.Rows, error)
type ExecSQLQueryRowFunc func(ctx context.Context, query string, args ...any) *sql.Row

//
//func SQLQueryContextWithOpt(ctx context.Context, queryFunc ExecSQLQueryFunc, query string, args []any, optionalFuncs... SQLOperationOptionalFunc) (*sql.Rows, error) {
//	// log, metrics, tracing, etc.
//
//	options := ApplySQLOperationOptionalFunc(optionalFuncs...)
//
//	row, err := queryFunc(ctx, query, args...)
//	// log, metrics, tracing, etc.
//
//	if err != nil {
//		if errors.As(err, &sql.ErrNoRows) {
//			err = errors.NewDataNotFoundErr()
//		}
//		return row, errors.LiftWithCtx(err, ctx)
//	}
//
//	// success
//	return row, nil
//}

func SQLQueryContext(ctx context.Context, queryFunc ExecSQLQueryFunc, query string, args ...any) (*sql.Rows, error) {
	// log, metrics, tracing, etc.

	row, err := queryFunc(ctx, query, args...)
	// log, metrics, tracing, etc.

	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			err = errors.NewDataNotFoundErr()
		}
		return row, errors.LiftWithCtx(err, ctx)
	}

	// success
	return row, nil
}

func SQLQueryRowContext(ctx context.Context, queryFunc ExecSQLQueryRowFunc, query string, args ...any) (*sql.Row, error) {
	// log, metrics, tracing, etc.

	row := queryFunc(ctx, query, args...)
	// log, metrics, tracing, etc.

	// success
	return row, nil
}

func BuildSelectQueryWithId[T string | int64](tableName string, id T, fields []string) (string, []any) {
	return BuildSelectQueryWithSingleCondition(tableName, fields, "id", id)
}

func BuildSelectQueryWithSingleCondition[T string | int64](tableName string, fields []string, whereField string, val T) (string, []any) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE $1 = $2`, strings.Join(fields, ", "), tableName)
	values := []any{whereField, val}
	return query, values
}

func BuildUpdateQueryWithId[T string | int64](id T, tableName string, updateFields UpdateFieldRequests) (string, []any) {
	setClause, values := updateFields.ToSQLSetClause()
	query := "UPDATE " + tableName + " SET " + setClause + " WHERE id = $" + string(len(values)+1)
	values = append(values, id)
	return query, values
}

func BuildDeleteQueryWithId[T string | int64](id T, tableName string) (string, []any) {
	query := "DELETE FROM " + tableName + " WHERE id = $1"
	values := []any{id}
	return query, values
}

func BuildInsertQuery(tableName string, fields map[string]any) (string, []any) {
	columns := ""
	placeholders := ""
	values := make([]any, 0)
	i := 1
	for col, val := range fields {
		if i > 1 {
			columns += ", "
			placeholders += ", "
		}
		columns += col
		placeholders += "$" + string(i)
		values = append(values, val)
		i++
	}
	query := "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + placeholders + ")"
	return query, values
}

// todo: check if postgre has upsert statement
func BuildUpsertQuery(tableName string, fields map[string]any) (string, []any) {
	columns := ""
	placeholders := ""
	values := make([]any, 0)
	i := 1
	for col, val := range fields {
		if i > 1 {
			columns += ", "
			placeholders += ", "
		}
		columns += col
		placeholders += "$" + string(i)
		values = append(values, val)
		i++
	}
	query := "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + placeholders + ")"
	return query, values
}

type UpdateFieldRequest struct {
	FieldName string
	NewValue  any
}

func NewUpdateFieldRequest(fieldName string, newValue any) UpdateFieldRequest {
	return UpdateFieldRequest{
		FieldName: fieldName,
		NewValue:  newValue,
	}
}

type UpdateFieldRequests []UpdateFieldRequest

func (requests *UpdateFieldRequests) Append(field string, value any) *UpdateFieldRequests {
	*requests = append(*requests, UpdateFieldRequest{
		FieldName: field,
		NewValue:  value,
	})
	return requests
}

func (requests UpdateFieldRequests) ToSQLSetClause() (string, []any) {
	setClause := ""
	values := make([]any, 0)
	for i, req := range requests {
		if i > 0 {
			setClause += ", "
		}
		setClause += req.FieldName + " = $" + string(i+1)
		values = append(values, req.NewValue)
	}

	// example: "field1 = $1, field2 = $2", [value1, value2]
	return setClause, values
}
