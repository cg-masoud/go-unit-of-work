package repository

import (
    "context"
    "database/sql"
)

type ItemRepository interface {
    DeleteItemsByOrderID(ctx context.Context, orderID int) error
}

type MySQLItemRepository struct {
    db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
    return &MySQLItemRepository{db: db}
}

func (r *MySQLItemRepository) DeleteItemsByOrderID(ctx context.Context, orderID int) error {
    tx, ok := ctx.Value("tx").(*sql.Tx)
    if !ok {
        return sql.ErrTxDone
    }

    _, err := tx.ExecContext(ctx, "DELETE FROM items WHERE order_id = ?", orderID)
    return err
}
