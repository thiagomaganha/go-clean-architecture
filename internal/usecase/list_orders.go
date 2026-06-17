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

type OrderRepositoryInterface interface {
	List(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error)
	Save(ctx context.Context, order *entity.Order) error
}

type ListOrdersUseCaseInterface interface {
	Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error)
}

type ListOrdersUseCase struct {
	OrderRepository OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository OrderRepositoryInterface) *ListOrdersUseCase {
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

	return u.OrderRepository.List(ctx, input)
}
