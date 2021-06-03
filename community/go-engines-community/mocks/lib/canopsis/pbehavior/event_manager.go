// Code generated by MockGen. DO NOT EDIT.
// Source: lib/canopsis/pbehavior/event_manager.go

// Package mock_pbehavior is a generated GoMock package.
package mock_pbehavior

import (
	pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockEventManager is a mock of EventManager interface
type MockEventManager struct {
	ctrl     *gomock.Controller
	recorder *MockEventManagerMockRecorder
}

// MockEventManagerMockRecorder is the mock recorder for MockEventManager
type MockEventManagerMockRecorder struct {
	mock *MockEventManager
}

// NewMockEventManager creates a new mock instance
func NewMockEventManager(ctrl *gomock.Controller) *MockEventManager {
	mock := &MockEventManager{ctrl: ctrl}
	mock.recorder = &MockEventManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventManager) EXPECT() *MockEventManagerMockRecorder {
	return m.recorder
}

// GetEvent mocks base method
func (m *MockEventManager) GetEvent(arg0 pbehavior.ResolveResult, arg1 types.Alarm, arg2 time.Time) types.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.Event)
	return ret0
}

// GetEvent indicates an expected call of GetEvent
func (mr *MockEventManagerMockRecorder) GetEvent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventManager)(nil).GetEvent), arg0, arg1, arg2)
}