package mocks

import (
	"microservice/internal/adapters/daos"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIOrderStatusDataSource is a mock of IOrderStatusDataSource interface.
type MockIOrderStatusDataSource struct {
	ctrl     *gomock.Controller
	recorder *MockIOrderStatusDataSourceMockRecorder
}

// MockIOrderStatusDataSourceMockRecorder is the mock recorder for MockIOrderStatusDataSource.
type MockIOrderStatusDataSourceMockRecorder struct {
	mock *MockIOrderStatusDataSource
}

// NewMockIOrderStatusDataSource creates a new mock instance.
func NewMockIOrderStatusDataSource(ctrl *gomock.Controller) *MockIOrderStatusDataSource {
	mock := &MockIOrderStatusDataSource{ctrl: ctrl}
	mock.recorder = &MockIOrderStatusDataSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderStatusDataSource) EXPECT() *MockIOrderStatusDataSourceMockRecorder {
	return m.recorder
}

// FindByID mocks base method.
func (m *MockIOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(daos.OrderStatusDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockIOrderStatusDataSourceMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockIOrderStatusDataSource)(nil).FindByID), id)
}

// FindByName mocks base method.
func (m *MockIOrderStatusDataSource) FindByName(name string) (daos.OrderStatusDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", name)
	ret0, _ := ret[0].(daos.OrderStatusDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName.
func (mr *MockIOrderStatusDataSourceMockRecorder) FindByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockIOrderStatusDataSource)(nil).FindByName), name)
}

// FindAll mocks base method.
func (m *MockIOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll")
	ret0, _ := ret[0].([]daos.OrderStatusDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIOrderStatusDataSourceMockRecorder) FindAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIOrderStatusDataSource)(nil).FindAll))
}
