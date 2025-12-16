package value_objects

import (
	"fmt"

	"microservice/internal/domain/exceptions"
)

type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value <= 0 {
		return Quantity{}, &exceptions.InvalidOrderItemData{
			Message: fmt.Sprintf("quantity must be greater than zero, got %d", value),
		}
	}
	return Quantity{value: value}, nil
}

func (q Quantity) Value() int {
	return q.value
}
