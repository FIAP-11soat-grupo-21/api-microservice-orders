package value_objects

import (
	"testing"

	"microservice/internal/domain/exceptions"
)

func TestNewProductID_ValidProductID(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"uuid format", "550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "product-123", "product-123"},
		{"numeric id", "12345", "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID, err := NewProductID(tt.value)
			if err != nil {
				t.Errorf("NewProductID(%v) unexpected error: %v", tt.value, err)
			}
			if productID.Value() != tt.expected {
				t.Errorf("NewProductID(%v).Value() = %v, want %v", tt.value, productID.Value(), tt.expected)
			}
		})
	}
}

func TestNewProductID_InvalidProductID(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProductID(tt.value)
			if err == nil {
				t.Errorf("NewProductID(%v) expected error, got nil", tt.value)
			}
			if _, ok := err.(*exceptions.InvalidOrderItemData); !ok {
				t.Errorf("NewProductID(%v) expected InvalidOrderItemData, got %T", tt.value, err)
			}
		})
	}
}
