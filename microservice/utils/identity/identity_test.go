package identity

import (
	"testing"
)

func TestNewUUIDV4(t *testing.T) {
	uuid1 := NewUUIDV4()
	uuid2 := NewUUIDV4()

	if uuid1 == "" {
		t.Error("NewUUIDV4() returned empty string")
	}
	if uuid1 == uuid2 {
		t.Error("NewUUIDV4() should generate unique UUIDs")
	}
	if !IsValidUUID(uuid1) {
		t.Errorf("NewUUIDV4() generated invalid UUID: %v", uuid1)
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		expected bool
	}{
		{"valid uuid v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid uuid v1", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"valid uuid lowercase", "f47ac10b-58cc-4372-a567-0e02b2c3d479", true},
		{"valid uuid uppercase", "F47AC10B-58CC-4372-A567-0E02B2C3D479", true},
		{"empty string", "", false},
		{"invalid format", "invalid-uuid", false},
		{"too short", "550e8400-e29b-41d4-a716", false},
		{"too long", "550e8400-e29b-41d4-a716-446655440000-extra", false},
		{"without dashes valid", "550e8400e29b41d4a716446655440000", true},
		{"random string", "not-a-uuid-at-all", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUUID(tt.uuid)
			if result != tt.expected {
				t.Errorf("IsValidUUID(%v) = %v, want %v", tt.uuid, result, tt.expected)
			}
		})
	}
}

func TestIsNotValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		expected bool
	}{
		{"valid uuid", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty string", "", true},
		{"invalid format", "invalid-uuid", true},
		{"random string", "not-a-uuid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotValidUUID(tt.uuid)
			if result != tt.expected {
				t.Errorf("IsNotValidUUID(%v) = %v, want %v", tt.uuid, result, tt.expected)
			}
		})
	}
}
