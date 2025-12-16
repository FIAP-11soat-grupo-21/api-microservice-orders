package presenters

import (
	"testing"
	"time"

	"microservice/internal/domain/entities"
)

func TestToOrderResponse(t *testing.T) {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	now := time.Now()
	order, _ := entities.NewOrderWithItems(
		"order-1",
		&customerID,
		20.0,
		*status,
		[]entities.OrderItem{*item},
		now,
		nil,
	)

	response := ToOrderResponse(*order)

	if response.ID != "order-1" {
		t.Errorf("ToOrderResponse() ID = %v, want order-1", response.ID)
	}
	if *response.CustomerID != customerID {
		t.Errorf("ToOrderResponse() CustomerID = %v, want %v", *response.CustomerID, customerID)
	}
	if response.Amount != 20.0 {
		t.Errorf("ToOrderResponse() Amount = %v, want 20.0", response.Amount)
	}
	if response.Status.ID != "status-1" {
		t.Errorf("ToOrderResponse() Status.ID = %v, want status-1", response.Status.ID)
	}
	if response.Status.Name != "Pending" {
		t.Errorf("ToOrderResponse() Status.Name = %v, want Pending", response.Status.Name)
	}
	if len(response.Items) != 1 {
		t.Errorf("ToOrderResponse() Items length = %v, want 1", len(response.Items))
	}
	if response.Items[0].ID != "item-1" {
		t.Errorf("ToOrderResponse() Items[0].ID = %v, want item-1", response.Items[0].ID)
	}
	if response.Items[0].ProductID != "product-1" {
		t.Errorf("ToOrderResponse() Items[0].ProductID = %v, want product-1", response.Items[0].ProductID)
	}
	if response.Items[0].Quantity != 2 {
		t.Errorf("ToOrderResponse() Items[0].Quantity = %v, want 2", response.Items[0].Quantity)
	}
	if response.Items[0].UnitPrice != 10.0 {
		t.Errorf("ToOrderResponse() Items[0].UnitPrice = %v, want 10.0", response.Items[0].UnitPrice)
	}
}

func TestToOrderResponse_WithUpdatedAt(t *testing.T) {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 1, 10.0)
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	order, _ := entities.NewOrderWithItems(
		"order-1",
		&customerID,
		10.0,
		*status,
		[]entities.OrderItem{*item},
		now,
		&updatedAt,
	)

	response := ToOrderResponse(*order)

	if response.UpdatedAt == nil {
		t.Error("ToOrderResponse() UpdatedAt should not be nil")
	}
}

func TestToOrderResponse_NilCustomerID(t *testing.T) {
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 1, 10.0)
	now := time.Now()
	order, _ := entities.NewOrderWithItems(
		"order-1",
		nil,
		10.0,
		*status,
		[]entities.OrderItem{*item},
		now,
		nil,
	)

	response := ToOrderResponse(*order)

	if response.CustomerID != nil {
		t.Errorf("ToOrderResponse() CustomerID = %v, want nil", response.CustomerID)
	}
}

func TestToOrderResponseList(t *testing.T) {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item1, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 1, 10.0)
	item2, _ := entities.NewOrderItem("item-2", "product-2", "order-2", 2, 20.0)
	now := time.Now()

	order1, _ := entities.NewOrderWithItems("order-1", &customerID, 10.0, *status, []entities.OrderItem{*item1}, now, nil)
	order2, _ := entities.NewOrderWithItems("order-2", &customerID, 40.0, *status, []entities.OrderItem{*item2}, now, nil)

	orders := []entities.Order{*order1, *order2}
	responses := ToOrderResponseList(orders)

	if len(responses) != 2 {
		t.Errorf("ToOrderResponseList() length = %v, want 2", len(responses))
	}
	if responses[0].ID != "order-1" {
		t.Errorf("ToOrderResponseList()[0].ID = %v, want order-1", responses[0].ID)
	}
	if responses[1].ID != "order-2" {
		t.Errorf("ToOrderResponseList()[1].ID = %v, want order-2", responses[1].ID)
	}
}

func TestToOrderResponseList_Empty(t *testing.T) {
	orders := []entities.Order{}
	responses := ToOrderResponseList(orders)

	if len(responses) != 0 {
		t.Errorf("ToOrderResponseList() length = %v, want 0", len(responses))
	}
}

func TestToOrderStatusResponse(t *testing.T) {
	status, _ := entities.NewOrderStatus("status-1", "Pending")

	response := ToOrderStatusResponse(*status)

	if response.ID != "status-1" {
		t.Errorf("ToOrderStatusResponse() ID = %v, want status-1", response.ID)
	}
	if response.Name != "Pending" {
		t.Errorf("ToOrderStatusResponse() Name = %v, want Pending", response.Name)
	}
}

func TestToOrderStatusResponseList(t *testing.T) {
	status1, _ := entities.NewOrderStatus("status-1", "Pending")
	status2, _ := entities.NewOrderStatus("status-2", "Confirmed")
	status3, _ := entities.NewOrderStatus("status-3", "Completed")

	statuses := []entities.OrderStatus{*status1, *status2, *status3}
	responses := ToOrderStatusResponseList(statuses)

	if len(responses) != 3 {
		t.Errorf("ToOrderStatusResponseList() length = %v, want 3", len(responses))
	}
	if responses[0].ID != "status-1" {
		t.Errorf("ToOrderStatusResponseList()[0].ID = %v, want status-1", responses[0].ID)
	}
	if responses[1].Name != "Confirmed" {
		t.Errorf("ToOrderStatusResponseList()[1].Name = %v, want Confirmed", responses[1].Name)
	}
}

func TestToOrderStatusResponseList_Empty(t *testing.T) {
	statuses := []entities.OrderStatus{}
	responses := ToOrderStatusResponseList(statuses)

	if len(responses) != 0 {
		t.Errorf("ToOrderStatusResponseList() length = %v, want 0", len(responses))
	}
}
