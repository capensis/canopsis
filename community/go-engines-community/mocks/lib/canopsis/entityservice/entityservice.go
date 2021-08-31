// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice (interfaces: Adapter,CountersCache,Storage)

// Package mock_entityservice is a generated GoMock package.
package mock_entityservice

import (
	context "context"
	entityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
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

// AddDepends mocks base method
func (m *MockAdapter) AddDepends(arg0 context.Context, arg1 string, arg2 []string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDepends", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDepends indicates an expected call of AddDepends
func (mr *MockAdapterMockRecorder) AddDepends(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDepends", reflect.TypeOf((*MockAdapter)(nil).AddDepends), arg0, arg1, arg2)
}

// GetAll mocks base method
func (m *MockAdapter) GetAll(arg0 context.Context) ([]entityservice.EntityService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]entityservice.EntityService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockAdapterMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAdapter)(nil).GetAll), arg0)
}

// GetByID mocks base method
func (m *MockAdapter) GetByID(arg0 context.Context, arg1 string) (*entityservice.EntityService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*entityservice.EntityService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockAdapterMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAdapter)(nil).GetByID), arg0, arg1)
}

// GetCounters mocks base method
func (m *MockAdapter) GetCounters(arg0 context.Context, arg1 string) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounters", arg0, arg1)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounters indicates an expected call of GetCounters
func (mr *MockAdapterMockRecorder) GetCounters(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounters", reflect.TypeOf((*MockAdapter)(nil).GetCounters), arg0, arg1)
}

// GetEnabled mocks base method
func (m *MockAdapter) GetEnabled(arg0 context.Context) ([]entityservice.EntityService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabled", arg0)
	ret0, _ := ret[0].([]entityservice.EntityService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnabled indicates an expected call of GetEnabled
func (mr *MockAdapterMockRecorder) GetEnabled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnabled", reflect.TypeOf((*MockAdapter)(nil).GetEnabled), arg0)
}

// GetValid mocks base method
func (m *MockAdapter) GetValid(arg0 context.Context) ([]entityservice.EntityService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValid", arg0)
	ret0, _ := ret[0].([]entityservice.EntityService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValid indicates an expected call of GetValid
func (mr *MockAdapterMockRecorder) GetValid(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValid", reflect.TypeOf((*MockAdapter)(nil).GetValid), arg0)
}

// RemoveDependByQuery mocks base method
func (m *MockAdapter) RemoveDependByQuery(arg0 context.Context, arg1 interface{}, arg2 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDependByQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveDependByQuery indicates an expected call of RemoveDependByQuery
func (mr *MockAdapterMockRecorder) RemoveDependByQuery(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDependByQuery", reflect.TypeOf((*MockAdapter)(nil).RemoveDependByQuery), arg0, arg1, arg2)
}

// RemoveDepends mocks base method
func (m *MockAdapter) RemoveDepends(arg0 context.Context, arg1 string, arg2 []string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDepends", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveDepends indicates an expected call of RemoveDepends
func (mr *MockAdapterMockRecorder) RemoveDepends(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDepends", reflect.TypeOf((*MockAdapter)(nil).RemoveDepends), arg0, arg1, arg2)
}

// UpdateBulk mocks base method
func (m *MockAdapter) UpdateBulk(arg0 context.Context, arg1 []mongo0.WriteModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBulk", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBulk indicates an expected call of UpdateBulk
func (mr *MockAdapterMockRecorder) UpdateBulk(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBulk", reflect.TypeOf((*MockAdapter)(nil).UpdateBulk), arg0, arg1)
}

// UpdateCounters mocks base method
func (m *MockAdapter) UpdateCounters(arg0 context.Context, arg1 string, arg2 entityservice.AlarmCounters) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounters", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounters indicates an expected call of UpdateCounters
func (mr *MockAdapterMockRecorder) UpdateCounters(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounters", reflect.TypeOf((*MockAdapter)(nil).UpdateCounters), arg0, arg1, arg2)
}

// MockCountersCache is a mock of CountersCache interface
type MockCountersCache struct {
	ctrl     *gomock.Controller
	recorder *MockCountersCacheMockRecorder
}

// MockCountersCacheMockRecorder is the mock recorder for MockCountersCache
type MockCountersCacheMockRecorder struct {
	mock *MockCountersCache
}

// NewMockCountersCache creates a new mock instance
func NewMockCountersCache(ctrl *gomock.Controller) *MockCountersCache {
	mock := &MockCountersCache{ctrl: ctrl}
	mock.recorder = &MockCountersCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCountersCache) EXPECT() *MockCountersCacheMockRecorder {
	return m.recorder
}

// ClearAll mocks base method
func (m *MockCountersCache) ClearAll(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearAll indicates an expected call of ClearAll
func (mr *MockCountersCacheMockRecorder) ClearAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearAll", reflect.TypeOf((*MockCountersCache)(nil).ClearAll), arg0)
}

// KeepOnly mocks base method
func (m *MockCountersCache) KeepOnly(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepOnly", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// KeepOnly indicates an expected call of KeepOnly
func (mr *MockCountersCacheMockRecorder) KeepOnly(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepOnly", reflect.TypeOf((*MockCountersCache)(nil).KeepOnly), arg0, arg1)
}

// Remove mocks base method
func (m *MockCountersCache) Remove(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockCountersCacheMockRecorder) Remove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockCountersCache)(nil).Remove), arg0, arg1)
}

// RemoveAndGet mocks base method
func (m *MockCountersCache) RemoveAndGet(arg0 context.Context, arg1 string) (*entityservice.AlarmCounters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAndGet", arg0, arg1)
	ret0, _ := ret[0].(*entityservice.AlarmCounters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveAndGet indicates an expected call of RemoveAndGet
func (mr *MockCountersCacheMockRecorder) RemoveAndGet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAndGet", reflect.TypeOf((*MockCountersCache)(nil).RemoveAndGet), arg0, arg1)
}

// Replace mocks base method
func (m *MockCountersCache) Replace(arg0 context.Context, arg1 string, arg2 entityservice.AlarmCounters) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replace", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replace indicates an expected call of Replace
func (mr *MockCountersCacheMockRecorder) Replace(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replace", reflect.TypeOf((*MockCountersCache)(nil).Replace), arg0, arg1, arg2)
}

// Update mocks base method
func (m *MockCountersCache) Update(arg0 context.Context, arg1 map[string]entityservice.AlarmCounters) (map[string]entityservice.AlarmCounters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(map[string]entityservice.AlarmCounters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockCountersCacheMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCountersCache)(nil).Update), arg0, arg1)
}

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockStorage) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockStorageMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorage)(nil).Delete), arg0, arg1)
}

// Get mocks base method
func (m *MockStorage) Get(arg0 context.Context, arg1 string) (*entityservice.ServiceData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*entityservice.ServiceData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockStorageMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), arg0, arg1)
}

// Load mocks base method
func (m *MockStorage) Load(arg0 context.Context) ([]entityservice.ServiceData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].([]entityservice.ServiceData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load
func (mr *MockStorageMockRecorder) Load(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStorage)(nil).Load), arg0)
}

// Save mocks base method
func (m *MockStorage) Save(arg0 context.Context, arg1 entityservice.ServiceData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockStorageMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStorage)(nil).Save), arg0, arg1)
}

// SaveAll mocks base method
func (m *MockStorage) SaveAll(arg0 context.Context, arg1 []entityservice.ServiceData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAll indicates an expected call of SaveAll
func (mr *MockStorageMockRecorder) SaveAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAll", reflect.TypeOf((*MockStorage)(nil).SaveAll), arg0, arg1)
}
