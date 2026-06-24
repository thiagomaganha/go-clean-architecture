package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
)

type WebOrderHandler struct {
	ListOrdersUseCase  usecase.ListOrdersUseCaseInterface
	CreateOrderUseCase usecase.CreateOrderUseCaseInterface
}

func NewWebOrderHandler(listOrdersUseCase usecase.ListOrdersUseCaseInterface, createOrderUseCase usecase.CreateOrderUseCaseInterface) *WebOrderHandler {
	return &WebOrderHandler{
		ListOrdersUseCase:  listOrdersUseCase,
		CreateOrderUseCase: createOrderUseCase,
	}
}

func (o *WebOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	dto := usecase.ListOrdersInput{
		Page:  page,
		Limit: limit,
		Query: query.Get("query"),
	}

	output, err := o.ListOrdersUseCase.Execute(r.Context(), dto)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, err error) {
	var valErr *entity.ValidationError
	if errors.As(err, &valErr) {
		http.Error(w, valErr.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func (o *WebOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	output, err := o.CreateOrderUseCase.Execute(r.Context(), input)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
