package value_objects

import "microservice/internal/domain/exceptions"

type Amount struct {
	value float64
}

func NewAmount(a float64) (Amount, error) {
	if a <= 0 {
		return Amount{}, &exceptions.AmountNotValidException{
			Message: "Amount must be greater than zero",
		}
	}
	return Amount{value: a}, nil
}

func (a Amount) Value() float64 {
	return a.value
}
