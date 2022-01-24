// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm (interfaces: Adapter,Service,EventProcessor,ActivationService,MetaAlarmEventProcessor)

// Package mock_alarm is a generated GoMock package.
package mock_alarm

import (
	context "context"
	reflect "reflect"
	time "time"

	config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	types "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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

// ArchiveResolvedAlarms mocks base method.
func (m *MockAdapter) ArchiveResolvedAlarms(arg0 context.Context, arg1 types.CpsTime) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveResolvedAlarms indicates an expected call of ArchiveResolvedAlarms.
func (mr *MockAdapterMockRecorder) ArchiveResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).ArchiveResolvedAlarms), arg0, arg1)
}

// CopyAlarmToResolvedCollection mocks base method.
func (m *MockAdapter) CopyAlarmToResolvedCollection(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyAlarmToResolvedCollection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyAlarmToResolvedCollection indicates an expected call of CopyAlarmToResolvedCollection.
func (mr *MockAdapterMockRecorder) CopyAlarmToResolvedCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyAlarmToResolvedCollection", reflect.TypeOf((*MockAdapter)(nil).CopyAlarmToResolvedCollection), arg0, arg1)
}

// CountResolvedAlarm mocks base method.
func (m *MockAdapter) CountResolvedAlarm(arg0 context.Context, arg1 []string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountResolvedAlarm", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountResolvedAlarm indicates an expected call of CountResolvedAlarm.
func (mr *MockAdapterMockRecorder) CountResolvedAlarm(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountResolvedAlarm", reflect.TypeOf((*MockAdapter)(nil).CountResolvedAlarm), arg0, arg1)
}

// DeleteArchivedResolvedAlarms mocks base method.
func (m *MockAdapter) DeleteArchivedResolvedAlarms(arg0 context.Context, arg1 types.CpsTime) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteArchivedResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteArchivedResolvedAlarms indicates an expected call of DeleteArchivedResolvedAlarms.
func (mr *MockAdapterMockRecorder) DeleteArchivedResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteArchivedResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).DeleteArchivedResolvedAlarms), arg0, arg1)
}

// DeleteResolvedAlarms mocks base method.
func (m *MockAdapter) DeleteResolvedAlarms(arg0 context.Context, arg1 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteResolvedAlarms", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteResolvedAlarms indicates an expected call of DeleteResolvedAlarms.
func (mr *MockAdapterMockRecorder) DeleteResolvedAlarms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteResolvedAlarms", reflect.TypeOf((*MockAdapter)(nil).DeleteResolvedAlarms), arg0, arg1)
}

// FindToCheckPbehaviorInfo mocks base method.
func (m *MockAdapter) FindToCheckPbehaviorInfo(arg0 context.Context, arg1 types.CpsTime, arg2 []string) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindToCheckPbehaviorInfo", arg0, arg1, arg2)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindToCheckPbehaviorInfo indicates an expected call of FindToCheckPbehaviorInfo.
func (mr *MockAdapterMockRecorder) FindToCheckPbehaviorInfo(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindToCheckPbehaviorInfo", reflect.TypeOf((*MockAdapter)(nil).FindToCheckPbehaviorInfo), arg0, arg1, arg2)
}

// GetAlarmByAlarmId mocks base method.
func (m *MockAdapter) GetAlarmByAlarmId(arg0 context.Context, arg1 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmByAlarmId", arg0, arg1)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmByAlarmId indicates an expected call of GetAlarmByAlarmId.
func (mr *MockAdapterMockRecorder) GetAlarmByAlarmId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmByAlarmId", reflect.TypeOf((*MockAdapter)(nil).GetAlarmByAlarmId), arg0, arg1)
}

// GetAlarmsByID mocks base method.
func (m *MockAdapter) GetAlarmsByID(arg0 context.Context, arg1 string) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsByID", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsByID indicates an expected call of GetAlarmsByID.
func (mr *MockAdapterMockRecorder) GetAlarmsByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsByID", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsByID), arg0, arg1)
}

