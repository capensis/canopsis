// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice (interfaces: Adapter)

// Package mock_entityservice is a generated GoMock package.
package mock_entityservice

import (
	context "context"
	reflect "reflect"

	entityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	gomock "github.com/golang/mock/gomock"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// MockAdapter is a mock of Adapter interface.
type MockAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockAdapterMockRecorder
}

// MockAdapterMockRecorder is the mock recorder for MockAdapter.
type MockAdapterMockRecorder struct {
	mock *MockAdapter
}

// NewMockAdapter creates a new mock instance.
func NewMockAdapter(ctrl *gomock.Controller) *MockAdapter {
	mock := &MockAdapter{ctrl: ctrl}
	mock.recorder = &MockAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdapter) EXPECT() *MockAdapterMockRecorder {
	return m.recorder
}

// GetValid mocks base method.
func (m *MockAdapter) GetValid(arg0 context.Context) ([]entityservice.EntityService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValid", arg0)
	ret0, _ := ret[0].([]entityservice.EntityService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValid indicates an expected call of GetValid.
func (mr *MockAdapterMockRecorder) GetValid(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValid", reflect.TypeOf((*MockAdapter)(nil).GetValid), arg0)
}

// UpdateBulk mocks base method.
func (m *MockAdapter) UpdateBulk(arg0 context.Context, arg1 []mongo.WriteModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBulk", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBulk indicates an expected call of UpdateBulk.
func (mr *MockAdapterMockRecorder) UpdateBulk(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBulk", reflect.TypeOf((*MockAdapter)(nil).UpdateBulk), arg0, arg1)
}
