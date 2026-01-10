package gateways

import (
	"errors"
	"testing"

	"microservice/internal/adapters/daos"
	"microservice/internal/test_helpers"
)

func TestNewOrderStatusGateway(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{}
	gateway := NewOrderStatusGateway(ds)

	if gateway == nil {
		t.Error("NewOrderStatusGateway() returned nil")
	}
}

func TestOrderStatusGateway_FindAll_Success(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(statuses) != 2 {
		t.Errorf("FindAll() length = %v, want 2", len(statuses))
	}
	if statuses[0].ID != "status-1" {
		t.Errorf("FindAll()[0].ID = %v, want status-1", statuses[0].ID)
	}
	if statuses[0].Name.Value() != "Pending" {
		t.Errorf("FindAll()[0].Name = %v, want Pending", statuses[0].Name.Value())
	}
}

func TestOrderStatusGateway_FindAll_Error(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return nil, errors.New("database error")
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error, got nil")
	}
}

func TestOrderStatusGateway_FindAll_InvalidName(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "ab"}, // Invalid: too short
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error for invalid name, got nil")
	}
}

func TestOrderStatusGateway_FindByID_Success(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	status, err := gateway.FindByID("status-1")

	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}
	if status == nil {
		t.Fatal("FindByID() returned nil")
	}
	if status.ID != "status-1" {
		t.Errorf("FindByID() ID = %v, want status-1", status.ID)
	}
	if status.Name.Value() != "Pending" {
		t.Errorf("FindByID() Name = %v, want Pending", status.Name.Value())
	}
}

func TestOrderStatusGateway_FindByID_Error(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("not found")
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("status-1")

	if err == nil {
		t.Error("FindByID() expected error, got nil")
	}
}

func TestOrderStatusGateway_FindByID_InvalidName(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "ab"}, nil // Invalid name
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("status-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid name, got nil")
	}
}


func TestOrderStatusGateway_FindAll_Empty(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(statuses) != 0 {
		t.Errorf("FindAll() length = %v, want 0", len(statuses))
	}
}

func TestOrderStatusGateway_FindAll_MultipleStatuses(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
				{ID: "status-3", Name: "Preparing"},
				{ID: "status-4", Name: "Ready"},
				{ID: "status-5", Name: "Delivered"},
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(statuses) != 5 {
		t.Errorf("FindAll() length = %v, want 5", len(statuses))
	}
	for i, status := range statuses {
		if status.ID != "status-"+string(rune(i+1+'0')) {
			t.Errorf("FindAll()[%d].ID = %v, want status-%d", i, status.ID, i+1)
		}
	}
}

func TestOrderStatusGateway_FindByID_DifferentIDs(t *testing.T) {
	tests := []struct {
		name     string
		statusID string
		statusName string
	}{
		{"Pending", "status-1", "Pending"},
		{"Confirmed", "status-2", "Confirmed"},
		{"Preparing", "status-3", "Preparing"},
		{"Ready", "status-4", "Ready"},
		{"Delivered", "status-5", "Delivered"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &test_helpers.MockOrderStatusDataSource{
				FindByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
					return daos.OrderStatusDAO{ID: tt.statusID, Name: tt.statusName}, nil
				},
			}

			gateway := NewOrderStatusGateway(ds)
			status, err := gateway.FindByID(tt.statusID)

			if err != nil {
				t.Errorf("FindByID() unexpected error: %v", err)
			}
			if status == nil {
				t.Fatal("FindByID() returned nil")
			}
			if status.ID != tt.statusID {
				t.Errorf("FindByID() ID = %v, want %v", status.ID, tt.statusID)
			}
			if status.Name.Value() != tt.statusName {
				t.Errorf("FindByID() Name = %v, want %v", status.Name.Value(), tt.statusName)
			}
		})
	}
}

func TestOrderStatusGateway_FindAll_PartialInvalidName(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "ab"}, // Invalid: too short
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindAll()

	if err == nil {
		t.Error("FindAll() expected error for invalid name in list, got nil")
	}
}

func TestOrderStatusGateway_FindByID_EmptyID(t *testing.T) {
	ds := &test_helpers.MockOrderStatusDataSource{
		FindByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			if id == "" {
				return daos.OrderStatusDAO{}, errors.New("empty id")
			}
			return daos.OrderStatusDAO{ID: id, Name: "Pending"}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	_, err := gateway.FindByID("")

	if err == nil {
		t.Error("FindByID() expected error for empty id, got nil")
	}
}

func TestOrderStatusGateway_FindAll_LongStatusName(t *testing.T) {
	longName := "VeryLongStatusNameThatIsStillValid"
	ds := &test_helpers.MockOrderStatusDataSource{
		FindAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: longName},
			}, nil
		},
	}

	gateway := NewOrderStatusGateway(ds)
	statuses, err := gateway.FindAll()

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(statuses) != 1 {
		t.Errorf("FindAll() length = %v, want 1", len(statuses))
	}
	if statuses[0].Name.Value() != longName {
		t.Errorf("FindAll()[0].Name = %v, want %v", statuses[0].Name.Value(), longName)
	}
}
