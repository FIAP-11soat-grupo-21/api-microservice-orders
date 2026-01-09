package use_cases

import (
	"microservice/internal/adapters/dtos"
	"microservice/internal/domain/entities"
	"microservice/internal/domain/exceptions"
)

type MockOrderGateway struct {
	orders             map[string]*entities.Order
	shouldFailFindByID bool
	shouldFailUpdate   bool
	shouldFailCreate   bool
}

func NewMockOrderGateway() *MockOrderGateway {
	return &MockOrderGateway{
		orders: make(map[string]*entities.Order),
	}
}

func (m *MockOrderGateway) AddOrder(order *entities.Order) {
	m.orders[order.ID] = order
}

func (m *MockOrderGateway) SetShouldFailFindByID(fail bool) {
	m.shouldFailFindByID = fail
}

func (m *MockOrderGateway) SetShouldFailUpdate(fail bool) {
	m.shouldFailUpdate = fail
}

func (m *MockOrderGateway) SetShouldFailCreate(fail bool) {
	m.shouldFailCreate = fail
}

func (m *MockOrderGateway) Create(order entities.Order) error {
	if m.shouldFailCreate {
		return &exceptions.InvalidOrderDataException{Message: "Create failed"}
	}
	m.orders[order.ID] = &order
	return nil
}

func (m *MockOrderGateway) FindByID(id string) (*entities.Order, error) {
	if m.shouldFailFindByID {
		return nil, &exceptions.OrderNotFoundException{}
	}

	order, exists := m.orders[id]
	if !exists {
		return nil, &exceptions.OrderNotFoundException{}
	}
	return order, nil
}

func (m *MockOrderGateway) FindAll(filter dtos.OrderFilterDTO) ([]entities.Order, error) {
	orders := make([]entities.Order, 0, len(m.orders))
	for _, order := range m.orders {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (m *MockOrderGateway) Update(order entities.Order) error {
	if m.shouldFailUpdate {
		return &exceptions.InvalidOrderDataException{Message: "Update failed"}
	}

	m.orders[order.ID] = &order
	return nil
}

func (m *MockOrderGateway) Delete(id string) error {
	delete(m.orders, id)
	return nil
}

type MockOrderStatusGateway struct {
	statuses           map[string]*entities.OrderStatus
	shouldFailFindByID bool
}

func NewMockOrderStatusGateway() *MockOrderStatusGateway {
	return &MockOrderStatusGateway{
		statuses: make(map[string]*entities.OrderStatus),
	}
}

func (m *MockOrderStatusGateway) AddStatus(status *entities.OrderStatus) {
	m.statuses[status.ID] = status
}

func (m *MockOrderStatusGateway) SetShouldFailFindByID(fail bool) {
	m.shouldFailFindByID = fail
}

func (m *MockOrderStatusGateway) FindAll() ([]entities.OrderStatus, error) {
	statuses := make([]entities.OrderStatus, 0, len(m.statuses))
	for _, status := range m.statuses {
		statuses = append(statuses, *status)
	}
	return statuses, nil
}

func (m *MockOrderStatusGateway) FindByID(id string) (*entities.OrderStatus, error) {
	if m.shouldFailFindByID {
		return nil, &exceptions.OrderStatusNotFoundException{}
	}

	status, exists := m.statuses[id]
	if !exists {
		return nil, &exceptions.OrderStatusNotFoundException{}
	}
	return status, nil
}
