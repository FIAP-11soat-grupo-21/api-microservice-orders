package gateways

import (
	"errors"
	"testing"
	"time"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
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

func TestNewOrderGateway(t *testing.T) {
	ds := &mockOrderDataSource{}
	gateway := NewOrderGateway(ds)

	if gateway == nil {
		t.Error("NewOrderGateway() returned nil")
	}
}

func createTestOrderEntity() entities.Order {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	now := time.Now()
	order, _ := entities.NewOrderWithItems("order-1", &customerID, 20.0, *status, []entities.OrderItem{*item}, now, nil)
	return *order
}

func TestOrderGateway_Create_Success(t *testing.T) {
	createCalled := false
	ds := &mockOrderDataSource{
		createFunc: func(order daos.OrderDAO) error {
			createCalled = true
			return nil
		},
	}

	gateway := NewOrderGateway(ds)
	order := createTestOrderEntity()

	err := gateway.Create(order)

	if err != nil {
		t.Errorf("Create() unexpected error: %v", err)
	}
	if !createCalled {
		t.Error("Create() should call datasource.Create()")
	}
}

func TestOrderGateway_Create_Error(t *testing.T) {
	ds := &mockOrderDataSource{
		createFunc: func(order daos.OrderDAO) error {
			return errors.New("database error")
		},
	}

	gateway := NewOrderGateway(ds)
	order := createTestOrderEntity()

	err := gateway.Create(order)

	if err == nil {
		t.Error("Create() expected error, got nil")
	}
}

func TestOrderGateway_FindByID_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	ds := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:         "order-1",
				CustomerID: &customerID,
				Amount:     20.0,
				Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
				},
				CreatedAt: now,
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	order, err := gateway.FindByID("order-1")

	if err != nil {
		t.Errorf("FindByID() unexpected error: %v", err)
	}
	if order == nil {
		t.Fatal("FindByID() returned nil")
	}
	if order.ID != "order-1" {
		t.Errorf("FindByID() ID = %v, want order-1", order.ID)
	}
	if *order.CustomerID != customerID {
		t.Errorf("FindByID() CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
	if order.Amount.Value() != 20.0 {
		t.Errorf("FindByID() Amount = %v, want 20.0", order.Amount.Value())
	}
	if len(order.Items) != 1 {
		t.Errorf("FindByID() Items length = %v, want 1", len(order.Items))
	}
}

func TestOrderGateway_FindByID_Error(t *testing.T) {
	ds := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{}, errors.New("not found")
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindByID("order-1")

	if err == nil {
		t.Error("FindByID() expected error, got nil")
	}
}

func TestOrderGateway_FindByID_InvalidStatus(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:        "order-1",
				Amount:    20.0,
				Status:    daos.OrderStatusDAO{ID: "status-1", Name: "ab"}, // Invalid name
				Items:     []daos.OrderItemDAO{},
				CreatedAt: now,
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindByID("order-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid status name, got nil")
	}
}

func TestOrderGateway_FindByID_InvalidItem(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:     "order-1",
				Amount: 20.0,
				Status: daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items: []daos.OrderItemDAO{
					{ID: "item-1", OrderID: "order-1", ProductID: "", Quantity: 2, UnitPrice: 10.0}, // Invalid: empty ProductID
				},
				CreatedAt: now,
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindByID("order-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid item, got nil")
	}
}

func TestOrderGateway_FindByID_InvalidAmount(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return daos.OrderDAO{
				ID:        "order-1",
				Amount:    -10.0, // Invalid amount
				Status:    daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
				Items:     []daos.OrderItemDAO{},
				CreatedAt: now,
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindByID("order-1")

	if err == nil {
		t.Error("FindByID() expected error for invalid amount, got nil")
	}
}

func TestOrderGateway_FindAll_Success(t *testing.T) {
	customerID := "customer-123"
	now := time.Now()

	ds := &mockOrderDataSource{
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

	gateway := NewOrderGateway(ds)
	orders, err := gateway.FindAll(dtos.OrderFilterDTO{})

	if err != nil {
		t.Errorf("FindAll() unexpected error: %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("FindAll() length = %v, want 1", len(orders))
	}
}

func TestOrderGateway_FindAll_Error(t *testing.T) {
	ds := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return nil, errors.New("database error")
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("FindAll() expected error, got nil")
	}
}

func TestOrderGateway_FindAll_InvalidStatus(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:        "order-1",
					Amount:    20.0,
					Status:    daos.OrderStatusDAO{ID: "status-1", Name: "ab"}, // Invalid
					Items:     []daos.OrderItemDAO{},
					CreatedAt: now,
				},
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("FindAll() expected error for invalid status, got nil")
	}
}

func TestOrderGateway_FindAll_InvalidItem(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:     "order-1",
					Amount: 20.0,
					Status: daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
					Items: []daos.OrderItemDAO{
						{ID: "item-1", OrderID: "order-1", ProductID: "", Quantity: 2, UnitPrice: 10.0}, // Invalid
					},
					CreatedAt: now,
				},
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("FindAll() expected error for invalid item, got nil")
	}
}

func TestOrderGateway_FindAll_InvalidOrderAmount(t *testing.T) {
	now := time.Now()

	ds := &mockOrderDataSource{
		findAllFunc: func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
			return []daos.OrderDAO{
				{
					ID:        "order-1",
					Amount:    -10.0, // Invalid amount
					Status:    daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
					Items:     []daos.OrderItemDAO{},
					CreatedAt: now,
				},
			}, nil
		},
	}

	gateway := NewOrderGateway(ds)
	_, err := gateway.FindAll(dtos.OrderFilterDTO{})

	if err == nil {
		t.Error("FindAll() expected error for invalid amount, got nil")
	}
}

func TestOrderGateway_Update_Success(t *testing.T) {
	updateCalled := false
	ds := &mockOrderDataSource{
		updateFunc: func(order daos.OrderDAO) error {
			updateCalled = true
			return nil
		},
	}

	gateway := NewOrderGateway(ds)
	order := createTestOrderEntity()

	err := gateway.Update(order)

	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}
	if !updateCalled {
		t.Error("Update() should call datasource.Update()")
	}
}

func TestOrderGateway_Update_Error(t *testing.T) {
	ds := &mockOrderDataSource{
		updateFunc: func(order daos.OrderDAO) error {
			return errors.New("database error")
		},
	}

	gateway := NewOrderGateway(ds)
	order := createTestOrderEntity()

	err := gateway.Update(order)

	if err == nil {
		t.Error("Update() expected error, got nil")
	}
}

func TestOrderGateway_Delete_Success(t *testing.T) {
	deleteCalled := false
	ds := &mockOrderDataSource{
		deleteFunc: func(id string) error {
			deleteCalled = true
			return nil
		},
	}

	gateway := NewOrderGateway(ds)
	err := gateway.Delete("order-1")

	if err != nil {
		t.Errorf("Delete() unexpected error: %v", err)
	}
	if !deleteCalled {
		t.Error("Delete() should call datasource.Delete()")
	}
}

func TestOrderGateway_Delete_Error(t *testing.T) {
	ds := &mockOrderDataSource{
		deleteFunc: func(id string) error {
			return errors.New("database error")
		},
	}

	gateway := NewOrderGateway(ds)
	err := gateway.Delete("order-1")

	if err == nil {
		t.Error("Delete() expected error, got nil")
	}
}
