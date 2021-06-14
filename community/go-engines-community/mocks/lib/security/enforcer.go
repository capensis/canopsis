// Code generated by MockGen. DO NOT EDIT.
// Source: lib/security/enforcer.go

// Package mock_security is a generated GoMock package.
package mock_security

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEnforcer is a mock of Enforcer interface
type MockEnforcer struct {
	ctrl     *gomock.Controller
	recorder *MockEnforcerMockRecorder
}

// MockEnforcerMockRecorder is the mock recorder for MockEnforcer
type MockEnforcerMockRecorder struct {
	mock *MockEnforcer
}

// NewMockEnforcer creates a new mock instance
func NewMockEnforcer(ctrl *gomock.Controller) *MockEnforcer {
	mock := &MockEnforcer{ctrl: ctrl}
	mock.recorder = &MockEnforcerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEnforcer) EXPECT() *MockEnforcerMockRecorder {
	return m.recorder
}

// Enforce mocks base method
func (m *MockEnforcer) Enforce(rvals ...interface{}) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range rvals {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Enforce", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enforce indicates an expected call of Enforce
func (mr *MockEnforcerMockRecorder) Enforce(rvals ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enforce", reflect.TypeOf((*MockEnforcer)(nil).Enforce), rvals...)
}

// StartAutoLoadPolicy mocks base method
func (m *MockEnforcer) StartAutoLoadPolicy(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartAutoLoadPolicy", arg0)
}

// StartAutoLoadPolicy indicates an expected call of StartAutoLoadPolicy
func (mr *MockEnforcerMockRecorder) StartAutoLoadPolicy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartAutoLoadPolicy", reflect.TypeOf((*MockEnforcer)(nil).StartAutoLoadPolicy), arg0)
}

// LoadPolicy mocks base method
func (m *MockEnforcer) LoadPolicy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadPolicy")
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadPolicy indicates an expected call of LoadPolicy
func (mr *MockEnforcerMockRecorder) LoadPolicy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadPolicy", reflect.TypeOf((*MockEnforcer)(nil).LoadPolicy))
}
