package use_cases

import (
	"errors"
	"testing"
	"time"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	"microservice/internal/adapters/gateways"
	"microservice/internal/domain/entities"
)

type testOrderDataSource struct {
	orders map[string]daos.OrderDAO
}

func newTestOrderDataSource() *testOrderDataSource {
	return &testOrderDataSource{
		orders: make(map[string]daos.OrderDAO),
	}
}

func (ds *testOrderDataSource) Create(order daos.OrderDAO) error {
	ds.orders[order.ID] = order
	return nil
}

func (ds *testOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	result := make([]daos.OrderDAO, 0, len(ds.orders))
	for _, order := range ds.orders {
		result = append(result, order)
	}
	return result, nil
}

func (ds *testOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	order, ok := ds.orders[id]
	if !ok {
		return daos.OrderDAO{}, errors.New("not found")
	}
	return order, nil
}

func (ds *testOrderDataSource) Update(order daos.OrderDAO) error {
	ds.orders[order.ID] = order
	return nil
}

func (ds *testOrderDataSource) Delete(id string) error {
	delete(ds.orders, id)
	return nil
}

type testOrderStatusDataSource struct {
	statuses map[string]daos.OrderStatusDAO
}

func newTestOrderStatusDataSource() *testOrderStatusDataSource {
	return &testOrderStatusDataSource{
		statuses: map[string]daos.OrderStatusDAO{
			INITIAL_ORDER_STATUS_ID: {ID: INITIAL_ORDER_STATUS_ID, Name: "Pending"},
			"status-2":              {ID: "status-2", Name: "Confirmed"},
		},
	}
}

func (ds *testOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	status, ok := ds.statuses[id]
	if !ok {
		return daos.OrderStatusDAO{}, errors.New("not found")
	}
	return status, nil
}

func (ds *testOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	result := make([]daos.OrderStatusDAO, 0, len(ds.statuses))
	for _, status := range ds.statuses {
		result = append(result, status)
	}
	return result, nil
}

func TestCreateOrderUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 2, Price: 10.0},
		{ProductID: "product-2", Quantity: 1, Price: 25.0},
	}

	order, err := uc.Execute(&customerID, items)

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
	if order.IsEmpty() {
		t.Error("Execute() returned empty order")
	}
	if *order.CustomerID != customerID {
		t.Errorf("Execute() CustomerID = %v, want %v", *order.CustomerID, customerID)
	}
	if len(order.Items) != 2 {
		t.Errorf("Execute() Items length = %v, want 2", len(order.Items))
	}

	expectedAmount := (2 * 10.0) + (1 * 25.0)
	if order.Amount.Value() != expectedAmount {
		t.Errorf("Execute() Amount = %v, want %v", order.Amount.Value(), expectedAmount)
	}
}

func TestFindAllOrdersUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	_, _ = createUC.Execute(&customerID, items)
	_, _ = createUC.Execute(&customerID, items)

	findAllUC := NewFindAllOrdersUseCase(*orderGateway)
	orders, err := findAllUC.Execute(dtos.OrderFilterDTO{})

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("Execute() returned %d orders, want 2", len(orders))
	}
}

func TestFindOrderByIDUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	findByIDUC := NewFindOrderByIDUseCase(*orderGateway)
	order, err := findByIDUC.Execute(createdOrder.ID)

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
	if order.ID != createdOrder.ID {
		t.Errorf("Execute() ID = %v, want %v", order.ID, createdOrder.ID)
	}
}

func TestUpdateOrderUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	updateUC := NewUpdateOrderUseCase(*orderGateway, *statusGateway)
	dto := dtos.UpdateOrderDTO{
		ID:       createdOrder.ID,
		StatusID: "status-2",
	}
	updatedOrder, err := updateUC.Execute(dto)

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
	if updatedOrder.Status.ID != "status-2" {
		t.Errorf("Execute() Status.ID = %v, want status-2", updatedOrder.Status.ID)
	}
	if updatedOrder.UpdatedAt == nil {
		t.Error("Execute() UpdatedAt should not be nil after update")
	}
}

