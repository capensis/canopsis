// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter (interfaces: RuleApplicator,RuleAdapter,RuleApplicatorContainer,ExternalDataGetter,Service,ActionProcessor,FailureService,EventCounter)

// Package mock_eventfilter is a generated GoMock package.
package mock_eventfilter

import (
	context "context"
	reflect "reflect"

	eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	time "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	gomock "github.com/golang/mock/gomock"
)

// MockRuleApplicator is a mock of RuleApplicator interface.
type MockRuleApplicator struct {
	ctrl     *gomock.Controller
	recorder *MockRuleApplicatorMockRecorder
}

// MockRuleApplicatorMockRecorder is the mock recorder for MockRuleApplicator.
type MockRuleApplicatorMockRecorder struct {
	mock *MockRuleApplicator
}

// NewMockRuleApplicator creates a new mock instance.
func NewMockRuleApplicator(ctrl *gomock.Controller) *MockRuleApplicator {
	mock := &MockRuleApplicator{ctrl: ctrl}
	mock.recorder = &MockRuleApplicatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRuleApplicator) EXPECT() *MockRuleApplicatorMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockRuleApplicator) Apply(arg0 context.Context, arg1 eventfilter.ParsedRule, arg2 types.Event, arg3 eventfilter.RegexMatch) (string, types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(types.Event)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Apply indicates an expected call of Apply.
func (mr *MockRuleApplicatorMockRecorder) Apply(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockRuleApplicator)(nil).Apply), arg0, arg1, arg2, arg3)
}

// MockRuleAdapter is a mock of RuleAdapter interface.
type MockRuleAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockRuleAdapterMockRecorder
}

// MockRuleAdapterMockRecorder is the mock recorder for MockRuleAdapter.
type MockRuleAdapterMockRecorder struct {
	mock *MockRuleAdapter
}

// NewMockRuleAdapter creates a new mock instance.
func NewMockRuleAdapter(ctrl *gomock.Controller) *MockRuleAdapter {
	mock := &MockRuleAdapter{ctrl: ctrl}
	mock.recorder = &MockRuleAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRuleAdapter) EXPECT() *MockRuleAdapterMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockRuleAdapter) GetAll(arg0 context.Context) ([]eventfilter.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]eventfilter.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRuleAdapterMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRuleAdapter)(nil).GetAll), arg0)
}

// GetByTypes mocks base method.
func (m *MockRuleAdapter) GetByTypes(arg0 context.Context, arg1 []string) ([]eventfilter.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTypes", arg0, arg1)
	ret0, _ := ret[0].([]eventfilter.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTypes indicates an expected call of GetByTypes.
func (mr *MockRuleAdapterMockRecorder) GetByTypes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTypes", reflect.TypeOf((*MockRuleAdapter)(nil).GetByTypes), arg0, arg1)
}

// MockRuleApplicatorContainer is a mock of RuleApplicatorContainer interface.
type MockRuleApplicatorContainer struct {
	ctrl     *gomock.Controller
	recorder *MockRuleApplicatorContainerMockRecorder
}

// MockRuleApplicatorContainerMockRecorder is the mock recorder for MockRuleApplicatorContainer.
type MockRuleApplicatorContainerMockRecorder struct {
	mock *MockRuleApplicatorContainer
}

// NewMockRuleApplicatorContainer creates a new mock instance.
func NewMockRuleApplicatorContainer(ctrl *gomock.Controller) *MockRuleApplicatorContainer {
	mock := &MockRuleApplicatorContainer{ctrl: ctrl}
	mock.recorder = &MockRuleApplicatorContainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRuleApplicatorContainer) EXPECT() *MockRuleApplicatorContainerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRuleApplicatorContainer) Get(arg0 string) (eventfilter.RuleApplicator, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(eventfilter.RuleApplicator)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRuleApplicatorContainerMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRuleApplicatorContainer)(nil).Get), arg0)
}

// Set mocks base method.
func (m *MockRuleApplicatorContainer) Set(arg0 string, arg1 eventfilter.RuleApplicator) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1)
}

// Set indicates an expected call of Set.
func (mr *MockRuleApplicatorContainerMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRuleApplicatorContainer)(nil).Set), arg0, arg1)
}

// MockExternalDataGetter is a mock of ExternalDataGetter interface.
type MockExternalDataGetter struct {
	ctrl     *gomock.Controller
	recorder *MockExternalDataGetterMockRecorder
}

// MockExternalDataGetterMockRecorder is the mock recorder for MockExternalDataGetter.
type MockExternalDataGetterMockRecorder struct {
	mock *MockExternalDataGetter
}

