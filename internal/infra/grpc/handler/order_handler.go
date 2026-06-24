package handler

import (
	"context"
	"errors"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/grpc/pb"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCaseInterface
	ListOrdersUseCase  usecase.ListOrdersUseCaseInterface
}

func NewOrderGrpcHandler(
	createOrderUseCase usecase.CreateOrderUseCaseInterface,
	listOrdersUseCase usecase.ListOrdersUseCaseInterface,
) *OrderGrpcHandler {
	return &OrderGrpcHandler{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (h *OrderGrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	input := usecase.CreateOrderInput{
		ID:     req.Id,
		Number: req.Number,
		Price:  req.Price,
		Tax:    req.Tax,
	}

	output, err := h.CreateOrderUseCase.Execute(ctx, input)
	if err != nil {
		return nil, grpcError(err)
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Number:     output.Number,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (h *OrderGrpcHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	input := usecase.ListOrdersInput{
		Page:  int(req.Page),
		Limit: int(req.Limit),
		Query: req.Query,
	}

	output, err := h.ListOrdersUseCase.Execute(ctx, input)
	if err != nil {
		return nil, grpcError(err)
	}

	orders := make([]*pb.Order, len(output.Orders))
	for i, o := range output.Orders {
		orders[i] = &pb.Order{
			Id:         o.ID,
			Number:     o.Number,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		}
	}

	return &pb.ListOrdersResponse{
		Orders: orders,
		Total:  int32(output.Total),
		Page:   int32(output.Page),
		Limit:  int32(output.Limit),
	}, nil
}

func grpcError(err error) error {
	var valErr *entity.ValidationError
	if errors.As(err, &valErr) {
		return status.Error(codes.InvalidArgument, valErr.Error())
	}
	return status.Error(codes.Internal, "internal server error")
}
