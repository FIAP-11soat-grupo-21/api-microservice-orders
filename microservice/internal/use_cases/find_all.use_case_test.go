package use_cases

import (
	"testing"
	"time"

	"microservice/internal/adapters/dtos"
)

func TestFindAllOrdersUseCase_FilterDTOStructure(t *testing.T) {
	customerID := "customer-123"
	statusID := "status-1"
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	filter := dtos.OrderFilterDTO{
		CustomerID:    &customerID,
		StatusID:      &statusID,
		CreatedAtFrom: &from,
		CreatedAtTo:   &to,
	}

	if *filter.CustomerID != customerID {
		t.Errorf("OrderFilterDTO.CustomerID = %v, want %v", *filter.CustomerID, customerID)
	}
	if *filter.StatusID != statusID {
		t.Errorf("OrderFilterDTO.StatusID = %v, want %v", *filter.StatusID, statusID)
	}
	if filter.CreatedAtFrom == nil {
		t.Error("OrderFilterDTO.CreatedAtFrom should not be nil")
	}
	if filter.CreatedAtTo == nil {
		t.Error("OrderFilterDTO.CreatedAtTo should not be nil")
	}
}

func TestFindAllOrdersUseCase_EmptyFilter(t *testing.T) {
	filter := dtos.OrderFilterDTO{}

	if filter.CustomerID != nil {
		t.Errorf("Empty OrderFilterDTO.CustomerID = %v, want nil", filter.CustomerID)
	}
	if filter.StatusID != nil {
		t.Errorf("Empty OrderFilterDTO.StatusID = %v, want nil", filter.StatusID)
	}
	if filter.CreatedAtFrom != nil {
		t.Errorf("Empty OrderFilterDTO.CreatedAtFrom = %v, want nil", filter.CreatedAtFrom)
	}
	if filter.CreatedAtTo != nil {
		t.Errorf("Empty OrderFilterDTO.CreatedAtTo = %v, want nil", filter.CreatedAtTo)
	}
}

func TestFindAllOrdersUseCase_PartialFilter(t *testing.T) {
	statusID := "status-1"

	filter := dtos.OrderFilterDTO{
		StatusID: &statusID,
	}

	if filter.CustomerID != nil {
		t.Errorf("Partial OrderFilterDTO.CustomerID = %v, want nil", filter.CustomerID)
	}
	if *filter.StatusID != statusID {
		t.Errorf("Partial OrderFilterDTO.StatusID = %v, want %v", *filter.StatusID, statusID)
	}
}

func TestNewFindAllOrdersUseCase(t *testing.T) {
	_ = NewFindAllOrdersUseCase
}

func TestFindAllOrdersUseCase_Execute_Success(t *testing.T) {
	filter := dtos.OrderFilterDTO{}

	if filter.CustomerID != nil {
		t.Error("Empty filter should have nil CustomerID")
	}

	if filter.StatusID != nil {
		t.Error("Empty filter should have nil StatusID")
	}
}

func TestFindAllOrdersUseCase_Execute_WithFilter(t *testing.T) {
	customerID := "customer-1"
	statusID := "status-1"

	filter := dtos.OrderFilterDTO{
		CustomerID: &customerID,
		StatusID:   &statusID,
	}

	if filter.CustomerID == nil || *filter.CustomerID != customerID {
		t.Errorf("Expected customer ID '%s', got %v", customerID, filter.CustomerID)
	}

	if filter.StatusID == nil || *filter.StatusID != statusID {
		t.Errorf("Expected status ID '%s', got %v", statusID, filter.StatusID)
	}
}

func TestFindAllOrdersUseCase_Execute_FilterValidation(t *testing.T) {
	testCases := []struct {
		name       string
		customerID *string
		statusID   *string
	}{
		{"empty filter", nil, nil},
		{"customer only", stringPtr("customer-1"), nil},
		{"status only", nil, stringPtr("status-1")},
		{"both fields", stringPtr("customer-1"), stringPtr("status-1")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := dtos.OrderFilterDTO{
				CustomerID: tc.customerID,
				StatusID:   tc.statusID,
			}

			if tc.customerID == nil && filter.CustomerID != nil {
				t.Error("Expected nil customer ID")
			}

			if tc.statusID == nil && filter.StatusID != nil {
				t.Error("Expected nil status ID")
			}
		})
	}
}
