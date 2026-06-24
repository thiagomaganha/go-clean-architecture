package entity

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	List(ctx context.Context, page, limit int, query string) ([]Order, int, error)
}
