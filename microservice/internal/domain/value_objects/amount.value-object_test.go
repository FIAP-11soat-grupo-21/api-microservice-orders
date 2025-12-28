package value_objects

import (
	"testing"

	"microservice/internal/domain/exceptions"
)

func TestNewAmount_ValidAmount(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"positive integer", 100.0, 100.0},
		{"positive decimal", 99.99, 99.99},
		{"small positive", 0.01, 0.01},
		{"large amount", 999999.99, 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := NewAmount(tt.value)
			if err != nil {
				t.Errorf("NewAmount(%v) unexpected error: %v", tt.value, err)
			}
			if amount.Value() != tt.expected {
				t.Errorf("NewAmount(%v).Value() = %v, want %v", tt.value, amount.Value(), tt.expected)
			}
		})
	}
}

func TestNewAmount_InvalidAmount(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"negative", -10.0},
		{"negative decimal", -0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAmount(tt.value)
			if err == nil {
				t.Errorf("NewAmount(%v) expected error, got nil", tt.value)
			}
			if _, ok := err.(*exceptions.AmountNotValidException); !ok {
				t.Errorf("NewAmount(%v) expected AmountNotValidException, got %T", tt.value, err)
			}
		})
	}
}
