package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockApiClient is a mock of IApiClient interface.
type MockApiClient struct {
	ctrl     *gomock.Controller
	recorder *MockApiClientMockRecorder
}

// MockApiClientMockRecorder is the mock recorder for MockApiClient.
type MockApiClientMockRecorder struct {
	mock *MockApiClient
}

// NewMockApiClient creates a new mock instance.
func NewMockApiClient(ctrl *gomock.Controller) *MockApiClient {
	mock := &MockApiClient{ctrl: ctrl}
	mock.recorder = &MockApiClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiClient) EXPECT() *MockApiClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockApiClient) Get(path string, obj interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", path, obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockApiClientMockRecorder) Get(path, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockApiClient)(nil).Get), path, obj)
}
