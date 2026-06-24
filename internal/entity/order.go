package entity

import "fmt"

type Order struct {
	ID         string
	Number     string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, num string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:     id,
		Number: num,
		Price:  price,
		Tax:    tax,
	}

	if err := order.IsValid(); err != nil {
		return nil, err
	}

	order.FinalPrice = order.Price + order.Tax

	return order, nil
}

func (o *Order) IsValid() error {
	if o.ID == "" {
		return &ValidationError{Message: "invalid id"}
	}
	if o.Number == "" {
		return &ValidationError{Message: "invalid number"}
	}
	if o.Price <= 0 {
		return &ValidationError{Message: fmt.Sprintf("invalid price: %v", o.Price)}
	}
	if o.Tax <= 0 {
		return &ValidationError{Message: fmt.Sprintf("invalid tax: %v", o.Tax)}
	}
	return nil
}
