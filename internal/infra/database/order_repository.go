package database

import (
	"context"
	"database/sql"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
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

func (o *OrderRepository) List(ctx context.Context, input usecase.ListOrdersInput) (usecase.ListOrdersOutput, error) {
	var total int
	if err := o.Db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM orders WHERE number LIKE ?",
		input.Query,
	).Scan(&total); err != nil {
		return usecase.ListOrdersOutput{}, err
	}

	var orders []entity.Order
	offset := (input.Page - 1) * input.Limit
	rows, err := o.Db.QueryContext(
		ctx,
		"SELECT id, number, price, tax, final_price FROM orders WHERE number LIKE ? LIMIT ? OFFSET ?",
		input.Query, input.Limit, offset,
	)
	if err != nil {
		return usecase.ListOrdersOutput{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Number, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return usecase.ListOrdersOutput{}, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return usecase.ListOrdersOutput{}, err
	}

	return usecase.ListOrdersOutput{
		Orders: orders,
		Total:  total,
		Page:   input.Page,
		Limit:  input.Limit,
	}, nil
}
