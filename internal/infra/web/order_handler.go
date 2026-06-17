package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
)

type WebOrderHandler struct {
	ListOrdersUseCase usecase.ListOrdersUseCaseInterface
}

func NewWebOrderHandler(listOrdersUseCase usecase.ListOrdersUseCaseInterface) *WebOrderHandler {
	return &WebOrderHandler{ListOrdersUseCase: listOrdersUseCase}
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
		http.Error(w, "failed to list orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
