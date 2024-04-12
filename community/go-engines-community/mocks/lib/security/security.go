// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security (interfaces: Enforcer,Provider,HttpProvider,UserProvider)

// Package mock_security is a generated GoMock package.
package mock_security

import (
	context "context"
	http "net/http"
	reflect "reflect"
	time "time"

	security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	gomock "github.com/golang/mock/gomock"
)

// MockEnforcer is a mock of Enforcer interface.
type MockEnforcer struct {
	ctrl     *gomock.Controller
	recorder *MockEnforcerMockRecorder
}

// MockEnforcerMockRecorder is the mock recorder for MockEnforcer.
type MockEnforcerMockRecorder struct {
	mock *MockEnforcer
}

// NewMockEnforcer creates a new mock instance.
func NewMockEnforcer(ctrl *gomock.Controller) *MockEnforcer {
	mock := &MockEnforcer{ctrl: ctrl}
	mock.recorder = &MockEnforcerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnforcer) EXPECT() *MockEnforcerMockRecorder {
	return m.recorder
}

// Enforce mocks base method.
func (m *MockEnforcer) Enforce(arg0 ...interface{}) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Enforce", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enforce indicates an expected call of Enforce.
func (mr *MockEnforcerMockRecorder) Enforce(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enforce", reflect.TypeOf((*MockEnforcer)(nil).Enforce), arg0...)
}

// GetPermissionsForUser mocks base method.
func (m *MockEnforcer) GetPermissionsForUser(arg0 string, arg1 ...string) ([][]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPermissionsForUser", varargs...)
	ret0, _ := ret[0].([][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPermissionsForUser indicates an expected call of GetPermissionsForUser.
func (mr *MockEnforcerMockRecorder) GetPermissionsForUser(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPermissionsForUser", reflect.TypeOf((*MockEnforcer)(nil).GetPermissionsForUser), varargs...)
}

// GetRolesForUser mocks base method.
func (m *MockEnforcer) GetRolesForUser(arg0 string, arg1 ...string) ([]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRolesForUser", varargs...)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRolesForUser indicates an expected call of GetRolesForUser.
func (mr *MockEnforcerMockRecorder) GetRolesForUser(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolesForUser", reflect.TypeOf((*MockEnforcer)(nil).GetRolesForUser), varargs...)
}

// HasPermissionForUser mocks base method.
func (m *MockEnforcer) HasPermissionForUser(arg0 string, arg1 ...string) bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HasPermissionForUser", varargs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasPermissionForUser indicates an expected call of HasPermissionForUser.
func (mr *MockEnforcerMockRecorder) HasPermissionForUser(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPermissionForUser", reflect.TypeOf((*MockEnforcer)(nil).HasPermissionForUser), varargs...)
}

// LoadPolicy mocks base method.
func (m *MockEnforcer) LoadPolicy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadPolicy")
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadPolicy indicates an expected call of LoadPolicy.
func (mr *MockEnforcerMockRecorder) LoadPolicy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadPolicy", reflect.TypeOf((*MockEnforcer)(nil).LoadPolicy))
}

// StartAutoLoadPolicy mocks base method.
func (m *MockEnforcer) StartAutoLoadPolicy(arg0 context.Context, arg1 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartAutoLoadPolicy", arg0, arg1)
}

// StartAutoLoadPolicy indicates an expected call of StartAutoLoadPolicy.
func (mr *MockEnforcerMockRecorder) StartAutoLoadPolicy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartAutoLoadPolicy", reflect.TypeOf((*MockEnforcer)(nil).StartAutoLoadPolicy), arg0, arg1)
}

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockProvider) Auth(arg0 context.Context, arg1, arg2 string) (*security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0, arg1, arg2)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth.
func (mr *MockProviderMockRecorder) Auth(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockProvider)(nil).Auth), arg0, arg1, arg2)
}

// GetName mocks base method.
func (m *MockProvider) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockProviderMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockProvider)(nil).GetName))
}

