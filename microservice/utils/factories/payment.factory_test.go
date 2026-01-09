package factories

import (
	"testing"
)

func TestNewProcessPaymentConfirmationUseCase(t *testing.T) {
	uc := NewProcessPaymentConfirmationUseCase()
	
	if uc == nil {
		t.Error("Expected use case to be created")
	}
}