// GetAlarmsWithCancelMark mocks base method.
func (m *MockAdapter) GetAlarmsWithCancelMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithCancelMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithCancelMark indicates an expected call of GetAlarmsWithCancelMark.
func (mr *MockAdapterMockRecorder) GetAlarmsWithCancelMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithCancelMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithCancelMark), arg0)
}

// GetAlarmsWithDoneMark mocks base method.
func (m *MockAdapter) GetAlarmsWithDoneMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithDoneMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithDoneMark indicates an expected call of GetAlarmsWithDoneMark.
func (mr *MockAdapterMockRecorder) GetAlarmsWithDoneMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithDoneMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithDoneMark), arg0)
}

// GetAlarmsWithFlappingStatus mocks base method.
func (m *MockAdapter) GetAlarmsWithFlappingStatus(arg0 context.Context) ([]types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithFlappingStatus", arg0)
	ret0, _ := ret[0].([]types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithFlappingStatus indicates an expected call of GetAlarmsWithFlappingStatus.
func (mr *MockAdapterMockRecorder) GetAlarmsWithFlappingStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithFlappingStatus", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithFlappingStatus), arg0)
}

// GetAlarmsWithSnoozeMark mocks base method.
func (m *MockAdapter) GetAlarmsWithSnoozeMark(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithSnoozeMark", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithSnoozeMark indicates an expected call of GetAlarmsWithSnoozeMark.
func (mr *MockAdapterMockRecorder) GetAlarmsWithSnoozeMark(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithSnoozeMark", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithSnoozeMark), arg0)
}

// GetAlarmsWithoutTicketByComponent mocks base method.
func (m *MockAdapter) GetAlarmsWithoutTicketByComponent(arg0 context.Context, arg1 string) ([]types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmsWithoutTicketByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlarmsWithoutTicketByComponent indicates an expected call of GetAlarmsWithoutTicketByComponent.
func (mr *MockAdapterMockRecorder) GetAlarmsWithoutTicketByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmsWithoutTicketByComponent", reflect.TypeOf((*MockAdapter)(nil).GetAlarmsWithoutTicketByComponent), arg0, arg1)
}

// GetAllOpenedResourceAlarmsByComponent mocks base method.
func (m *MockAdapter) GetAllOpenedResourceAlarmsByComponent(arg0 context.Context, arg1 string) ([]types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOpenedResourceAlarmsByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOpenedResourceAlarmsByComponent indicates an expected call of GetAllOpenedResourceAlarmsByComponent.
func (mr *MockAdapterMockRecorder) GetAllOpenedResourceAlarmsByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOpenedResourceAlarmsByComponent", reflect.TypeOf((*MockAdapter)(nil).GetAllOpenedResourceAlarmsByComponent), arg0, arg1)
}

// GetCountOpenedAlarmsByIDs mocks base method.
func (m *MockAdapter) GetCountOpenedAlarmsByIDs(arg0 context.Context, arg1 []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountOpenedAlarmsByIDs", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountOpenedAlarmsByIDs indicates an expected call of GetCountOpenedAlarmsByIDs.
func (mr *MockAdapterMockRecorder) GetCountOpenedAlarmsByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountOpenedAlarmsByIDs", reflect.TypeOf((*MockAdapter)(nil).GetCountOpenedAlarmsByIDs), arg0, arg1)
}

// GetLastAlarm mocks base method.
func (m *MockAdapter) GetLastAlarm(arg0 context.Context, arg1, arg2, arg3 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAlarm", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAlarm indicates an expected call of GetLastAlarm.
func (mr *MockAdapterMockRecorder) GetLastAlarm(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAlarm", reflect.TypeOf((*MockAdapter)(nil).GetLastAlarm), arg0, arg1, arg2, arg3)
}

// GetLastAlarmByEntityID mocks base method.
func (m *MockAdapter) GetLastAlarmByEntityID(arg0 context.Context, arg1 string) (*types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAlarmByEntityID", arg0, arg1)
	ret0, _ := ret[0].(*types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAlarmByEntityID indicates an expected call of GetLastAlarmByEntityID.
func (mr *MockAdapterMockRecorder) GetLastAlarmByEntityID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAlarmByEntityID", reflect.TypeOf((*MockAdapter)(nil).GetLastAlarmByEntityID), arg0, arg1)
}

// GetLastAlarmWithEntity mocks base method.
func (m *MockAdapter) GetLastAlarmWithEntity(arg0 context.Context, arg1, arg2, arg3 string) (types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAlarmWithEntity", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAlarmWithEntity indicates an expected call of GetLastAlarmWithEntity.
func (mr *MockAdapterMockRecorder) GetLastAlarmWithEntity(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAlarmWithEntity", reflect.TypeOf((*MockAdapter)(nil).GetLastAlarmWithEntity), arg0, arg1, arg2, arg3)
}

// GetOpenedAlarm mocks base method.
func (m *MockAdapter) GetOpenedAlarm(arg0 context.Context, arg1, arg2, arg3 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarm", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarm indicates an expected call of GetOpenedAlarm.
func (mr *MockAdapterMockRecorder) GetOpenedAlarm(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarm", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarm), arg0, arg1, arg2, arg3)
}

// GetOpenedAlarmByAlarmId mocks base method.
func (m *MockAdapter) GetOpenedAlarmByAlarmId(arg0 context.Context, arg1 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmByAlarmId", arg0, arg1)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmByAlarmId indicates an expected call of GetOpenedAlarmByAlarmId.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmByAlarmId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmByAlarmId", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmByAlarmId), arg0, arg1)
}

// GetOpenedAlarmsByAlarmIDs mocks base method.
func (m *MockAdapter) GetOpenedAlarmsByAlarmIDs(arg0 context.Context, arg1 []string, arg2 *[]types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByAlarmIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsByAlarmIDs indicates an expected call of GetOpenedAlarmsByAlarmIDs.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByAlarmIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByAlarmIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByAlarmIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsByConnectorIdleRules mocks base method.
func (m *MockAdapter) GetOpenedAlarmsByConnectorIdleRules(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByConnectorIdleRules", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsByConnectorIdleRules indicates an expected call of GetOpenedAlarmsByConnectorIdleRules.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByConnectorIdleRules(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByConnectorIdleRules", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByConnectorIdleRules), arg0)
}

// GetOpenedAlarmsByIDs mocks base method.
func (m *MockAdapter) GetOpenedAlarmsByIDs(arg0 context.Context, arg1 []string, arg2 *[]types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsByIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsByIDs indicates an expected call of GetOpenedAlarmsByIDs.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsByIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsByIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsByIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithEntity mocks base method.
func (m *MockAdapter) GetOpenedAlarmsWithEntity(arg0 context.Context) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntity", arg0)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsWithEntity indicates an expected call of GetOpenedAlarmsWithEntity.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntity", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntity), arg0)
}

// GetOpenedAlarmsWithEntityAfter mocks base method.
func (m *MockAdapter) GetOpenedAlarmsWithEntityAfter(arg0 context.Context, arg1 types.CpsTime) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntityAfter", arg0, arg1)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsWithEntityAfter indicates an expected call of GetOpenedAlarmsWithEntityAfter.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntityAfter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntityAfter", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntityAfter), arg0, arg1)
}

