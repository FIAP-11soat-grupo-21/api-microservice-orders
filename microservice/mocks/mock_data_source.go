package mocks

import (
	daos "microservice/internal/adapters/daos"
	"microservice/internal/adapters/dtos"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIDataSource is a mock of IDataSource interface.
type MockIDataSource struct {
	ctrl     *gomock.Controller
	recorder *MockIDataSourceMockRecorder
}

// MockIDataSourceMockRecorder is the mock recorder for MockIDataSource.
type MockIDataSourceMockRecorder struct {
	mock *MockIDataSource
}

// NewMockIDataSource creates a new mock instance.
func NewMockIDataSource(ctrl *gomock.Controller) *MockIDataSource {
	mock := &MockIDataSource{ctrl: ctrl}
	mock.recorder = &MockIDataSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDataSource) EXPECT() *MockIDataSourceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIDataSource) Create(order daos.OrderDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIDataSourceMockRecorder) Create(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIDataSource)(nil).Create), order)
}

// FindAll mocks base method.
func (m *MockIDataSource) FindAll(filter dtos.OrderFilterDTO) ([]daos.OrderDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", filter)
	ret0, _ := ret[0].([]daos.OrderDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIDataSourceMockRecorder) FindAll(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIDataSource)(nil).FindAll), filter)
}

// FindByID mocks base method.
func (m *MockIDataSource) FindByID(id string) (daos.OrderDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(daos.OrderDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockIDataSourceMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockIDataSource)(nil).FindByID), id)
}

// Update mocks base method.
func (m *MockIDataSource) Update(order daos.OrderDAO) (daos.OrderDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", order)
	ret0, _ := ret[0].(daos.OrderDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIDataSourceMockRecorder) Update(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIDataSource)(nil).Update), order)
}

// Delete mocks base method.
func (m *MockIDataSource) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIDataSourceMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIDataSource)(nil).Delete), id)
}
