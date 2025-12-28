package value_objects

import "microservice/internal/domain/exceptions"

type ProductID struct {
	value string
}

func NewProductID(value string) (ProductID, error) {
	if value == "" {
		return ProductID{}, &exceptions.InvalidOrderItemData{
			Message: "Product ID cannot be empty",
		}
	}
	return ProductID{value: value}, nil
}

func (p ProductID) Value() string {
	return p.value
}
