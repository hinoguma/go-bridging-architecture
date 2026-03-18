package model

import (
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"time"
)

type HasCreatedAt struct {
	CreatedAt time.Time
}

type HasUpdatedAt struct {
	UpdatedAt time.Time
}

type DBItem struct {
	HasCreatedAt
	HasUpdatedAt
}

// to implement DBRecorder interface
func (model DBItem) DBRecordable() {}

func (model *DBItem) Create(t time.Time) {
	model.CreatedAt = t
	model.UpdatedAt = t
}

func (model *DBItem) Update(t time.Time) {
	model.UpdatedAt = t
}

type DBRecorder interface {
	DBRecordable()
}

type HasName struct {
	Name string
}

type DBTransactionID string

func (value DBTransactionID) String() string {
	return string(value)
}

func IssueDBTransactionID() DBTransactionID {
	return DBTransactionID(crosscutting.IssueRandomStrID())
}

const ErrorTypeDBTransactionNotFound errors.ErrorType = "DBTransactionNotFound"

func NewTransactionNotFoundError(id DBTransactionID) error {
	err := errors.New("transaction not found")
	err = errors.AddType(err, ErrorTypeDBTransactionNotFound)
	return errors.AddTagString(err, "transaction_id", id.String())
}

type DBTransactionBeginRequest struct {
	ID *DBTransactionID
}

func (model DBTransactionBeginRequest) HasDBTransactionID() bool {
	return model.ID != nil
}

func (model DBTransactionBeginRequest) GetTransactionID() DBTransactionID {
	if model.ID == nil {
		return ""
	}
	return *model.ID
}

func (model DBTransactionBeginRequest) IssueTransactionIDIfNotExists() DBTransactionID {
	if model.ID == nil {
		return IssueDBTransactionID()
	}
	return *model.ID
}

type UseDBTransaction struct {
	TxID *DBTransactionID
}

func (model UseDBTransaction) HasDBTransactionID() bool {
	return model.TxID != nil
}

func (model UseDBTransaction) GetTransactionID() DBTransactionID {
	if model.TxID == nil {
		return ""
	}
	return *model.TxID
}

func (model UseDBTransaction) SetTransactionID(txId DBTransactionID) UseDBTransaction {
	model.TxID = &txId
	return model
}

type HasDBTransactionID struct {
	DBTransactionID DBTransactionID
}

type DBOperationOptions struct {
	TransactionID *DBTransactionID
}

func (options DBOperationOptions) HasDBTransactionID() bool {
	return options.TransactionID != nil
}

func (options DBOperationOptions) GetTransactionID() DBTransactionID {
	if options.TransactionID == nil {
		return ""
	}
	return *options.TransactionID
}

type DBOperationOptionalFunc func(options *DBOperationOptions)

func WithDBTransactionID(txId DBTransactionID) DBOperationOptionalFunc {
	return func(options *DBOperationOptions) {
		options.TransactionID = &txId
	}
}

func ApplyDBOperationOptionalFuncs(optFns ...DBOperationOptionalFunc) DBOperationOptions {
	options := DBOperationOptions{}
	for _, optFn := range optFns {
		optFn(&options)
	}
	return options
}
