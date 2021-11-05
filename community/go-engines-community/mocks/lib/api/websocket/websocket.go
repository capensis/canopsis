// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket (interfaces: Upgrader,Connection,Authorizer)

// Package mock_websocket is a generated GoMock package.
package mock_websocket

import (
	websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	gomock "github.com/golang/mock/gomock"
	net "net"
	http "net/http"
	reflect "reflect"
	time "time"
)

// MockUpgrader is a mock of Upgrader interface
type MockUpgrader struct {
	ctrl     *gomock.Controller
	recorder *MockUpgraderMockRecorder
}

// MockUpgraderMockRecorder is the mock recorder for MockUpgrader
type MockUpgraderMockRecorder struct {
	mock *MockUpgrader
}

// NewMockUpgrader creates a new mock instance
func NewMockUpgrader(ctrl *gomock.Controller) *MockUpgrader {
	mock := &MockUpgrader{ctrl: ctrl}
	mock.recorder = &MockUpgraderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpgrader) EXPECT() *MockUpgraderMockRecorder {
	return m.recorder
}

// Upgrade mocks base method
func (m *MockUpgrader) Upgrade(arg0 http.ResponseWriter, arg1 *http.Request, arg2 http.Header) (websocket.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upgrade", arg0, arg1, arg2)
	ret0, _ := ret[0].(websocket.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upgrade indicates an expected call of Upgrade
func (mr *MockUpgraderMockRecorder) Upgrade(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockUpgrader)(nil).Upgrade), arg0, arg1, arg2)
}

// MockConnection is a mock of Connection interface
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnection)(nil).Close))
}

// ReadJSON mocks base method
func (m *MockConnection) ReadJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadJSON indicates an expected call of ReadJSON
func (mr *MockConnectionMockRecorder) ReadJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadJSON", reflect.TypeOf((*MockConnection)(nil).ReadJSON), arg0)
}

// RemoteAddr mocks base method
func (m *MockConnection) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr
func (mr *MockConnectionMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConnection)(nil).RemoteAddr))
}

// SetPongHandler mocks base method
func (m *MockConnection) SetPongHandler(arg0 func(string) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPongHandler", arg0)
}

// SetPongHandler indicates an expected call of SetPongHandler
func (mr *MockConnectionMockRecorder) SetPongHandler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPongHandler", reflect.TypeOf((*MockConnection)(nil).SetPongHandler), arg0)
}

// SetReadDeadline mocks base method
func (m *MockConnection) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline
func (mr *MockConnectionMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockConnection)(nil).SetReadDeadline), arg0)
}

// WriteControl mocks base method
func (m *MockConnection) WriteControl(arg0 int, arg1 []byte, arg2 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteControl", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteControl indicates an expected call of WriteControl
func (mr *MockConnectionMockRecorder) WriteControl(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteControl", reflect.TypeOf((*MockConnection)(nil).WriteControl), arg0, arg1, arg2)
}

// WriteJSON mocks base method
func (m *MockConnection) WriteJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteJSON indicates an expected call of WriteJSON
func (mr *MockConnectionMockRecorder) WriteJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteJSON", reflect.TypeOf((*MockConnection)(nil).WriteJSON), arg0)
}

// MockAuthorizer is a mock of Authorizer interface
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// Auth mocks base method
func (m *MockAuthorizer) Auth(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth
func (mr *MockAuthorizerMockRecorder) Auth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockAuthorizer)(nil).Auth), arg0, arg1)
}

// Exists mocks base method
func (m *MockAuthorizer) Exists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockAuthorizerMockRecorder) Exists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockAuthorizer)(nil).Exists), arg0)
}