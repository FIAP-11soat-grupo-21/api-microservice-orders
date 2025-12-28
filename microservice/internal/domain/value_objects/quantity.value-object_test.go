package value_objects

import (
	"testing"

	"microservice/internal/domain/exceptions"
)

func TestNewQuantity_ValidQuantity(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"one item", 1, 1},
		{"multiple items", 5, 5},
		{"large quantity", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quantity, err := NewQuantity(tt.value)
			if err != nil {
				t.Errorf("NewQuantity(%v) unexpected error: %v", tt.value, err)
			}
			if quantity.Value() != tt.expected {
				t.Errorf("NewQuantity(%v).Value() = %v, want %v", tt.value, quantity.Value(), tt.expected)
			}
		})
	}
}

func TestNewQuantity_InvalidQuantity(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{"zero", 0},
		{"negative", -1},
		{"large negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQuantity(tt.value)
			if err == nil {
				t.Errorf("NewQuantity(%v) expected error, got nil", tt.value)
			}
			if _, ok := err.(*exceptions.InvalidOrderItemData); !ok {
				t.Errorf("NewQuantity(%v) expected InvalidOrderItemData, got %T", tt.value, err)
			}
		})
	}
}
