package db

import (
    "context"
    "database/sql"
    "sync"
)

type UnitOfWork struct {
    tx       *sql.Tx
    released bool
    mu       sync.Mutex
}

func NewUnitOfWork(db *sql.DB) (*UnitOfWork, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }

    return &UnitOfWork{
        tx:       tx,
        released: false,
    }, nil
}

func (u *UnitOfWork) Commit() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.released {
        return nil
    }

    err := u.tx.Commit()
    u.released = true
    return err
}

func (u *UnitOfWork) Rollback() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.released {
        return nil
    }

    err := u.tx.Rollback()
    u.released = true
    return err
}

func (u *UnitOfWork) Context(ctx context.Context) context.Context {
    return context.WithValue(ctx, "tx", u.tx)
}
