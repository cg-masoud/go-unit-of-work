package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error)
	Exec(ctx context.Context, sql string, args ...interface{}) error
}

type sqlQuerier struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlQuerier(db *sql.DB) Querier {
	return &sqlQuerier{db: db}
}

func (q *sqlQuerier) Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	if q.tx != nil {
		return q.tx.QueryContext(ctx, sql, args...)
	}
	return q.db.QueryContext(ctx, sql, args...)
}

func (q *sqlQuerier) Exec(ctx context.Context, sql string, args ...interface{}) error {
	if q.tx != nil {
		_, err := q.tx.ExecContext(ctx, sql, args...)
		return err
	}
	_, err := q.db.ExecContext(ctx, sql, args...)
	return err
}

type UnitOfWork struct {
	withTX   bool
	released *atomic.Bool
	initOnce *sync.Once

	dbConn     *sql.DB
	tx         *sql.Tx
	querier    Querier
	committed  bool
	rolledBack bool
}

// NewUnitOfWork creates a new UnitOfWork instance.
func NewUnitOfWork(dbConn *sql.DB) (*UnitOfWork, error) {
	return &UnitOfWork{
		withTX:     false,
		released:   &atomic.Bool{},
		initOnce:   &sync.Once{},
		committed:  false,
		rolledBack: false,
		dbConn:     dbConn,
	}, nil
}

// WithTX marks the UnitOfWork to use a transaction.
func (uow *UnitOfWork) WithTX() {
	uow.withTX = true
}

// GetQuerier initializes the connection or transaction and returns the querier.
func (uow *UnitOfWork) GetQuerier(ctx context.Context) (Querier, error) {
	if uow.released.Load() {
		return nil, fmt.Errorf("UnitOfWork has already been released")
	}

	var err error
	uow.initOnce.Do(func() {
		if uow.withTX {
			uow.tx, err = uow.dbConn.BeginTx(ctx, nil)
			if err != nil {
				return
			}
			q := NewSqlQuerier(uow.dbConn)
			q.(*sqlQuerier).tx = uow.tx
			uow.querier = q
		} else {
			uow.querier = NewSqlQuerier(uow.dbConn)
		}
	})

	return uow.querier, err
}

// Finalize commits or rolls back the transaction and releases the connection.
func (uow *UnitOfWork) Finalize(ctx context.Context, isSuccess bool) error {
	if uow.released.Load() {
		return nil
	}
	if uow.tx == nil {
		return nil
	}

	defer func() {
		uow.tx.Rollback()
		uow.released.Store(true)
	}()

	if isSuccess {
		return uow.tx.Commit()
	}
	return nil
}

// Context returns the context from gin.Context with transaction
func (uow *UnitOfWork) Context(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if uow.tx != nil {
		return context.WithValue(ctx, "tx", uow.tx)
	}
	return ctx
}

// Rollback rolls back the transaction
func (uow *UnitOfWork) Rollback() error {
	if uow.tx != nil && !uow.committed && !uow.rolledBack {
		uow.rolledBack = true
		return uow.tx.Rollback()
	}
	return nil
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit() error {
	if uow.tx != nil && !uow.committed && !uow.rolledBack {
		uow.committed = true
		return uow.tx.Commit()
	}
	return nil
}