// GetOpenedAlarmsWithEntityByAlarmIDs mocks base method.
func (m *MockAdapter) GetOpenedAlarmsWithEntityByAlarmIDs(arg0 context.Context, arg1 []string, arg2 *[]types.AlarmWithEntity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntityByAlarmIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsWithEntityByAlarmIDs indicates an expected call of GetOpenedAlarmsWithEntityByAlarmIDs.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntityByAlarmIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntityByAlarmIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntityByAlarmIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithEntityByIDs mocks base method.
func (m *MockAdapter) GetOpenedAlarmsWithEntityByIDs(arg0 context.Context, arg1 []string, arg2 *[]types.AlarmWithEntity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithEntityByIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOpenedAlarmsWithEntityByIDs indicates an expected call of GetOpenedAlarmsWithEntityByIDs.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithEntityByIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithEntityByIDs", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithEntityByIDs), arg0, arg1, arg2)
}

// GetOpenedAlarmsWithLastDatesBefore mocks base method.
func (m *MockAdapter) GetOpenedAlarmsWithLastDatesBefore(arg0 context.Context, arg1 types.CpsTime) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedAlarmsWithLastDatesBefore", arg0, arg1)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedAlarmsWithLastDatesBefore indicates an expected call of GetOpenedAlarmsWithLastDatesBefore.
func (mr *MockAdapterMockRecorder) GetOpenedAlarmsWithLastDatesBefore(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedAlarmsWithLastDatesBefore", reflect.TypeOf((*MockAdapter)(nil).GetOpenedAlarmsWithLastDatesBefore), arg0, arg1)
}

// GetOpenedMetaAlarm mocks base method.
func (m *MockAdapter) GetOpenedMetaAlarm(arg0 context.Context, arg1, arg2 string) (types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedMetaAlarm", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedMetaAlarm indicates an expected call of GetOpenedMetaAlarm.
func (mr *MockAdapterMockRecorder) GetOpenedMetaAlarm(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedMetaAlarm", reflect.TypeOf((*MockAdapter)(nil).GetOpenedMetaAlarm), arg0, arg1, arg2)
}

// GetOpenedMetaAlarmWithEntity mocks base method.
func (m *MockAdapter) GetOpenedMetaAlarmWithEntity(arg0 context.Context, arg1, arg2 string) (types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenedMetaAlarmWithEntity", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenedMetaAlarmWithEntity indicates an expected call of GetOpenedMetaAlarmWithEntity.
func (mr *MockAdapterMockRecorder) GetOpenedMetaAlarmWithEntity(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenedMetaAlarmWithEntity", reflect.TypeOf((*MockAdapter)(nil).GetOpenedMetaAlarmWithEntity), arg0, arg1, arg2)
}

// GetUnacknowledgedAlarmsByComponent mocks base method.
func (m *MockAdapter) GetUnacknowledgedAlarmsByComponent(arg0 context.Context, arg1 string) ([]types.AlarmWithEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnacknowledgedAlarmsByComponent", arg0, arg1)
	ret0, _ := ret[0].([]types.AlarmWithEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnacknowledgedAlarmsByComponent indicates an expected call of GetUnacknowledgedAlarmsByComponent.
func (mr *MockAdapterMockRecorder) GetUnacknowledgedAlarmsByComponent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnacknowledgedAlarmsByComponent", reflect.TypeOf((*MockAdapter)(nil).GetUnacknowledgedAlarmsByComponent), arg0, arg1)
}

// GetWorstAlarmState mocks base method.
func (m *MockAdapter) GetWorstAlarmState(arg0 context.Context, arg1 []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorstAlarmState", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorstAlarmState indicates an expected call of GetWorstAlarmState.
func (mr *MockAdapterMockRecorder) GetWorstAlarmState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorstAlarmState", reflect.TypeOf((*MockAdapter)(nil).GetWorstAlarmState), arg0, arg1)
}

// Insert mocks base method.
func (m *MockAdapter) Insert(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockAdapterMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAdapter)(nil).Insert), arg0, arg1)
}

// MassPartialUpdateOpen mocks base method.
func (m *MockAdapter) MassPartialUpdateOpen(arg0 context.Context, arg1 *types.Alarm, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassPartialUpdateOpen", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassPartialUpdateOpen indicates an expected call of MassPartialUpdateOpen.
func (mr *MockAdapterMockRecorder) MassPartialUpdateOpen(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassPartialUpdateOpen", reflect.TypeOf((*MockAdapter)(nil).MassPartialUpdateOpen), arg0, arg1, arg2)
}

// PartialMassUpdateOpen mocks base method.
func (m *MockAdapter) PartialMassUpdateOpen(arg0 context.Context, arg1 []types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialMassUpdateOpen", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PartialMassUpdateOpen indicates an expected call of PartialMassUpdateOpen.
func (mr *MockAdapterMockRecorder) PartialMassUpdateOpen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialMassUpdateOpen", reflect.TypeOf((*MockAdapter)(nil).PartialMassUpdateOpen), arg0, arg1)
}

// PartialUpdateOpen mocks base method.
func (m *MockAdapter) PartialUpdateOpen(arg0 context.Context, arg1 *types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialUpdateOpen", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PartialUpdateOpen indicates an expected call of PartialUpdateOpen.
func (mr *MockAdapterMockRecorder) PartialUpdateOpen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialUpdateOpen", reflect.TypeOf((*MockAdapter)(nil).PartialUpdateOpen), arg0, arg1)
}

// Update mocks base method.
func (m *MockAdapter) Update(arg0 context.Context, arg1 types.Alarm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAdapterMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAdapter)(nil).Update), arg0, arg1)
}

