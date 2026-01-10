package test_helpers

import (
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type MockOrderDataSource struct {
	CreateFunc   func(order daos.OrderDAO) error
	FindAllFunc  func(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error)
	FindByIDFunc func(id string) (daos.OrderDAO, error)
	UpdateFunc   func(order daos.OrderDAO) error
	DeleteFunc   func(id string) error
}

func (m *MockOrderDataSource) Create(order daos.OrderDAO) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(order)
	}
	return nil
}

func (m *MockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(filter)
	}
	return []daos.OrderDAO{}, nil
}

func (m *MockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return daos.OrderDAO{}, nil
}

func (m *MockOrderDataSource) Update(order daos.OrderDAO) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(order)
	}
	return nil
}

func (m *MockOrderDataSource) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

type MockOrderStatusDataSource struct {
	FindByIDFunc func(id string) (daos.OrderStatusDAO, error)
	FindAllFunc  func() ([]daos.OrderStatusDAO, error)
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return daos.OrderStatusDAO{}, nil
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}
