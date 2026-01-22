package mocks

import (
	"microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIOrderDataSource is a mock of IOrderDataSource interface.
type MockIOrderDataSource struct {
	ctrl     *gomock.Controller
	recorder *MockIOrderDataSourceMockRecorder
}

// MockIOrderDataSourceMockRecorder is the mock recorder for MockIOrderDataSource.
type MockIOrderDataSourceMockRecorder struct {
	mock *MockIOrderDataSource
}

// NewMockIOrderDataSource creates a new mock instance.
func NewMockIOrderDataSource(ctrl *gomock.Controller) *MockIOrderDataSource {
	mock := &MockIOrderDataSource{ctrl: ctrl}
	mock.recorder = &MockIOrderDataSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderDataSource) EXPECT() *MockIOrderDataSourceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIOrderDataSource) Create(order daos.OrderDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIOrderDataSourceMockRecorder) Create(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIOrderDataSource)(nil).Create), order)
}

// FindAll mocks base method.
func (m *MockIOrderDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", filter)
	ret0, _ := ret[0].([]daos.OrderDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIOrderDataSourceMockRecorder) FindAll(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIOrderDataSource)(nil).FindAll), filter)
}

// FindByID mocks base method.
func (m *MockIOrderDataSource) FindByID(id string) (daos.OrderDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(daos.OrderDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockIOrderDataSourceMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockIOrderDataSource)(nil).FindByID), id)
}

// Update mocks base method.
func (m *MockIOrderDataSource) Update(order daos.OrderDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIOrderDataSourceMockRecorder) Update(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIOrderDataSource)(nil).Update), order)
}

// Delete mocks base method.
func (m *MockIOrderDataSource) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIOrderDataSourceMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIOrderDataSource)(nil).Delete), id)
}