// MockHttpProvider is a mock of HttpProvider interface.
type MockHttpProvider struct {
	ctrl     *gomock.Controller
	recorder *MockHttpProviderMockRecorder
}

// MockHttpProviderMockRecorder is the mock recorder for MockHttpProvider.
type MockHttpProviderMockRecorder struct {
	mock *MockHttpProvider
}

// NewMockHttpProvider creates a new mock instance.
func NewMockHttpProvider(ctrl *gomock.Controller) *MockHttpProvider {
	mock := &MockHttpProvider{ctrl: ctrl}
	mock.recorder = &MockHttpProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpProvider) EXPECT() *MockHttpProviderMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockHttpProvider) Auth(arg0 *http.Request) (*security.User, error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// Auth indicates an expected call of Auth.
func (mr *MockHttpProviderMockRecorder) Auth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockHttpProvider)(nil).Auth), arg0)
}

// MockUserProvider is a mock of UserProvider interface.
type MockUserProvider struct {
	ctrl     *gomock.Controller
	recorder *MockUserProviderMockRecorder
}

// MockUserProviderMockRecorder is the mock recorder for MockUserProvider.
type MockUserProviderMockRecorder struct {
	mock *MockUserProvider
}

// NewMockUserProvider creates a new mock instance.
func NewMockUserProvider(ctrl *gomock.Controller) *MockUserProvider {
	mock := &MockUserProvider{ctrl: ctrl}
	mock.recorder = &MockUserProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserProvider) EXPECT() *MockUserProviderMockRecorder {
	return m.recorder
}

// FindByAuthApiKey mocks base method.
func (m *MockUserProvider) FindByAuthApiKey(arg0 context.Context, arg1 string) (*security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByAuthApiKey", arg0, arg1)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByAuthApiKey indicates an expected call of FindByAuthApiKey.
func (mr *MockUserProviderMockRecorder) FindByAuthApiKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByAuthApiKey", reflect.TypeOf((*MockUserProvider)(nil).FindByAuthApiKey), arg0, arg1)
}

// FindByExternalSource mocks base method.
func (m *MockUserProvider) FindByExternalSource(arg0 context.Context, arg1, arg2 string) (*security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByExternalSource", arg0, arg1, arg2)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByExternalSource indicates an expected call of FindByExternalSource.
func (mr *MockUserProviderMockRecorder) FindByExternalSource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByExternalSource", reflect.TypeOf((*MockUserProvider)(nil).FindByExternalSource), arg0, arg1, arg2)
}

// FindByID mocks base method.
func (m *MockUserProvider) FindByID(arg0 context.Context, arg1 string) (*security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserProviderMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserProvider)(nil).FindByID), arg0, arg1)
}

// FindByUsername mocks base method.
func (m *MockUserProvider) FindByUsername(arg0 context.Context, arg1 string) (*security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *MockUserProviderMockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*MockUserProvider)(nil).FindByUsername), arg0, arg1)
}

// FindWithoutPermission mocks base method.
func (m *MockUserProvider) FindWithoutPermission(arg0 context.Context, arg1 string) ([]security.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindWithoutPermission", arg0, arg1)
	ret0, _ := ret[0].([]security.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWithoutPermission indicates an expected call of FindWithoutPermission.
func (mr *MockUserProviderMockRecorder) FindWithoutPermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWithoutPermission", reflect.TypeOf((*MockUserProvider)(nil).FindWithoutPermission), arg0, arg1)
}

// Save mocks base method.
func (m *MockUserProvider) Save(arg0 context.Context, arg1 *security.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserProviderMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserProvider)(nil).Save), arg0, arg1)
}

// UpdateHashedPassword mocks base method.
func (m *MockUserProvider) UpdateHashedPassword(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHashedPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHashedPassword indicates an expected call of UpdateHashedPassword.
func (mr *MockUserProviderMockRecorder) UpdateHashedPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHashedPassword", reflect.TypeOf((*MockUserProvider)(nil).UpdateHashedPassword), arg0, arg1, arg2)
}
