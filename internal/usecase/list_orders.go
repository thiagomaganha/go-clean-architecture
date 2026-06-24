package usecase

import (
	"context"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
)

type ListOrdersInput struct {
	Page  int
	Limit int
	Query string
}

type ListOrdersOutput struct {
	Orders []entity.Order
	Total  int
	Page   int
	Limit  int
}

type ListOrdersUseCaseInterface interface {
	Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error)
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewListOrdersUseCase(orderRepository entity.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: orderRepository}
}

func (u *ListOrdersUseCase) Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error) {

	if input.Page == 0 {
		input.Page = 1
	}

	if input.Limit == 0 {
		input.Limit = 10
	}

	if input.Query == "" {
		input.Query = "%"
	}

	orders, total, err := u.OrderRepository.List(ctx, input.Page, input.Limit, input.Query)
	if err != nil {
		return ListOrdersOutput{}, err
	}

	return ListOrdersOutput{
		Orders: orders,
		Total:  total,
		Page:   input.Page,
		Limit:  input.Limit,
	}, nil
}
