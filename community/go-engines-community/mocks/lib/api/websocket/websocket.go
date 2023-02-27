// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket (interfaces: Upgrader,Connection,Authorizer,Hub)

// Package mock_websocket is a generated GoMock package.
package mock_websocket

import (
	context "context"
	net "net"
	http "net/http"
	reflect "reflect"
	time "time"

	websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	gomock "github.com/golang/mock/gomock"
)

// MockUpgrader is a mock of Upgrader interface.
type MockUpgrader struct {
	ctrl     *gomock.Controller
	recorder *MockUpgraderMockRecorder
}

// MockUpgraderMockRecorder is the mock recorder for MockUpgrader.
type MockUpgraderMockRecorder struct {
	mock *MockUpgrader
}

// NewMockUpgrader creates a new mock instance.
func NewMockUpgrader(ctrl *gomock.Controller) *MockUpgrader {
	mock := &MockUpgrader{ctrl: ctrl}
	mock.recorder = &MockUpgraderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpgrader) EXPECT() *MockUpgraderMockRecorder {
	return m.recorder
}

// Upgrade mocks base method.
func (m *MockUpgrader) Upgrade(arg0 http.ResponseWriter, arg1 *http.Request, arg2 http.Header) (websocket.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upgrade", arg0, arg1, arg2)
	ret0, _ := ret[0].(websocket.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upgrade indicates an expected call of Upgrade.
func (mr *MockUpgraderMockRecorder) Upgrade(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockUpgrader)(nil).Upgrade), arg0, arg1, arg2)
}

// MockConnection is a mock of Connection interface.
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection.
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance.
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnection)(nil).Close))
}

// ReadJSON mocks base method.
func (m *MockConnection) ReadJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadJSON indicates an expected call of ReadJSON.
func (mr *MockConnectionMockRecorder) ReadJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadJSON", reflect.TypeOf((*MockConnection)(nil).ReadJSON), arg0)
}

// RemoteAddr mocks base method.
func (m *MockConnection) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockConnectionMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConnection)(nil).RemoteAddr))
}

// SetPongHandler mocks base method.
func (m *MockConnection) SetPongHandler(arg0 func(string) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPongHandler", arg0)
}

// SetPongHandler indicates an expected call of SetPongHandler.
func (mr *MockConnectionMockRecorder) SetPongHandler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPongHandler", reflect.TypeOf((*MockConnection)(nil).SetPongHandler), arg0)
}

// SetReadDeadline mocks base method.
func (m *MockConnection) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockConnectionMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockConnection)(nil).SetReadDeadline), arg0)
}

// WriteControl mocks base method.
func (m *MockConnection) WriteControl(arg0 int, arg1 []byte, arg2 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteControl", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteControl indicates an expected call of WriteControl.
func (mr *MockConnectionMockRecorder) WriteControl(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteControl", reflect.TypeOf((*MockConnection)(nil).WriteControl), arg0, arg1, arg2)
}

// WriteJSON mocks base method.
func (m *MockConnection) WriteJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteJSON indicates an expected call of WriteJSON.
func (mr *MockConnectionMockRecorder) WriteJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteJSON", reflect.TypeOf((*MockConnection)(nil).WriteJSON), arg0)
}

// MockAuthorizer is a mock of Authorizer interface.
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer.
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance.
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// AddGroup mocks base method.
func (m *MockAuthorizer) AddGroup(arg0 string, arg1 []string, arg2 websocket.GroupCheck) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddGroup indicates an expected call of AddGroup.
func (mr *MockAuthorizerMockRecorder) AddGroup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGroup", reflect.TypeOf((*MockAuthorizer)(nil).AddGroup), arg0, arg1, arg2)
}

// AddRoom mocks base method.
func (m *MockAuthorizer) AddRoom(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRoom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRoom indicates an expected call of AddRoom.
func (mr *MockAuthorizerMockRecorder) AddRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRoom", reflect.TypeOf((*MockAuthorizer)(nil).AddRoom), arg0, arg1)
}

// Authenticate mocks base method.
func (m *MockAuthorizer) Authenticate(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockAuthorizerMockRecorder) Authenticate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuthorizer)(nil).Authenticate), arg0, arg1)
}

// Authorize mocks base method.
func (m *MockAuthorizer) Authorize(arg0 context.Context, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockAuthorizerMockRecorder) Authorize(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockAuthorizer)(nil).Authorize), arg0, arg1, arg2)
}

// GetGroupIds mocks base method.
func (m *MockAuthorizer) GetGroupIds(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupIds", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetGroupIds indicates an expected call of GetGroupIds.
func (mr *MockAuthorizerMockRecorder) GetGroupIds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupIds", reflect.TypeOf((*MockAuthorizer)(nil).GetGroupIds), arg0)
}

