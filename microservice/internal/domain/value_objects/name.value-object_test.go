package value_objects

import (
	"strings"
	"testing"
)

func TestNewName_ValidName(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"minimum length", "abc", "abc"},
		{"normal name", "Pending", "Pending"},
		{"long name", "Order Status Name", "Order Status Name"},
		{"exactly 100 chars", strings.Repeat("a", 100), strings.Repeat("a", 100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := NewName(tt.value)
			if err != nil {
				t.Errorf("NewName(%v) unexpected error: %v", tt.value, err)
			}
			if name.Value() != tt.expected {
				t.Errorf("NewName(%v).Value() = %v, want %v", tt.value, name.Value(), tt.expected)
			}
		})
	}
}

func TestNewName_InvalidName(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectedErr string
	}{
		{"empty", "", "at least 3 characters"},
		{"too short", "ab", "at least 3 characters"},
		{"single char", "a", "at least 3 characters"},
		{"too long", strings.Repeat("a", 101), "at most 100 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewName(tt.value)
			if err == nil {
				t.Errorf("NewName(%v) expected error, got nil", tt.value)
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("NewName(%v) error = %v, want error containing %v", tt.value, err, tt.expectedErr)
			}
		})
	}
}
