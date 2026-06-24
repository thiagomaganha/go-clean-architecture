package graph

import "github.com/thiagomaganha/go-clean-architecture/internal/usecase"

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCaseInterface
	ListOrdersUseCase  usecase.ListOrdersUseCaseInterface
}
