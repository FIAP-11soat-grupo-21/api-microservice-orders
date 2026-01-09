package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
)

type MockOrderDataSource struct {
	mock.Mock
}

func (m *MockOrderDataSource) Create(order daos.OrderDAO) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	args := m.Called(filter)
	return args.Get(0).([]daos.OrderDAO), args.Error(1)
}

func (m *MockOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	args := m.Called(id)
	return args.Get(0).(daos.OrderDAO), args.Error(1)
}

func (m *MockOrderDataSource) Update(order daos.OrderDAO) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderDataSource) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockOrderStatusDataSource struct {
	mock.Mock
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	args := m.Called(id)
	return args.Get(0).(daos.OrderStatusDAO), args.Error(1)
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	args := m.Called()
	return args.Get(0).([]daos.OrderStatusDAO), args.Error(1)
}

func TestIOrderDataSource_Interface(t *testing.T) {
	var dataSource IOrderDataSource = &MockOrderDataSource{}
	assert.NotNil(t, dataSource)

	mockDS := &MockOrderDataSource{}

	order := daos.OrderDAO{ID: "test-order"}
	mockDS.On("Create", order).Return(nil)
	err := mockDS.Create(order)
	assert.NoError(t, err)
	mockDS.AssertExpectations(t)

	mockDS2 := &MockOrderDataSource{}
	filter := dtos.OrderFilterDTO{}
	expectedOrders := []daos.OrderDAO{{ID: "order1"}, {ID: "order2"}}
	mockDS2.On("FindAll", filter).Return(expectedOrders, nil)
	orders, err := mockDS2.FindAll(filter)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	mockDS2.AssertExpectations(t)

	mockDS3 := &MockOrderDataSource{}
	expectedOrder := daos.OrderDAO{ID: "test-order"}
	mockDS3.On("FindByID", "test-order").Return(expectedOrder, nil)
	foundOrder, err := mockDS3.FindByID("test-order")
	assert.NoError(t, err)
	assert.Equal(t, "test-order", foundOrder.ID)
	mockDS3.AssertExpectations(t)

	mockDS4 := &MockOrderDataSource{}
	updateOrder := daos.OrderDAO{ID: "update-order"}
	mockDS4.On("Update", updateOrder).Return(nil)
	err = mockDS4.Update(updateOrder)
	assert.NoError(t, err)
	mockDS4.AssertExpectations(t)

	mockDS5 := &MockOrderDataSource{}
	mockDS5.On("Delete", "delete-order").Return(nil)
	err = mockDS5.Delete("delete-order")
	assert.NoError(t, err)
	mockDS5.AssertExpectations(t)
}

func TestIOrderStatusDataSource_Interface(t *testing.T) {
	var statusDataSource IOrderStatusDataSource = &MockOrderStatusDataSource{}
	assert.NotNil(t, statusDataSource)

	mockSDS := &MockOrderStatusDataSource{}

	expectedStatus := daos.OrderStatusDAO{ID: "status-1", Name: "pending"}
	mockSDS.On("FindByID", "status-1").Return(expectedStatus, nil)
	status, err := mockSDS.FindByID("status-1")
	assert.NoError(t, err)
	assert.Equal(t, "status-1", status.ID)
	assert.Equal(t, "pending", status.Name)
	mockSDS.AssertExpectations(t)

	mockSDS2 := &MockOrderStatusDataSource{}
	expectedStatuses := []daos.OrderStatusDAO{
		{ID: "status-1", Name: "pending"},
		{ID: "status-2", Name: "confirmed"},
	}
	mockSDS2.On("FindAll").Return(expectedStatuses, nil)
	statuses, err := mockSDS2.FindAll()
	assert.NoError(t, err)
	assert.Len(t, statuses, 2)
	assert.Equal(t, "pending", statuses[0].Name)
	assert.Equal(t, "confirmed", statuses[1].Name)
	mockSDS2.AssertExpectations(t)
}

func TestInterfaceCompatibility(t *testing.T) {
	var orderDS IOrderDataSource
	var statusDS IOrderStatusDataSource

	orderDS = &MockOrderDataSource{}
	statusDS = &MockOrderStatusDataSource{}

	assert.NotNil(t, orderDS)
	assert.NotNil(t, statusDS)

	assert.NotPanics(t, func() {
		_ = orderDS.Create
		_ = orderDS.FindAll
		_ = orderDS.FindByID
		_ = orderDS.Update
		_ = orderDS.Delete

		_ = statusDS.FindByID
		_ = statusDS.FindAll
	})
}
