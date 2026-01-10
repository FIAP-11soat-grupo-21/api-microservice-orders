package test_helpers

import (
	"time"

	"microservice/internal/adapters/daos"
	"microservice/internal/domain/entities"
)

func CreateTestOrderDAO() daos.OrderDAO {
	customerID := "customer-123"
	now := time.Now()
	return daos.OrderDAO{
		ID:         "order-1",
		CustomerID: &customerID,
		Amount:     20.0,
		Status:     daos.OrderStatusDAO{ID: "status-1", Name: "Pending"},
		Items: []daos.OrderItemDAO{
			{ID: "item-1", OrderID: "order-1", ProductID: "product-1", Quantity: 2, UnitPrice: 10.0},
		},
		CreatedAt: now,
	}
}

func CreateTestOrderEntity() entities.Order {
	customerID := "customer-123"
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	item, _ := entities.NewOrderItem("item-1", "product-1", "order-1", 2, 10.0)
	now := time.Now()
	order, _ := entities.NewOrderWithItems("order-1", &customerID, 20.0, *status, []entities.OrderItem{*item}, now, nil)
	return *order
}

func CreateTestOrderStatusDAO() daos.OrderStatusDAO {
	return daos.OrderStatusDAO{ID: "status-1", Name: "Pending"}
}

func CreateTestOrderStatusEntity() entities.OrderStatus {
	status, _ := entities.NewOrderStatus("status-1", "Pending")
	return *status
}
