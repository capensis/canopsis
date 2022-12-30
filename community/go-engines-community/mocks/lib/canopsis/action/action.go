// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action (interfaces: Adapter,DelayedScenarioManager,DelayedScenarioStorage,ScenarioExecutionStorage,ScenarioStorage,WorkerPool)

// Package mock_action is a generated GoMock package.
package mock_action

import (
	context "context"
	reflect "reflect"

	action "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	gomock "github.com/golang/mock/gomock"
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

// GetEnabled mocks base method.
func (m *MockAdapter) GetEnabled(arg0 context.Context) ([]action.Scenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabled", arg0)
	ret0, _ := ret[0].([]action.Scenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnabled indicates an expected call of GetEnabled.
func (mr *MockAdapterMockRecorder) GetEnabled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnabled", reflect.TypeOf((*MockAdapter)(nil).GetEnabled), arg0)
}

// GetEnabledByIDs mocks base method.
func (m *MockAdapter) GetEnabledByIDs(arg0 context.Context, arg1 []string) ([]action.Scenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabledByIDs", arg0, arg1)
	ret0, _ := ret[0].([]action.Scenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnabledByIDs indicates an expected call of GetEnabledByIDs.
func (mr *MockAdapterMockRecorder) GetEnabledByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnabledByIDs", reflect.TypeOf((*MockAdapter)(nil).GetEnabledByIDs), arg0, arg1)
}

// GetEnabledById mocks base method.
func (m *MockAdapter) GetEnabledById(arg0 context.Context, arg1 string) (action.Scenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabledById", arg0, arg1)
	ret0, _ := ret[0].(action.Scenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnabledById indicates an expected call of GetEnabledById.
func (mr *MockAdapterMockRecorder) GetEnabledById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnabledById", reflect.TypeOf((*MockAdapter)(nil).GetEnabledById), arg0, arg1)
}

// MockDelayedScenarioManager is a mock of DelayedScenarioManager interface.
type MockDelayedScenarioManager struct {
	ctrl     *gomock.Controller
	recorder *MockDelayedScenarioManagerMockRecorder
}

// MockDelayedScenarioManagerMockRecorder is the mock recorder for MockDelayedScenarioManager.
type MockDelayedScenarioManagerMockRecorder struct {
	mock *MockDelayedScenarioManager
}

// NewMockDelayedScenarioManager creates a new mock instance.
func NewMockDelayedScenarioManager(ctrl *gomock.Controller) *MockDelayedScenarioManager {
	mock := &MockDelayedScenarioManager{ctrl: ctrl}
	mock.recorder = &MockDelayedScenarioManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDelayedScenarioManager) EXPECT() *MockDelayedScenarioManagerMockRecorder {
	return m.recorder
}

// AddDelayedScenario mocks base method.
func (m *MockDelayedScenarioManager) AddDelayedScenario(arg0 context.Context, arg1 types.Alarm, arg2 action.Scenario, arg3 action.AdditionalData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDelayedScenario", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDelayedScenario indicates an expected call of AddDelayedScenario.
func (mr *MockDelayedScenarioManagerMockRecorder) AddDelayedScenario(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDelayedScenario", reflect.TypeOf((*MockDelayedScenarioManager)(nil).AddDelayedScenario), arg0, arg1, arg2, arg3)
}

// PauseDelayedScenarios mocks base method.
func (m *MockDelayedScenarioManager) PauseDelayedScenarios(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseDelayedScenarios", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PauseDelayedScenarios indicates an expected call of PauseDelayedScenarios.
func (mr *MockDelayedScenarioManagerMockRecorder) PauseDelayedScenarios(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseDelayedScenarios", reflect.TypeOf((*MockDelayedScenarioManager)(nil).PauseDelayedScenarios), arg0, arg1)
}

// ResumeDelayedScenarios mocks base method.
func (m *MockDelayedScenarioManager) ResumeDelayedScenarios(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResumeDelayedScenarios", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResumeDelayedScenarios indicates an expected call of ResumeDelayedScenarios.
func (mr *MockDelayedScenarioManagerMockRecorder) ResumeDelayedScenarios(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResumeDelayedScenarios", reflect.TypeOf((*MockDelayedScenarioManager)(nil).ResumeDelayedScenarios), arg0, arg1)
}

// Run mocks base method.
func (m *MockDelayedScenarioManager) Run(arg0 context.Context) (<-chan action.DelayedScenarioTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(<-chan action.DelayedScenarioTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockDelayedScenarioManagerMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockDelayedScenarioManager)(nil).Run), arg0)
}

// MockDelayedScenarioStorage is a mock of DelayedScenarioStorage interface.
type MockDelayedScenarioStorage struct {
	ctrl     *gomock.Controller
	recorder *MockDelayedScenarioStorageMockRecorder
}

// MockDelayedScenarioStorageMockRecorder is the mock recorder for MockDelayedScenarioStorage.
type MockDelayedScenarioStorageMockRecorder struct {
	mock *MockDelayedScenarioStorage
}

// NewMockDelayedScenarioStorage creates a new mock instance.
func NewMockDelayedScenarioStorage(ctrl *gomock.Controller) *MockDelayedScenarioStorage {
	mock := &MockDelayedScenarioStorage{ctrl: ctrl}
	mock.recorder = &MockDelayedScenarioStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDelayedScenarioStorage) EXPECT() *MockDelayedScenarioStorageMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockDelayedScenarioStorage) Add(arg0 context.Context, arg1 action.DelayedScenario) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockDelayedScenarioStorageMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockDelayedScenarioStorage)(nil).Add), arg0, arg1)
}

