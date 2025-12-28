package value_objects

import "fmt"

type Name struct {
	value string
}

func NewName(value string) (Name, error) {
	if len(value) < 3 {
		return Name{}, fmt.Errorf("name must be at least 3 characters long")
	}
	if len(value) > 100 {
		return Name{}, fmt.Errorf("name must be at most 100 characters long")
	}
	return Name{value: value}, nil
}

func (n *Name) Value() string {
	return n.value
}
