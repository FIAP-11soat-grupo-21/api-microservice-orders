package value_objects

import "fmt"

type UnitPrice struct {
	value float64
}

func NewUnitPrice(value float64) (UnitPrice, error) {
	if value <= 0 {
		return UnitPrice{}, fmt.Errorf("unit price must be greater than 0")
	}
	return UnitPrice{value: value}, nil
}

func (p *UnitPrice) Value() float64 {
	return p.value
}
