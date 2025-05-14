package repository

import (
	"context"
	"database/sql"
)

type OrderRepository interface {
	DeleteOrder(ctx context.Context, orderID int) error
}

type MySQLOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &MySQLOrderRepository{db: db}
}

func (r *MySQLOrderRepository) DeleteOrder(ctx context.Context, orderID int) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		// If no transaction in context, use the database connection
		_, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = ?", orderID)
		return err
	}

	_, err := tx.ExecContext(ctx, "DELETE FROM orders WHERE id = ?", orderID)
	return err
}
