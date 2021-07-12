// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity (interfaces: Adapter)

// Package mock_entity is a generated GoMock package.
package mock_entity

import (
	context "context"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	gomock "github.com/golang/mock/gomock"
	mongo0 "go.mongodb.org/mongo-driver/mongo"
	reflect "reflect"
)

// MockAdapter is a mock of Adapter interface
type MockAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockAdapterMockRecorder
}

// MockAdapterMockRecorder is the mock recorder for MockAdapter
type MockAdapterMockRecorder struct {
	mock *MockAdapter
}

// NewMockAdapter creates a new mock instance
func NewMockAdapter(ctrl *gomock.Controller) *MockAdapter {
	mock := &MockAdapter{ctrl: ctrl}
	mock.recorder = &MockAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdapter) EXPECT() *MockAdapterMockRecorder {
	return m.recorder
}

// AddImpactByQuery mocks base method
func (m *MockAdapter) AddImpactByQuery(arg0 context.Context, arg1 interface{}, arg2 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImpactByQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddImpactByQuery indicates an expected call of AddImpactByQuery
func (mr *MockAdapterMockRecorder) AddImpactByQuery(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImpactByQuery", reflect.TypeOf((*MockAdapter)(nil).AddImpactByQuery), arg0, arg1, arg2)
}

// AddImpacts mocks base method
func (m *MockAdapter) AddImpacts(arg0 context.Context, arg1, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImpacts", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImpacts indicates an expected call of AddImpacts
func (mr *MockAdapterMockRecorder) AddImpacts(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImpacts", reflect.TypeOf((*MockAdapter)(nil).AddImpacts), arg0, arg1, arg2)
}

// AddInfos mocks base method
func (m *MockAdapter) AddInfos(arg0 context.Context, arg1 string, arg2 map[string]types.Info) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddInfos", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddInfos indicates an expected call of AddInfos
func (mr *MockAdapterMockRecorder) AddInfos(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddInfos", reflect.TypeOf((*MockAdapter)(nil).AddInfos), arg0, arg1, arg2)
}

// Bulk mocks base method
func (m *MockAdapter) Bulk(arg0 context.Context, arg1 []mongo0.WriteModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bulk", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Bulk indicates an expected call of Bulk
func (mr *MockAdapterMockRecorder) Bulk(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bulk", reflect.TypeOf((*MockAdapter)(nil).Bulk), arg0, arg1)
}

// Count mocks base method
func (m *MockAdapter) Count(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockAdapterMockRecorder) Count(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockAdapter)(nil).Count), arg0)
}

// FindByIDs mocks base method
func (m *MockAdapter) FindByIDs(arg0 context.Context, arg1 []string) ([]types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDs", arg0, arg1)
	ret0, _ := ret[0].([]types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDs indicates an expected call of FindByIDs
func (mr *MockAdapterMockRecorder) FindByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDs", reflect.TypeOf((*MockAdapter)(nil).FindByIDs), arg0, arg1)
}

// FindComponentForResource mocks base method
func (m *MockAdapter) FindComponentForResource(arg0 context.Context, arg1 string) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindComponentForResource", arg0, arg1)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindComponentForResource indicates an expected call of FindComponentForResource
func (mr *MockAdapterMockRecorder) FindComponentForResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindComponentForResource", reflect.TypeOf((*MockAdapter)(nil).FindComponentForResource), arg0, arg1)
}

// FindConnectorForComponent mocks base method
func (m *MockAdapter) FindConnectorForComponent(arg0 context.Context, arg1 string) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindConnectorForComponent", arg0, arg1)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindConnectorForComponent indicates an expected call of FindConnectorForComponent
func (mr *MockAdapterMockRecorder) FindConnectorForComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindConnectorForComponent", reflect.TypeOf((*MockAdapter)(nil).FindConnectorForComponent), arg0, arg1)
}

// FindConnectorForResource mocks base method
func (m *MockAdapter) FindConnectorForResource(arg0 context.Context, arg1 string) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindConnectorForResource", arg0, arg1)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindConnectorForResource indicates an expected call of FindConnectorForResource
func (mr *MockAdapterMockRecorder) FindConnectorForResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindConnectorForResource", reflect.TypeOf((*MockAdapter)(nil).FindConnectorForResource), arg0, arg1)
}

// Get mocks base method
func (m *MockAdapter) Get(arg0 context.Context, arg1 string) (types.Entity, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(types.Entity)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockAdapterMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdapter)(nil).Get), arg0, arg1)
}

// GetAllWithLastUpdateDateBefore mocks base method
func (m *MockAdapter) GetAllWithLastUpdateDateBefore(arg0 context.Context, arg1 types.CpsTime, arg2 []string) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithLastUpdateDateBefore", arg0, arg1, arg2)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithLastUpdateDateBefore indicates an expected call of GetAllWithLastUpdateDateBefore
func (mr *MockAdapterMockRecorder) GetAllWithLastUpdateDateBefore(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithLastUpdateDateBefore", reflect.TypeOf((*MockAdapter)(nil).GetAllWithLastUpdateDateBefore), arg0, arg1, arg2)
}

