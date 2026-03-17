package rdb

import (
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"context"
	"database/sql"
	"sync"
)

type IsolationLevel string

const (
	IsolationLevelDefault         IsolationLevel = "Default"
	IsolationLevelReadUncommitted IsolationLevel = "ReadUncommitted"
	IsolationLevelReadCommitted   IsolationLevel = "ReadCommitted"
	IsolationLevelWriteCommitted  IsolationLevel = "WriteCommitted"
	IsolationLevelRepeatableRead  IsolationLevel = "RepeatableRead"
	IsolationLevelSnapshot        IsolationLevel = "Snapshot"
	IsolationLevelSerializable    IsolationLevel = "Serializable"
	IsolationLevelLinearizable    IsolationLevel = "Linearizable"
)

func (value IsolationLevel) SQLLevel() sql.IsolationLevel {
	switch value {
	case IsolationLevelDefault:
		return sql.LevelDefault
	case IsolationLevelReadUncommitted:
		return sql.LevelReadUncommitted
	case IsolationLevelReadCommitted:
		return sql.LevelReadCommitted
	case IsolationLevelWriteCommitted:
		return sql.LevelWriteCommitted
	case IsolationLevelRepeatableRead:
		return sql.LevelRepeatableRead
	case IsolationLevelSnapshot:
		return sql.LevelSnapshot
	case IsolationLevelSerializable:
		return sql.LevelSerializable
	case IsolationLevelLinearizable:
		return sql.LevelLinearizable
	}
	return sql.LevelDefault
}

type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

func GenerateTxID() string {
	return crosscutting.IssueRandomStrID()
}

/********************************************
	Connection Pool
 ********************************************/

var txConMngSingle TxConnectionPoolIF = NewTxConnectionManager()

type TransactionConnection struct {
	id string
	tx *sql.Tx
}

func NewTransactionConnection(id string, tx *sql.Tx) TransactionConnection {
	return TransactionConnection{
		id: id,
		tx: tx,
	}
}

func (conn TransactionConnection) HasTx() bool {
	return conn.tx != nil
}
func (conn TransactionConnection) GetID() string {
	return conn.id
}

type TxConnectionPoolIF interface {
	Get(id string) (TransactionConnection, bool)
	Set(con TransactionConnection)
}

type txConnectionManager struct {
	connections map[string]TransactionConnection
	mu          sync.Mutex
}

func NewTxConnectionManager() TxConnectionPoolIF {
	return &txConnectionManager{
		connections: make(map[string]TransactionConnection),
		mu:          sync.Mutex{},
	}
}

func (mng *txConnectionManager) Get(id string) (TransactionConnection, bool) {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if mng.connections == nil {
		return TransactionConnection{}, false
	}
	con, ok := mng.connections[id]
	return con, ok
}

func (mng *txConnectionManager) Set(con TransactionConnection) {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if mng.connections == nil {
		mng.connections = make(map[string]TransactionConnection)
	}
	if con.id == "" {
		return
	}
	mng.connections[con.id] = con
}

type globalTxConnectionPool struct {
}

func (g globalTxConnectionPool) Get(id string) (TransactionConnection, bool) {
	return txConMngSingle.Get(id)
}

func (g globalTxConnectionPool) Set(con TransactionConnection) {
	txConMngSingle.Set(con)
}

func NewGlobalTxConnectionPool() TxConnectionPoolIF {
	return globalTxConnectionPool{}
}

func GlobalTxConnectionPool() TxConnectionPoolIF {
	return txConMngSingle
}

/********************************************
	Transaction Manager
 ********************************************/

type TransactionManagerIF interface {
	Begin(ctx context.Context, id string, opts *TxOptions) (TransactionConnection, error)
	Rollback(ctx context.Context, id string) error
	Commit(ctx context.Context, id string) error
}

type transactionManager struct {
	client *sql.DB
}

func NewTransactionManager(client *sql.DB) TransactionManagerIF {
	return &transactionManager{client: client}
}

func (mng transactionManager) Begin(ctx context.Context, id string, opts *TxOptions) (TransactionConnection, error) {
	var sqlOpts *sql.TxOptions
	if opts != nil {
		sqlOpts = &sql.TxOptions{
			Isolation: opts.Isolation.SQLLevel(),
			ReadOnly:  opts.ReadOnly,
		}
	}
	tx, err := mng.client.BeginTx(ctx, sqlOpts)
	if err != nil {
		return TransactionConnection{}, errors.LiftWithCtx(err, ctx)
	}
	con := TransactionConnection{
		id: id,
		tx: tx,
	}
	txConMngSingle.Set(con)
	return con, nil
}

func (mng transactionManager) Rollback(ctx context.Context, id string) error {
	conn, ok := txConMngSingle.Get(id)
	if !ok || !conn.HasTx() {
		err := errors.NewWithCtx("Not Found tx connection.", ctx)
		err = errors.AddTagString(err, "txId", id)
		return err
	}
	return conn.tx.Rollback()
}

func (mng transactionManager) Commit(ctx context.Context, id string) error {
	conn, ok := txConMngSingle.Get(id)
	if !ok || !conn.HasTx() {
		err := errors.NewWithCtx("Not Found tx connection.", ctx)
		err = errors.AddTagString(err, "txId", id)
		return err
	}
	return conn.tx.Commit()
}