// UpdateLastEventDate mocks base method.
func (m *MockAdapter) UpdateLastEventDate(arg0 context.Context, arg1 []string, arg2 types.CpsTime) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastEventDate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastEventDate indicates an expected call of UpdateLastEventDate.
func (mr *MockAdapterMockRecorder) UpdateLastEventDate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastEventDate", reflect.TypeOf((*MockAdapter)(nil).UpdateLastEventDate), arg0, arg1, arg2)
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

// ResolveCancels mocks base method.
func (m *MockService) ResolveCancels(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveCancels", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveCancels indicates an expected call of ResolveCancels.
func (mr *MockServiceMockRecorder) ResolveCancels(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveCancels", reflect.TypeOf((*MockService)(nil).ResolveCancels), arg0, arg1)
}

// ResolveClosed mocks base method.
func (m *MockService) ResolveClosed(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveClosed", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveClosed indicates an expected call of ResolveClosed.
func (mr *MockServiceMockRecorder) ResolveClosed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveClosed", reflect.TypeOf((*MockService)(nil).ResolveClosed), arg0)
}

// ResolveDone mocks base method.
func (m *MockService) ResolveDone(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveDone", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveDone indicates an expected call of ResolveDone.
func (mr *MockServiceMockRecorder) ResolveDone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveDone", reflect.TypeOf((*MockService)(nil).ResolveDone), arg0)
}

// ResolveSnoozes mocks base method.
func (m *MockService) ResolveSnoozes(arg0 context.Context, arg1 config.AlarmConfig) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveSnoozes", arg0, arg1)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveSnoozes indicates an expected call of ResolveSnoozes.
func (mr *MockServiceMockRecorder) ResolveSnoozes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveSnoozes", reflect.TypeOf((*MockService)(nil).ResolveSnoozes), arg0, arg1)
}

