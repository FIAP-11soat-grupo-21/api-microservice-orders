package controllers

import (
	"errors"
	"testing"
	"time"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type mockOrderDataSource struct {
	createFunc   func(order daos.OrderDAO) error
	findAllFunc  func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error)
	findByIDFunc func(id string) (daos.OrderDAO, error)
	updateFunc   func(order daos.OrderDAO) error
	deleteFunc   func(id string) error
}

func (m *mockOrderDataSource) Create(order daos.OrderDAO) error {
	if m.createFunc != nil {
		return m.createFunc(order)
	}
	return nil
}

func (m *mockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return []daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (m *mockOrderDataSource) Update(order daos.OrderDAO) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func (m *mockOrderDataSource) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

type mockOrderStatusDataSource struct {
	findByIDFunc func(id string) (daos.OrderStatusDAO, error)
	findAllFunc  func() ([]daos.OrderStatusDAO, error)
}

func (m *mockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderStatusDAO{}, nil
}

func (m *mockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}

func TestNewOrderController(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	if controller == nil {
		t.Error("NewOrderController() returned nil")
	}
}

func TestOrderController_Create_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{
		createFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}, nil
		},
	}

	controller := NewOrderController(orderDS, statusDS)

	dto := dtos.CreateOrderDTO{
		CustomerID: nil,
		Items: []dtos.CreateOrderItemDTO{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}

	response, err := controller.Create(dto)

	if err != nil {
		t.Errorf("Create() unexpected error: %v", err)
	}
	if response.ID == "" {
		t.Error("Create() returned empty ID")
	}
}

func TestOrderController_Create_StatusNotFound(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, errors.New("not found")
		},
	}

	controller := NewOrderController(orderDS, statusDS)

	dto := dtos.CreateOrderDTO{
		Items: []dtos.CreateOrderItemDTO{
			{ProductID: "product-1", Quantity: 2, Price: 10.0},
		},
	}

	_, err := controller.Create(dto)

	if err == nil {
		t.Error("Create() expected error when status not found, got nil")
	}
}

func TestOrderController_FindAll_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:         "order-1",
					CustomerID: &customerID,
					Amount:     20.0,
					Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
					Items: []daos.OrderItemDAO{
						{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
					},
					CreatedAt: now,
				},
			}, nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	responses, err := controller.FindAll(dtos.OrderFilterDTO{})

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Errorf("FindAll() length = %v, want 1", len(responses))
	}
}

func TestOrderController_FindAll_Error(t *testing.T) {
	orderDS := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return nil, errors.New("database error")
		},
	}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	_, err := controller.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("FindAll() expected error, got nil")
	}
}

func TestOrderController_FindByID_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "550e8400-e29b-41d4-a716-446655440000", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	response, err := controller.FindByID("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}
	if response.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("FindByID() ID = %v, want 550e8400-e29b-41d4-a716-446655440000", response.ID)
	}
}

func TestOrderController_FindByID_InvalidID(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	_, err := controller.FindByID("invalid-id")

	if err == nil {
		t.Error("FindByID() expected error for invalid ID, got nil")
	}
}

func TestOrderController_Update_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "550e8400-e29b-41d4-a716-446655440000", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
		updateFunc: func(order daos.OrderDAO) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{ID: "status-2", Name: "Confirmed"}, nil
		},
	}

	controller := NewOrderController(orderDS, statusDS)

	dto := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: "status-2",
	}

	response, err := controller.Update(dto)

	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}
	if response.Status.ID != "status-2" {
		t.Errorf("Update() Status.ID = %v, want status-2", response.Status.ID)
	}
}

func TestOrderController_Update_InvalidID(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	dto := dtos.UpdateOrderDTO{
		ID:       "invalid-id",
		StatusID: "status-2",
	}

	_, err := controller.Update(dto)

	if err == nil {
		t.Error("Update() expected error for invalid ID, got nil")
	}
}

func TestOrderController_Delete_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	orderDS := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "550e8400-e29b-41d4-a716-446655440000",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items:      []daos.OrderItemDAO{},
				CreatedAt:  now,
			}, nil
		},
		deleteFunc: func(id string) error {
			return nil
		},
	}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	err := controller.Delete("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("Delete() unexpected error: %v", err)
	}
}

func TestOrderController_Delete_InvalidID(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{}

	controller := NewOrderController(orderDS, statusDS)

	err := controller.Delete("invalid-id")

	if err == nil {
		t.Error("Delete() expected error for invalid ID, got nil")
	}
}

func TestOrderController_FindAllStatus_Success(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return []daos.OrderStatusDAO{
				{ID: "status-1", Name: "Pending"},
				{ID: "status-2", Name: "Confirmed"},
			}, nil
		},
	}

	controller := NewOrderController(orderDS, statusDS)

	responses, err := controller.FindAllStatus()

	if err != nil {
		t.Errorf("FindAllStatus() unexpected error: %v", err)
	}
	if len(responses) != 2 {
		t.Errorf("FindAllStatus() length = %v, want 2", len(responses))
	}
}

func TestOrderController_FindAllStatus_Error(t *testing.T) {
	orderDS := &mockOrderDataSource{}
	statusDS := &mockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return nil, errors.New("database error")
		},
	}

	controller := NewOrderController(orderDS, statusDS)

	_, err := controller.FindAllStatus()

	if err == nil {
		t.Error("FindAllStatus() expected error, got nil")
	}
}