// GetEntityByID mocks base method
func (m *MockAdapter) GetEntityByID(arg0 context.Context, arg1 string) (types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByID", arg0, arg1)
	ret0, _ := ret[0].(types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByID indicates an expected call of GetEntityByID
func (mr *MockAdapterMockRecorder) GetEntityByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByID", reflect.TypeOf((*MockAdapter)(nil).GetEntityByID), arg0, arg1)
}

// GetIDs mocks base method
func (m *MockAdapter) GetIDs(arg0 context.Context, arg1 map[string]interface{}, arg2 *[]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetIDs indicates an expected call of GetIDs
func (mr *MockAdapterMockRecorder) GetIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIDs", reflect.TypeOf((*MockAdapter)(nil).GetIDs), arg0, arg1, arg2)
}

// GetImpactedServicesInfo mocks base method
func (m *MockAdapter) GetImpactedServicesInfo(arg0 context.Context) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImpactedServicesInfo", arg0)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImpactedServicesInfo indicates an expected call of GetImpactedServicesInfo
func (mr *MockAdapterMockRecorder) GetImpactedServicesInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImpactedServicesInfo", reflect.TypeOf((*MockAdapter)(nil).GetImpactedServicesInfo), arg0)
}

// GetWithIdleSince mocks base method
func (m *MockAdapter) GetWithIdleSince(arg0 context.Context) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithIdleSince", arg0)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithIdleSince indicates an expected call of GetWithIdleSince
func (mr *MockAdapterMockRecorder) GetWithIdleSince(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithIdleSince", reflect.TypeOf((*MockAdapter)(nil).GetWithIdleSince), arg0)
}

// Insert mocks base method
func (m *MockAdapter) Insert(arg0 context.Context, arg1 types.Entity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockAdapterMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAdapter)(nil).Insert), arg0, arg1)
}

// RemoveImpactByQuery mocks base method
func (m *MockAdapter) RemoveImpactByQuery(arg0 context.Context, arg1 interface{}, arg2 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveImpactByQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveImpactByQuery indicates an expected call of RemoveImpactByQuery
func (mr *MockAdapterMockRecorder) RemoveImpactByQuery(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveImpactByQuery", reflect.TypeOf((*MockAdapter)(nil).RemoveImpactByQuery), arg0, arg1, arg2)
}

// RemoveImpacts mocks base method
func (m *MockAdapter) RemoveImpacts(arg0 context.Context, arg1, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveImpacts", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveImpacts indicates an expected call of RemoveImpacts
func (mr *MockAdapterMockRecorder) RemoveImpacts(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveImpacts", reflect.TypeOf((*MockAdapter)(nil).RemoveImpacts), arg0, arg1, arg2)
}

// Update mocks base method
func (m *MockAdapter) Update(arg0 context.Context, arg1 types.Entity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockAdapterMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAdapter)(nil).Update), arg0, arg1)
}

// UpdateComponentInfos mocks base method
func (m *MockAdapter) UpdateComponentInfos(arg0 context.Context, arg1, arg2 string) (map[string]types.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComponentInfos", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]types.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComponentInfos indicates an expected call of UpdateComponentInfos
func (mr *MockAdapterMockRecorder) UpdateComponentInfos(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComponentInfos", reflect.TypeOf((*MockAdapter)(nil).UpdateComponentInfos), arg0, arg1, arg2)
}

// UpdateComponentInfosByComponent mocks base method
func (m *MockAdapter) UpdateComponentInfosByComponent(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComponentInfosByComponent", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComponentInfosByComponent indicates an expected call of UpdateComponentInfosByComponent
func (mr *MockAdapterMockRecorder) UpdateComponentInfosByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComponentInfosByComponent", reflect.TypeOf((*MockAdapter)(nil).UpdateComponentInfosByComponent), arg0, arg1)
}

// UpdateIdleFields mocks base method
func (m *MockAdapter) UpdateIdleFields(arg0 context.Context, arg1 string, arg2 *types.CpsTime, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIdleFields", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateIdleFields indicates an expected call of UpdateIdleFields
func (mr *MockAdapterMockRecorder) UpdateIdleFields(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIdleFields", reflect.TypeOf((*MockAdapter)(nil).UpdateIdleFields), arg0, arg1, arg2, arg3)
}

// UpdateLastEventDate mocks base method
func (m *MockAdapter) UpdateLastEventDate(arg0 context.Context, arg1 []string, arg2 types.CpsTime) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastEventDate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastEventDate indicates an expected call of UpdateLastEventDate
func (mr *MockAdapterMockRecorder) UpdateLastEventDate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastEventDate", reflect.TypeOf((*MockAdapter)(nil).UpdateLastEventDate), arg0, arg1, arg2)
}

// UpsertMany mocks base method
func (m *MockAdapter) UpsertMany(arg0 context.Context, arg1 []types.Entity) (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMany", arg0, arg1)
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertMany indicates an expected call of UpsertMany
func (mr *MockAdapterMockRecorder) UpsertMany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMany", reflect.TypeOf((*MockAdapter)(nil).UpsertMany), arg0, arg1)
}