// Delete mocks base method.
func (m *MockDelayedScenarioStorage) Delete(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockDelayedScenarioStorageMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDelayedScenarioStorage)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockDelayedScenarioStorage) Get(arg0 context.Context, arg1 string) (*action.DelayedScenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*action.DelayedScenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDelayedScenarioStorageMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDelayedScenarioStorage)(nil).Get), arg0, arg1)
}

// GetAll mocks base method.
func (m *MockDelayedScenarioStorage) GetAll(arg0 context.Context) ([]action.DelayedScenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]action.DelayedScenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockDelayedScenarioStorageMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockDelayedScenarioStorage)(nil).GetAll), arg0)
}

// Update mocks base method.
func (m *MockDelayedScenarioStorage) Update(arg0 context.Context, arg1 action.DelayedScenario) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockDelayedScenarioStorageMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDelayedScenarioStorage)(nil).Update), arg0, arg1)
}

// MockScenarioExecutionStorage is a mock of ScenarioExecutionStorage interface.
type MockScenarioExecutionStorage struct {
	ctrl     *gomock.Controller
	recorder *MockScenarioExecutionStorageMockRecorder
}

// MockScenarioExecutionStorageMockRecorder is the mock recorder for MockScenarioExecutionStorage.
type MockScenarioExecutionStorageMockRecorder struct {
	mock *MockScenarioExecutionStorage
}

// NewMockScenarioExecutionStorage creates a new mock instance.
func NewMockScenarioExecutionStorage(ctrl *gomock.Controller) *MockScenarioExecutionStorage {
	mock := &MockScenarioExecutionStorage{ctrl: ctrl}
	mock.recorder = &MockScenarioExecutionStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScenarioExecutionStorage) EXPECT() *MockScenarioExecutionStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockScenarioExecutionStorage) Create(arg0 context.Context, arg1 action.ScenarioExecution) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockScenarioExecutionStorageMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).Create), arg0, arg1)
}

// Del mocks base method.
func (m *MockScenarioExecutionStorage) Del(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockScenarioExecutionStorageMockRecorder) Del(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).Del), arg0, arg1)
}

// Get mocks base method.
func (m *MockScenarioExecutionStorage) Get(arg0 context.Context, arg1 string) (*action.ScenarioExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*action.ScenarioExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockScenarioExecutionStorageMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).Get), arg0, arg1)
}

// GetAbandoned mocks base method.
func (m *MockScenarioExecutionStorage) GetAbandoned(arg0 context.Context) ([]action.ScenarioExecution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbandoned", arg0)
	ret0, _ := ret[0].([]action.ScenarioExecution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAbandoned indicates an expected call of GetAbandoned.
func (mr *MockScenarioExecutionStorageMockRecorder) GetAbandoned(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbandoned", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).GetAbandoned), arg0)
}

// Inc mocks base method.
func (m *MockScenarioExecutionStorage) Inc(arg0 context.Context, arg1 string, arg2 int64, arg3 bool) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inc", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inc indicates an expected call of Inc.
func (mr *MockScenarioExecutionStorageMockRecorder) Inc(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).Inc), arg0, arg1, arg2, arg3)
}