func TestDeleteOrderUseCase_Execute_Integration(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	deleteUC := NewDeleteOrderUseCase(*orderGateway)
	err := deleteUC.Execute(createdOrder.ID)

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}

	findByIDUC := NewFindOrderByIDUseCase(*orderGateway)
	_, err = findByIDUC.Execute(createdOrder.ID)
	if err == nil {
		t.Error("Execute() order should be deleted")
	}
}

func TestFindAllOrderStatusUseCase_Execute_Integration(t *testing.T) {
	statusDS := newTestOrderStatusDataSource()
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewFindAllOrderStatusUseCase(*statusGateway)
	statuses, err := uc.Execute()

	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}
	if len(statuses) != 2 {
		t.Errorf("Execute() returned %d statuses, want 2", len(statuses))
	}
}

func createTestOrderForUseCases(id string, customerID *string) *entities.Order {
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", id, 2, 10.0)
	now := time.Now()
	order, _ := entities.NewOrderWithItems(id, customerID, 20.0, *status, []entities.OrderItem{*item}, now, nil)
	return order
}


type errorOrderDataSource struct{}

func (ds *errorOrderDataSource) Create(order daos.OrderDAO) error {
	return errors.New("create error")
}

func (ds *errorOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return nil, errors.New("find all error")
}

func (ds *errorOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	return daos.OrderDAO{}, errors.New("find by id error")
}

func (ds *errorOrderDataSource) Update(order daos.OrderDAO) error {
	return errors.New("update error")
}

func (ds *errorOrderDataSource) Delete(id string) error {
	return errors.New("delete error")
}

type errorOrderStatusDataSource struct{}

func (ds *errorOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	return daos.OrderStatusDAO{}, errors.New("status not found")
}

func (ds *errorOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	return nil, errors.New("find all error")
}

func TestCreateOrderUseCase_Execute_StatusNotFound(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := &errorOrderStatusDataSource{}

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Execute() expected error when status not found, got nil")
	}
}

func TestCreateOrderUseCase_Execute_InvalidItem(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "", Quantity: 1, Price: 10.0}, // Invalid: empty product ID
	}

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Execute() expected error for invalid item, got nil")
	}
}

func TestCreateOrderUseCase_Execute_GatewayError(t *testing.T) {
	orderDS := &errorOrderDataSource{}
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Execute() expected error when gateway fails, got nil")
	}
}

func TestFindOrderByIDUseCase_Execute_InvalidID(t *testing.T) {
	orderDS := newTestOrderDataSource()
	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewFindOrderByIDUseCase(*orderGateway)

	_, err := uc.Execute("invalid-uuid")

	if err == nil {
		t.Error("Execute() expected error for invalid UUID, got nil")
	}
}

func TestFindOrderByIDUseCase_Execute_NotFound(t *testing.T) {
	orderDS := newTestOrderDataSource()
	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewFindOrderByIDUseCase(*orderGateway)

	_, err := uc.Execute("550e8400-e29b-41d4-a716-446655440000")

	if err == nil {
		t.Error("Execute() expected error when order not found, got nil")
	}
}

func TestUpdateOrderUseCase_Execute_InvalidID(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewUpdateOrderUseCase(*orderGateway, *statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       "invalid-uuid",
		StatusID: "status-2",
	}

	_, err := uc.Execute(dto)

	if err == nil {
		t.Error("Execute() expected error for invalid UUID, got nil")
	}
}

func TestUpdateOrderUseCase_Execute_OrderNotFound(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewUpdateOrderUseCase(*orderGateway, *statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: "status-2",
	}

	_, err := uc.Execute(dto)

	if err == nil {
		t.Error("Execute() expected error when order not found, got nil")
	}
}

func TestUpdateOrderUseCase_Execute_StatusNotFound(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	errorStatusDS := &errorOrderStatusDataSource{}
	errorStatusGateway := gateways.NewOrderStatusGateway(errorStatusDS)

	updateUC := NewUpdateOrderUseCase(*orderGateway, *errorStatusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       createdOrder.ID,
		StatusID: "invalid-status",
	}

	_, err := updateUC.Execute(dto)

	if err == nil {
		t.Error("Execute() expected error when status not found, got nil")
	}
}

