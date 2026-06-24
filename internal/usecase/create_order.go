package usecase

import (
	"context"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
)

type CreateOrderInput struct {
	ID     string
	Number string
	Price  float64
	Tax    float64
}

type CreateOrderOutput struct {
	ID         string
	Number     string
	Price      float64
	Tax        float64
	FinalPrice float64
}

type CreateOrderUseCaseInterface interface {
	Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error)
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{OrderRepository: orderRepository}
}

func (u *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	order, err := entity.NewOrder(input.ID, input.Number, input.Price, input.Tax)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	if err := u.OrderRepository.Save(ctx, order); err != nil {
		return CreateOrderOutput{}, err
	}

	return CreateOrderOutput{
		ID:         order.ID,
		Number:     order.Number,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