// UpdateFlappingAlarms mocks base method.
func (m *MockService) UpdateFlappingAlarms(arg0 context.Context) ([]types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFlappingAlarms", arg0)
	ret0, _ := ret[0].([]types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFlappingAlarms indicates an expected call of UpdateFlappingAlarms.
func (mr *MockServiceMockRecorder) UpdateFlappingAlarms(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlappingAlarms", reflect.TypeOf((*MockService)(nil).UpdateFlappingAlarms), arg0)
}

// MockEventProcessor is a mock of EventProcessor interface.
type MockEventProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockEventProcessorMockRecorder
}

// MockEventProcessorMockRecorder is the mock recorder for MockEventProcessor.
type MockEventProcessorMockRecorder struct {
	mock *MockEventProcessor
}

// NewMockEventProcessor creates a new mock instance.
func NewMockEventProcessor(ctrl *gomock.Controller) *MockEventProcessor {
	mock := &MockEventProcessor{ctrl: ctrl}
	mock.recorder = &MockEventProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventProcessor) EXPECT() *MockEventProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockEventProcessor) Process(arg0 context.Context, arg1 *types.Event) (types.AlarmChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1)
	ret0, _ := ret[0].(types.AlarmChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockEventProcessorMockRecorder) Process(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockEventProcessor)(nil).Process), arg0, arg1)
}

// MockActivationService is a mock of ActivationService interface.
type MockActivationService struct {
	ctrl     *gomock.Controller
	recorder *MockActivationServiceMockRecorder
}