// NewMockExternalDataGetter creates a new mock instance.
func NewMockExternalDataGetter(ctrl *gomock.Controller) *MockExternalDataGetter {
	mock := &MockExternalDataGetter{ctrl: ctrl}
	mock.recorder = &MockExternalDataGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternalDataGetter) EXPECT() *MockExternalDataGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockExternalDataGetter) Get(arg0 context.Context, arg1, arg2 string, arg3 types.Event, arg4 eventfilter.ParsedExternalDataParameters, arg5 eventfilter.Template) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockExternalDataGetterMockRecorder) Get(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockExternalDataGetter)(nil).Get), arg0, arg1, arg2, arg3, arg4, arg5)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// LoadRules mocks base method.
func (m *MockService) LoadRules(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadRules", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadRules indicates an expected call of LoadRules.
func (mr *MockServiceMockRecorder) LoadRules(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadRules", reflect.TypeOf((*MockService)(nil).LoadRules), arg0, arg1)
}

// ProcessEvent mocks base method.
func (m *MockService) ProcessEvent(arg0 context.Context, arg1 types.Event) (types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessEvent", arg0, arg1)
	ret0, _ := ret[0].(types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessEvent indicates an expected call of ProcessEvent.
func (mr *MockServiceMockRecorder) ProcessEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessEvent", reflect.TypeOf((*MockService)(nil).ProcessEvent), arg0, arg1)
}

// MockActionProcessor is a mock of ActionProcessor interface.
type MockActionProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockActionProcessorMockRecorder
}

// MockActionProcessorMockRecorder is the mock recorder for MockActionProcessor.
type MockActionProcessorMockRecorder struct {
	mock *MockActionProcessor
}

// NewMockActionProcessor creates a new mock instance.
func NewMockActionProcessor(ctrl *gomock.Controller) *MockActionProcessor {
	mock := &MockActionProcessor{ctrl: ctrl}
	mock.recorder = &MockActionProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActionProcessor) EXPECT() *MockActionProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockActionProcessor) Process(arg0 context.Context, arg1 string, arg2 eventfilter.ParsedAction, arg3 types.Event, arg4 eventfilter.RegexMatch, arg5 map[string]interface{}) (types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockActionProcessorMockRecorder) Process(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockActionProcessor)(nil).Process), arg0, arg1, arg2, arg3, arg4, arg5)
}

// MockFailureService is a mock of FailureService interface.
type MockFailureService struct {
	ctrl     *gomock.Controller
	recorder *MockFailureServiceMockRecorder
}

// MockFailureServiceMockRecorder is the mock recorder for MockFailureService.
type MockFailureServiceMockRecorder struct {
	mock *MockFailureService
}

// NewMockFailureService creates a new mock instance.
func NewMockFailureService(ctrl *gomock.Controller) *MockFailureService {
	mock := &MockFailureService{ctrl: ctrl}
	mock.recorder = &MockFailureServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFailureService) EXPECT() *MockFailureServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockFailureService) Add(arg0 string, arg1 int64, arg2 string, arg3 *types.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0, arg1, arg2, arg3)
}

// Add indicates an expected call of Add.
func (mr *MockFailureServiceMockRecorder) Add(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockFailureService)(nil).Add), arg0, arg1, arg2, arg3)
}

// Run mocks base method.
func (m *MockFailureService) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockFailureServiceMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockFailureService)(nil).Run), arg0)
}

// MockEventCounter is a mock of EventCounter interface.
type MockEventCounter struct {
	ctrl     *gomock.Controller
	recorder *MockEventCounterMockRecorder
}

// MockEventCounterMockRecorder is the mock recorder for MockEventCounter.
type MockEventCounterMockRecorder struct {
	mock *MockEventCounter
}

// NewMockEventCounter creates a new mock instance.
func NewMockEventCounter(ctrl *gomock.Controller) *MockEventCounter {
	mock := &MockEventCounter{ctrl: ctrl}
	mock.recorder = &MockEventCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventCounter) EXPECT() *MockEventCounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockEventCounter) Add(arg0 string, arg1 time.CpsTime) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0, arg1)
}

// Add indicates an expected call of Add.
func (mr *MockEventCounterMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockEventCounter)(nil).Add), arg0, arg1)
}

// Run mocks base method.
func (m *MockEventCounter) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockEventCounterMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockEventCounter)(nil).Run), arg0)
}