func TestUpdateOrderUseCase_Execute_GatewayError(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	hybridDS := &hybridOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return orderDS.FindByID(id)
		},
		updateFunc: func(order daos.OrderDAO) error {
			return errors.New("update error")
		},
	}
	errorOrderGateway := gateways.NewOrderGateway(hybridDS)

	updateUC := NewUpdateOrderUseCase(*errorOrderGateway, *statusGateway)

	dto := dtos.UpdateOrderDTO{
		ID:       createdOrder.ID,
		StatusID: "status-2",
	}

	_, err := updateUC.Execute(dto)

	if err == nil {
		t.Error("Execute() expected error when gateway fails, got nil")
	}
}

type hybridOrderDataSource struct {
	findByIDFunc func(id string) (daos.OrderDAO, error)
	updateFunc   func(order daos.OrderDAO) error
}

func (ds *hybridOrderDataSource) Create(order daos.OrderDAO) error {
	return nil
}

func (ds *hybridOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return []daos.OrderDAO{}, nil
}

func (ds *hybridOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	if ds.findByIDFunc != nil {
		return ds.findByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (ds *hybridOrderDataSource) Update(order daos.OrderDAO) error {
	if ds.updateFunc != nil {
		return ds.updateFunc(order)
	}
	return nil
}

func (ds *hybridOrderDataSource) Delete(id string) error {
	return nil
}

func TestDeleteOrderUseCase_Execute_InvalidID(t *testing.T) {
	orderDS := newTestOrderDataSource()
	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewDeleteOrderUseCase(*orderGateway)

	err := uc.Execute("invalid-uuid")

	if err == nil {
		t.Error("Execute() expected error for invalid UUID, got nil")
	}
}

func TestDeleteOrderUseCase_Execute_NotFound(t *testing.T) {
	orderDS := newTestOrderDataSource()
	orderGateway := gateways.NewOrderGateway(orderDS)

	uc := NewDeleteOrderUseCase(*orderGateway)

	err := uc.Execute("550e8400-e29b-41d4-a716-446655440000")

	if err == nil {
		t.Error("Execute() expected error when order not found, got nil")
	}
}

func TestDeleteOrderUseCase_Execute_GatewayError(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	createUC := NewCreateOrderUseCase(*orderGateway, *statusGateway)
	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{
		{ProductID: "product-1", Quantity: 1, Price: 10.0},
	}
	createdOrder, _ := createUC.Execute(&customerID, items)

	hybridDS := &hybridOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return orderDS.FindByID(id)
		},
	}
	errorDeleteDS := &errorDeleteOrderDataSource{
		findByIDFunc: func(id string) (daos.OrderDAO, error) {
			return orderDS.FindByID(id)
		},
	}
	errorOrderGateway := gateways.NewOrderGateway(errorDeleteDS)

	deleteUC := NewDeleteOrderUseCase(*errorOrderGateway)

	err := deleteUC.Execute(createdOrder.ID)

	if err == nil {
		t.Error("Execute() expected error when gateway fails, got nil")
	}

	_ = hybridDS
}

type errorDeleteOrderDataSource struct {
	findByIDFunc func(id string) (daos.OrderDAO, error)
}

func (ds *errorDeleteOrderDataSource) Create(order daos.OrderDAO) error {
	return nil
}

func (ds *errorDeleteOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	return []daos.OrderDAO{}, nil
}

func (ds *errorDeleteOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	if ds.findByIDFunc != nil {
		return ds.findByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (ds *errorDeleteOrderDataSource) Update(order daos.OrderDAO) error {
	return nil
}

func (ds *errorDeleteOrderDataSource) Delete(id string) error {
	return errors.New("delete error")
}


func TestCreateOrderUseCase_Execute_EmptyItems(t *testing.T) {
	orderDS := newTestOrderDataSource()
	statusDS := newTestOrderStatusDataSource()

	orderGateway := gateways.NewOrderGateway(orderDS)
	statusGateway := gateways.NewOrderStatusGateway(statusDS)

	uc := NewCreateOrderUseCase(*orderGateway, *statusGateway)

	customerID := "customer-123"
	items := []dtos.CreateOrderItemDTO{} // Empty items - CalcTotalAmount will fail

	_, err := uc.Execute(&customerID, items)

	if err == nil {
		t.Error("Execute() expected error for empty items (CalcTotalAmount fails), got nil")
	}
}
