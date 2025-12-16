package value_objects

import (
	"testing"
)

func TestNewUnitPrice_ValidPrice(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"positive integer", 10.0, 10.0},
		{"positive decimal", 19.99, 19.99},
		{"small positive", 0.01, 0.01},
		{"large price", 9999.99, 9999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := NewUnitPrice(tt.value)
			if err != nil {
				t.Errorf("NewUnitPrice(%v) unexpected error: %v", tt.value, err)
			}
			if price.Value() != tt.expected {
				t.Errorf("NewUnitPrice(%v).Value() = %v, want %v", tt.value, price.Value(), tt.expected)
			}
		})
	}
}

func TestNewUnitPrice_InvalidPrice(t *testing.T) {
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
			_, err := NewUnitPrice(tt.value)
			if err == nil {
				t.Errorf("NewUnitPrice(%v) expected error, got nil", tt.value)
			}
		})
	}
}
