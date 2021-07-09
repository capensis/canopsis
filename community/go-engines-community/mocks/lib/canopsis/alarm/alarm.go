// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm (interfaces: Adapter,Service,EventProcessor)

// Package mock_alarm is a generated GoMock package.
package mock_alarm

import (
	context "context"
	config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
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

// ArchiveResolvedAlarms mocks base method
func (m *MockAdapter) ArchiveResolvedAlarms(arg0 context.Context, arg1 time.Duration) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveResolvedAlarms indicates an expected call of ArchiveResolvedAlarms
func (mr *MockAdapterMockRecorder) ArchiveResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).ArchiveResolvedAlarms), arg0, arg1)
}

// CopyAlarmToResolvedCollection mocks base method
func (m *MockAdapter) CopyAlarmToResolvedCollection(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyAlarmToResolvedCollection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyAlarmToResolvedCollection indicates an expected call of CopyAlarmToResolvedCollection
func (mr *MockAdapterMockRecorder) CopyAlarmToResolvedCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyAlarmToResolvedCollection", reflect.TypeOf((*MockAdapter)(nil).CopyAlarmToResolvedCollection), arg0, arg1)
}

// CountResolvedAlarm mocks base method
func (m *MockAdapter) CountResolvedAlarm(arg0 context.Context, arg1 []string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountResolvedAlarm", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountResolvedAlarm indicates an expected call of CountResolvedAlarm
func (mr *MockAdapterMockRecorder) CountResolvedAlarm(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountResolvedAlarm", reflect.TypeOf((*MockAdapter)(nil).CountResolvedAlarm), arg0, arg1)
}

// DeleteArchivedResolvedAlarms mocks base method
func (m *MockAdapter) DeleteArchivedResolvedAlarms(arg0 context.Context, arg1 time.Duration) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteArchivedResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteArchivedResolvedAlarms indicates an expected call of DeleteArchivedResolvedAlarms
func (mr *MockAdapterMockRecorder) DeleteArchivedResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteArchivedResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).DeleteArchivedResolvedAlarms), arg0, arg1)
}

// DeleteResolvedAlarms mocks base method
func (m *MockAdapter) DeleteResolvedAlarms(arg0 context.Context, arg1 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteResolvedAlarms indicates an expected call of DeleteResolvedAlarms
func (mr *MockAdapterMockRecorder) DeleteResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).DeleteResolvedAlarms), arg0, arg1)
}

// GetAlarmsByID mocks base method
func (m *MockAdapter) GetAlarmsByID(arg0 context.Context, arg1 string) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsByID", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsByID indicates an expected call of GetAlarmsByID
func (mr *MockAdapterMockRecorder) GetAlarmsByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsByID", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsByID), arg0, arg1)
}

// GetAlarmsWithCancelMark mocks base method
func (m *MockAdapter) GetAlarmsWithCancelMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithCancelMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithCancelMark indicates an expected call of GetAlarmsWithCancelMark
func (mr *MockAdapterMockRecorder) GetAlarmsWithCancelMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithCancelMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithCancelMark), arg0)
}

// GetAlarmsWithDoneMark mocks base method
func (m *MockAdapter) GetAlarmsWithDoneMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithDoneMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithDoneMark indicates an expected call of GetAlarmsWithDoneMark
func (mr *MockAdapterMockRecorder) GetAlarmsWithDoneMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithDoneMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithDoneMark), arg0)
}

// GetAlarmsWithFlappingStatus mocks base method
func (m *MockAdapter) GetAlarmsWithFlappingStatus(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithFlappingStatus", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithFlappingStatus indicates an expected call of GetAlarmsWithFlappingStatus
func (mr *MockAdapterMockRecorder) GetAlarmsWithFlappingStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithFlappingStatus", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithFlappingStatus), arg0)
}

// GetAlarmsWithSnoozeMark mocks base method
func (m *MockAdapter) GetAlarmsWithSnoozeMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithSnoozeMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithSnoozeMark indicates an expected call of GetAlarmsWithSnoozeMark
func (mr *MockAdapterMockRecorder) GetAlarmsWithSnoozeMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithSnoozeMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithSnoozeMark), arg0)
}

// GetAlarmsWithoutTicketByComponent mocks base method
func (m *MockAdapter) GetAlarmsWithoutTicketByComponent(arg0 context.Context, arg1 string) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithoutTicketByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithoutTicketByComponent indicates an expected call of GetAlarmsWithoutTicketByComponent
func (mr *MockAdapterMockRecorder) GetAlarmsWithoutTicketByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithoutTicketByComponent", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithoutTicketByComponent), arg0, arg1)
}

// GetAllOpenedResourceAlarmsByComponent mocks base method
func (m *MockAdapter) GetAllOpenedResourceAlarmsByComponent(arg0 context.Context, arg1 string) ([]types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOpenedResourceAlarmsByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOpenedResourceAlarmsByComponent indicates an expected call of GetAllOpenedResourceAlarmsByComponent
func (mr *MockAdapterMockRecorder) GetAllOpenedResourceAlarmsByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOpenedResourceAlarmsByComponent", reflect.TypeOf((*MockAdapter)(nil).GetAllOpenedResourceAlarmsByComponent), arg0, arg1)
}