// Update mocks base method.
func (m *MockScenarioExecutionStorage) Update(arg0 context.Context, arg1 action.ScenarioExecution) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockScenarioExecutionStorageMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScenarioExecutionStorage)(nil).Update), arg0, arg1)
}

// MockScenarioStorage is a mock of ScenarioStorage interface.
type MockScenarioStorage struct {
	ctrl     *gomock.Controller
	recorder *MockScenarioStorageMockRecorder
}

// MockScenarioStorageMockRecorder is the mock recorder for MockScenarioStorage.
type MockScenarioStorageMockRecorder struct {
	mock *MockScenarioStorage
}

// NewMockScenarioStorage creates a new mock instance.
func NewMockScenarioStorage(ctrl *gomock.Controller) *MockScenarioStorage {
	mock := &MockScenarioStorage{ctrl: ctrl}
	mock.recorder = &MockScenarioStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScenarioStorage) EXPECT() *MockScenarioStorageMockRecorder {
	return m.recorder
}

// GetScenario mocks base method.
func (m *MockScenarioStorage) GetScenario(arg0 string) *action.Scenario {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScenario", arg0)
	ret0, _ := ret[0].(*action.Scenario)
	return ret0
}

// GetScenario indicates an expected call of GetScenario.
func (mr *MockScenarioStorageMockRecorder) GetScenario(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScenario", reflect.TypeOf((*MockScenarioStorage)(nil).GetScenario), arg0)
}

// GetTriggeredScenarios mocks base method.
func (m *MockScenarioStorage) GetTriggeredScenarios(arg0 []string, arg1 types.Alarm) ([]action.Scenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTriggeredScenarios", arg0, arg1)
	ret0, _ := ret[0].([]action.Scenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTriggeredScenarios indicates an expected call of GetTriggeredScenarios.
func (mr *MockScenarioStorageMockRecorder) GetTriggeredScenarios(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTriggeredScenarios", reflect.TypeOf((*MockScenarioStorage)(nil).GetTriggeredScenarios), arg0, arg1)
}

// ReloadScenarios mocks base method.
func (m *MockScenarioStorage) ReloadScenarios(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadScenarios", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadScenarios indicates an expected call of ReloadScenarios.
func (mr *MockScenarioStorageMockRecorder) ReloadScenarios(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadScenarios", reflect.TypeOf((*MockScenarioStorage)(nil).ReloadScenarios), arg0)
}

// RunDelayedScenarios mocks base method.
func (m *MockScenarioStorage) RunDelayedScenarios(arg0 context.Context, arg1 []string, arg2 types.Alarm, arg3 types.Entity, arg4 action.AdditionalData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunDelayedScenarios", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunDelayedScenarios indicates an expected call of RunDelayedScenarios.
func (mr *MockScenarioStorageMockRecorder) RunDelayedScenarios(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunDelayedScenarios", reflect.TypeOf((*MockScenarioStorage)(nil).RunDelayedScenarios), arg0, arg1, arg2, arg3, arg4)
}

// MockWorkerPool is a mock of WorkerPool interface.
type MockWorkerPool struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerPoolMockRecorder
}

// MockWorkerPoolMockRecorder is the mock recorder for MockWorkerPool.
type MockWorkerPoolMockRecorder struct {
	mock *MockWorkerPool
}

// NewMockWorkerPool creates a new mock instance.
func NewMockWorkerPool(ctrl *gomock.Controller) *MockWorkerPool {
	mock := &MockWorkerPool{ctrl: ctrl}
	mock.recorder = &MockWorkerPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkerPool) EXPECT() *MockWorkerPoolMockRecorder {
	return m.recorder
}

// RunWorkers mocks base method.
func (m *MockWorkerPool) RunWorkers(arg0 context.Context, arg1 <-chan action.Task) (<-chan action.TaskResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunWorkers", arg0, arg1)
	ret0, _ := ret[0].(<-chan action.TaskResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunWorkers indicates an expected call of RunWorkers.
func (mr *MockWorkerPoolMockRecorder) RunWorkers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunWorkers", reflect.TypeOf((*MockWorkerPool)(nil).RunWorkers), arg0, arg1)
}
