// Code generated by MockGen. DO NOT EDIT.
// Source: lib/canopsis/alarm/activation_service.go

// Package mock_alarm is a generated GoMock package.
package mock_alarm

import (
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockActivationService is a mock of ActivationService interface
type MockActivationService struct {
	ctrl     *gomock.Controller
	recorder *MockActivationServiceMockRecorder
}

// MockActivationServiceMockRecorder is the mock recorder for MockActivationService
type MockActivationServiceMockRecorder struct {
	mock *MockActivationService
}

// NewMockActivationService creates a new mock instance
func NewMockActivationService(ctrl *gomock.Controller) *MockActivationService {
	mock := &MockActivationService{ctrl: ctrl}
	mock.recorder = &MockActivationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockActivationService) EXPECT() *MockActivationServiceMockRecorder {
	return m.recorder
}

// Process mocks base method
func (m *MockActivationService) Process(arg0 *types.Alarm) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process
func (mr *MockActivationServiceMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockActivationService)(nil).Process), arg0)
}