// GetCountOpenedAlarmsByIDs mocks base method
func (m *MockAdapter) GetCountOpenedAlarmsByIDs(arg0 context.Context, arg1 []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountOpenedAlarmsByIDs", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountOpenedAlarmsByIDs indicates an expected call of GetCountOpenedAlarmsByIDs
func (mr *MockAdapterMockRecorder) GetCountOpenedAlarmsByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountOpenedAlarmsByIDs", reflect.TypeOf((*MockAdapter)(nil).GetCountOpenedAlarmsByIDs), arg0, arg1)
}

// GetLastAlarm mocks base method
func (m *MockAdapter) GetLastAlarm(arg0 context.Context, arg1, arg2, arg3 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAlarm", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAlarm indicates an expected call of GetLastAlarm
func (mr *MockAdapterMockRecorder) GetLastAlarm(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAlarm", reflect.TypeOf((*MockAdapter)(nil).GetLastAlarm), arg0, arg1, arg2, arg3)
}

// GetLastAlarmByEntityID mocks base method
func (m *MockAdapter) GetLastAlarmByEntityID(arg0 context.Context, arg1 string) (*types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAlarmByEntityID", arg0, arg1)
	ret0, _ := ret[0].(*types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAlarmByEntityID indicates an expected call of GetLastAlarmByEntityID
func (mr *MockAdapterMockRecorder) GetLastAlarmByEntityID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAlarmByEntityID", reflect.TypeOf((*MockAdapter)(nil).GetLastAlarmByEntityID), arg0, arg1)
}

// GetOpenedAlarm mocks base method
func (m *MockAdapter) GetOpenedAlarm(arg0 context.Context, arg1, arg2, arg3 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarm", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarm indicates an expected call of GetOpenedAlarm
func (mr *MockAdapterMockRecorder) GetOpenedAlarm(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarm", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarm), arg0, arg1, arg2, arg3)
}

// GetOpenedAlarmByAlarmId mocks base method
func (m *MockAdapter) GetOpenedAlarmByAlarmId(arg0 context.Context, arg1 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmByAlarmId", arg0, arg1)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmByAlarmId indicates an expected call of GetOpenedAlarmByAlarmId
func (mr *MockAdapterMockRecorder) GetOpenedAlarmByAlarmId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmByAlarmId", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmByAlarmId), arg0, arg1)
}

// GetOpenedAlarmsByAlarmIDs mocks base method
func (m *MockAdapter) GetOpenedAlarmsByAlarmIDs(arg0 context.Context, arg1 []string, arg2 *[]types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByAlarmIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsByAlarmIDs indicates an expected call of GetOpenedAlarmsByAlarmIDs
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByAlarmIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByAlarmIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByAlarmIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsByConnectorIdleRules mocks base method
func (m *MockAdapter) GetOpenedAlarmsByConnectorIdleRules(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByConnectorIdleRules", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsByConnectorIdleRules indicates an expected call of GetOpenedAlarmsByConnectorIdleRules
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByConnectorIdleRules(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByConnectorIdleRules", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByConnectorIdleRules), arg0)
}

// GetOpenedAlarmsByIDs mocks base method
func (m *MockAdapter) GetOpenedAlarmsByIDs(arg0 context.Context, arg1 []string, arg2 *[]types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsByIDs indicates an expected call of GetOpenedAlarmsByIDs
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithEntityByAlarmIDs mocks base method
func (m *MockAdapter) GetOpenedAlarmsWithEntityByAlarmIDs(arg0 context.Context, arg1 []string, arg2 *[]types.AlarmWithEntity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntityByAlarmIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsWithEntityByAlarmIDs indicates an expected call of GetOpenedAlarmsWithEntityByAlarmIDs
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntityByAlarmIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntityByAlarmIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntityByAlarmIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithEntityByIDs mocks base method
func (m *MockAdapter) GetOpenedAlarmsWithEntityByIDs(arg0 context.Context, arg1 []string, arg2 *[]types.AlarmWithEntity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntityByIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsWithEntityByIDs indicates an expected call of GetOpenedAlarmsWithEntityByIDs
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntityByIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntityByIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntityByIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithLastDatesBefore mocks base method
func (m *MockAdapter) GetOpenedAlarmsWithLastDatesBefore(arg0 context.Context, arg1 types.CpsTime) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithLastDatesBefore", arg0, arg1)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsWithLastDatesBefore indicates an expected call of GetOpenedAlarmsWithLastDatesBefore
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithLastDatesBefore(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithLastDatesBefore", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithLastDatesBefore), arg0, arg1)
}

// GetOpenedMetaAlarm mocks base method
func (m *MockAdapter) GetOpenedMetaAlarm(arg0 context.Context, arg1, arg2 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedMetaAlarm", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedMetaAlarm indicates an expected call of GetOpenedMetaAlarm
func (mr *MockAdapterMockRecorder) GetOpenedMetaAlarm(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedMetaAlarm", reflect.TypeOf((*MockAdapter)(nil).GetOpenedMetaAlarm), arg0, arg1, arg2)
}

// GetUnacknowledgedAlarmsByComponent mocks base method
func (m *MockAdapter) GetUnacknowledgedAlarmsByComponent(arg0 context.Context, arg1 string) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnacknowledgedAlarmsByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnacknowledgedAlarmsByComponent indicates an expected call of GetUnacknowledgedAlarmsByComponent
func (mr *MockAdapterMockRecorder) GetUnacknowledgedAlarmsByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnacknowledgedAlarmsByComponent", reflect.TypeOf((*MockAdapter)(nil).GetUnacknowledgedAlarmsByComponent), arg0, arg1)
}

// GetUnresolved mocks base method
func (m *MockAdapter) GetUnresolved(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnresolved", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnresolved indicates an expected call of GetUnresolved
func (mr *MockAdapterMockRecorder) GetUnresolved(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnresolved", reflect.TypeOf((*MockAdapter)(nil).GetUnresolved), arg0)
}

// Insert mocks base method
func (m *MockAdapter) Insert(arg0 context.Context, arg1 types.Alarm) error {
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

// MassPartialUpdateOpen mocks base method
func (m *MockAdapter) MassPartialUpdateOpen(arg0 context.Context, arg1 *types.Alarm, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassPartialUpdateOpen", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassPartialUpdateOpen indicates an expected call of MassPartialUpdateOpen
func (mr *MockAdapterMockRecorder) MassPartialUpdateOpen(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassPartialUpdateOpen", reflect.TypeOf((*MockAdapter)(nil).MassPartialUpdateOpen), arg0, arg1, arg2)
}

// MassUpdate mocks base method
func (m *MockAdapter) MassUpdate(arg0 context.Context, arg1 []types.Alarm, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassUpdate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassUpdate indicates an expected call of MassUpdate
func (mr *MockAdapterMockRecorder) MassUpdate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassUpdate", reflect.TypeOf((*MockAdapter)(nil).MassUpdate), arg0, arg1, arg2)
}

// PartialUpdateOpen mocks base method
func (m *MockAdapter) PartialUpdateOpen(arg0 context.Context, arg1 *types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialUpdateOpen", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PartialUpdateOpen indicates an expected call of PartialUpdateOpen
func (mr *MockAdapterMockRecorder) PartialUpdateOpen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialUpdateOpen", reflect.TypeOf((*MockAdapter)(nil).PartialUpdateOpen), arg0, arg1)
}

// Update mocks base method
func (m *MockAdapter) Update(arg0 context.Context, arg1 types.Alarm) error {
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

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// ResolveAlarms mocks base method
func (m *MockService) ResolveAlarms(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveAlarms", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveAlarms indicates an expected call of ResolveAlarms
func (mr *MockServiceMockRecorder) ResolveAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveAlarms", reflect.TypeOf((*MockService)(nil).ResolveAlarms), arg0, arg1)
}

// ResolveCancels mocks base method
func (m *MockService) ResolveCancels(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveCancels", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveCancels indicates an expected call of ResolveCancels
func (mr *MockServiceMockRecorder) ResolveCancels(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveCancels", reflect.TypeOf((*MockService)(nil).ResolveCancels), arg0, arg1)
}

// ResolveDone mocks base method
func (m *MockService) ResolveDone(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveDone", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveDone indicates an expected call of ResolveDone
func (mr *MockServiceMockRecorder) ResolveDone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveDone", reflect.TypeOf((*MockService)(nil).ResolveDone), arg0)
}

// ResolveSnoozes mocks base method
func (m *MockService) ResolveSnoozes(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveSnoozes", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveSnoozes indicates an expected call of ResolveSnoozes
func (mr *MockServiceMockRecorder) ResolveSnoozes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveSnoozes", reflect.TypeOf((*MockService)(nil).ResolveSnoozes), arg0, arg1)
}

// UpdateFlappingAlarms mocks base method
func (m *MockService) UpdateFlappingAlarms(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFlappingAlarms", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFlappingAlarms indicates an expected call of UpdateFlappingAlarms
func (mr *MockServiceMockRecorder) UpdateFlappingAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlappingAlarms", reflect.TypeOf((*MockService)(nil).UpdateFlappingAlarms), arg0, arg1)
}

// MockEventProcessor is a mock of EventProcessor interface
type MockEventProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockEventProcessorMockRecorder
}

// MockEventProcessorMockRecorder is the mock recorder for MockEventProcessor
type MockEventProcessorMockRecorder struct {
	mock *MockEventProcessor
}

// NewMockEventProcessor creates a new mock instance
func NewMockEventProcessor(ctrl *gomock.Controller) *MockEventProcessor {
	mock := &MockEventProcessor{ctrl: ctrl}
	mock.recorder = &MockEventProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventProcessor) EXPECT() *MockEventProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method
func (m *MockEventProcessor) Process(arg0 context.Context, arg1 *types.Event) (types.AlarmChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1)
	ret0, _ := ret[0].(types.AlarmChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process
func (mr *MockEventProcessorMockRecorder) Process(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockEventProcessor)(nil).Process), arg0, arg1)
}
