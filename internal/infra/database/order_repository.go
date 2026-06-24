package database

import (
	"context"
	"database/sql"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(ctx context.Context, order *entity.Order) error {
	_, err := r.Db.ExecContext(
		ctx,
		"INSERT INTO orders (id, number, price, tax, final_price) VALUES (?, ?, ?, ?, ?)",
		order.ID, order.Number, order.Price, order.Tax, order.FinalPrice,
	)
	return err
}

func (o *OrderRepository) List(ctx context.Context, page int, limit int, query string) ([]entity.Order, int, error) {
	var total int
	if err := o.Db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM orders WHERE number LIKE ?",
		query,
	).Scan(&total); err != nil {
		return []entity.Order{}, 0, err
	}

	var orders []entity.Order
	offset := (page - 1) * limit
	rows, err := o.Db.QueryContext(
		ctx,
		"SELECT id, number, price, tax, final_price FROM orders WHERE number LIKE ? LIMIT ? OFFSET ?",
		query, limit, offset,
	)
	if err != nil {
		return []entity.Order{}, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Number, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return []entity.Order{}, 0, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return []entity.Order{}, 0, err
	}

	return orders, total, nil
}
