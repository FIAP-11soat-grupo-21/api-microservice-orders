package exceptions

import (
	"testing"
)

func TestInvalidOrderItemData_Error(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"with custom message", "Quantity must be positive", "Quantity must be positive"},
		{"with empty message", "", "Invalid order item data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &InvalidOrderItemData{Message: tt.message}
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}