// MockActivationServiceMockRecorder is the mock recorder for MockActivationService.
type MockActivationServiceMockRecorder struct {
	mock *MockActivationService
}

// NewMockActivationService creates a new mock instance.
func NewMockActivationService(ctrl *gomock.Controller) *MockActivationService {
	mock := &MockActivationService{ctrl: ctrl}
	mock.recorder = &MockActivationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActivationService) EXPECT() *MockActivationServiceMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockActivationService) Process(arg0 *types.Alarm) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockActivationServiceMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockActivationService)(nil).Process), arg0)
}

// MockMetaAlarmEventProcessor is a mock of MetaAlarmEventProcessor interface.
type MockMetaAlarmEventProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockMetaAlarmEventProcessorMockRecorder
}

// MockMetaAlarmEventProcessorMockRecorder is the mock recorder for MockMetaAlarmEventProcessor.
type MockMetaAlarmEventProcessorMockRecorder struct {
	mock *MockMetaAlarmEventProcessor
}

// NewMockMetaAlarmEventProcessor creates a new mock instance.
func NewMockMetaAlarmEventProcessor(ctrl *gomock.Controller) *MockMetaAlarmEventProcessor {
	mock := &MockMetaAlarmEventProcessor{ctrl: ctrl}
	mock.recorder = &MockMetaAlarmEventProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetaAlarmEventProcessor) EXPECT() *MockMetaAlarmEventProcessorMockRecorder {
	return m.recorder
}

// CreateMetaAlarm mocks base method.
func (m *MockMetaAlarmEventProcessor) CreateMetaAlarm(arg0 context.Context, arg1 types.Event) (*types.Alarm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMetaAlarm", arg0, arg1)
	ret0, _ := ret[0].(*types.Alarm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMetaAlarm indicates an expected call of CreateMetaAlarm.
func (mr *MockMetaAlarmEventProcessorMockRecorder) CreateMetaAlarm(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMetaAlarm", reflect.TypeOf((*MockMetaAlarmEventProcessor)(nil).CreateMetaAlarm), arg0, arg1)
}

// Process mocks base method.
func (m *MockMetaAlarmEventProcessor) Process(arg0 context.Context, arg1 types.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockMetaAlarmEventProcessorMockRecorder) Process(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockMetaAlarmEventProcessor)(nil).Process), arg0, arg1)
}

// ProcessAckResources mocks base method.
func (m *MockMetaAlarmEventProcessor) ProcessAckResources(arg0 context.Context, arg1 types.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessAckResources", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessAckResources indicates an expected call of ProcessAckResources.
func (mr *MockMetaAlarmEventProcessorMockRecorder) ProcessAckResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessAckResources", reflect.TypeOf((*MockMetaAlarmEventProcessor)(nil).ProcessAckResources), arg0, arg1)
}

// ProcessAxeRpc mocks base method.
func (m *MockMetaAlarmEventProcessor) ProcessAxeRpc(arg0 context.Context, arg1 types.RPCAxeEvent, arg2 types.RPCAxeResultEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessAxeRpc", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessAxeRpc indicates an expected call of ProcessAxeRpc.
func (mr *MockMetaAlarmEventProcessorMockRecorder) ProcessAxeRpc(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessAxeRpc", reflect.TypeOf((*MockMetaAlarmEventProcessor)(nil).ProcessAxeRpc), arg0, arg1, arg2)
}

// ProcessWebhookRpc mocks base method.
func (m *MockMetaAlarmEventProcessor) ProcessWebhookRpc(arg0 context.Context, arg1 types.RPCWebhookEvent, arg2 string, arg3 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessWebhookRpc", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessWebhookRpc indicates an expected call of ProcessWebhookRpc.
func (mr *MockMetaAlarmEventProcessorMockRecorder) ProcessWebhookRpc(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessWebhookRpc", reflect.TypeOf((*MockMetaAlarmEventProcessor)(nil).ProcessWebhookRpc), arg0, arg1, arg2, arg3)
}