// RemoveGroupRoom mocks base method.
func (m *MockAuthorizer) RemoveGroupRoom(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveGroupRoom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveGroupRoom indicates an expected call of RemoveGroupRoom.
func (mr *MockAuthorizerMockRecorder) RemoveGroupRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveGroupRoom", reflect.TypeOf((*MockAuthorizer)(nil).RemoveGroupRoom), arg0, arg1)
}

// MockHub is a mock of Hub interface.
type MockHub struct {
	ctrl     *gomock.Controller
	recorder *MockHubMockRecorder
}

// MockHubMockRecorder is the mock recorder for MockHub.
type MockHubMockRecorder struct {
	mock *MockHub
}

// NewMockHub creates a new mock instance.
func NewMockHub(ctrl *gomock.Controller) *MockHub {
	mock := &MockHub{ctrl: ctrl}
	mock.recorder = &MockHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHub) EXPECT() *MockHubMockRecorder {
	return m.recorder
}

// CloseGroupRoom mocks base method.
func (m *MockHub) CloseGroupRoom(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseGroupRoom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseGroupRoom indicates an expected call of CloseGroupRoom.
func (mr *MockHubMockRecorder) CloseGroupRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseGroupRoom", reflect.TypeOf((*MockHub)(nil).CloseGroupRoom), arg0, arg1)
}

// CloseGroupRoomAndNotify mocks base method.
func (m *MockHub) CloseGroupRoomAndNotify(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseGroupRoomAndNotify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseGroupRoomAndNotify indicates an expected call of CloseGroupRoomAndNotify.
func (mr *MockHubMockRecorder) CloseGroupRoomAndNotify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseGroupRoomAndNotify", reflect.TypeOf((*MockHub)(nil).CloseGroupRoomAndNotify), arg0, arg1)
}

// Connect mocks base method.
func (m *MockHub) Connect(arg0 http.ResponseWriter, arg1 *http.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockHubMockRecorder) Connect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockHub)(nil).Connect), arg0, arg1)
}

// GetConnectedGroupIds mocks base method.
func (m *MockHub) GetConnectedGroupIds(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectedGroupIds", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetConnectedGroupIds indicates an expected call of GetConnectedGroupIds.
func (mr *MockHubMockRecorder) GetConnectedGroupIds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectedGroupIds", reflect.TypeOf((*MockHub)(nil).GetConnectedGroupIds), arg0)
}

// GetConnections mocks base method.
func (m *MockHub) GetConnections() []websocket.UserConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnections")
	ret0, _ := ret[0].([]websocket.UserConnection)
	return ret0
}

// GetConnections indicates an expected call of GetConnections.
func (mr *MockHubMockRecorder) GetConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnections", reflect.TypeOf((*MockHub)(nil).GetConnections))
}

// GetGroupIds mocks base method.
func (m *MockHub) GetGroupIds(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupIds", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetGroupIds indicates an expected call of GetGroupIds.
func (mr *MockHubMockRecorder) GetGroupIds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupIds", reflect.TypeOf((*MockHub)(nil).GetGroupIds), arg0)
}

// GetUsers mocks base method.
func (m *MockHub) GetUsers() map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers")
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockHubMockRecorder) GetUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockHub)(nil).GetUsers))
}

// RegisterGroup mocks base method.
func (m *MockHub) RegisterGroup(arg0 string, arg1 websocket.GroupCheck, arg2 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterGroup", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterGroup indicates an expected call of RegisterGroup.
func (mr *MockHubMockRecorder) RegisterGroup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterGroup", reflect.TypeOf((*MockHub)(nil).RegisterGroup), varargs...)
}

// RegisterRoom mocks base method.
func (m *MockHub) RegisterRoom(arg0 string, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterRoom", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterRoom indicates an expected call of RegisterRoom.
func (mr *MockHubMockRecorder) RegisterRoom(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRoom", reflect.TypeOf((*MockHub)(nil).RegisterRoom), varargs...)
}

// Send mocks base method.
func (m *MockHub) Send(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", arg0, arg1)
}

// Send indicates an expected call of Send.
func (mr *MockHubMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockHub)(nil).Send), arg0, arg1)
}

// SendGroupRoom mocks base method.
func (m *MockHub) SendGroupRoom(arg0, arg1 string, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendGroupRoom", arg0, arg1, arg2)
}

// SendGroupRoom indicates an expected call of SendGroupRoom.
func (mr *MockHubMockRecorder) SendGroupRoom(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendGroupRoom", reflect.TypeOf((*MockHub)(nil).SendGroupRoom), arg0, arg1, arg2)
}

// Start mocks base method.
func (m *MockHub) Start(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0)
}

// Start indicates an expected call of Start.
func (mr *MockHubMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockHub)(nil).Start), arg0)
}
