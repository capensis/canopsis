// Code generated by MockGen. DO NOT EDIT.
// Source: git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo (interfaces: DbCollection,DbClient,SingleResultHelper,Cursor)

// Package mock_mongo is a generated GoMock package.
package mock_mongo

import (
	context "context"
	mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	gomock "github.com/golang/mock/gomock"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo0 "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	reflect "reflect"
	time "time"
)

// MockDbCollection is a mock of DbCollection interface
type MockDbCollection struct {
	ctrl     *gomock.Controller
	recorder *MockDbCollectionMockRecorder
}

// MockDbCollectionMockRecorder is the mock recorder for MockDbCollection
type MockDbCollectionMockRecorder struct {
	mock *MockDbCollection
}

// NewMockDbCollection creates a new mock instance
func NewMockDbCollection(ctrl *gomock.Controller) *MockDbCollection {
	mock := &MockDbCollection{ctrl: ctrl}
	mock.recorder = &MockDbCollectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDbCollection) EXPECT() *MockDbCollectionMockRecorder {
	return m.recorder
}

// Aggregate mocks base method
func (m *MockDbCollection) Aggregate(arg0 context.Context, arg1 interface{}, arg2 ...*options.AggregateOptions) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Aggregate", varargs...)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Aggregate indicates an expected call of Aggregate
func (mr *MockDbCollectionMockRecorder) Aggregate(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Aggregate", reflect.TypeOf((*MockDbCollection)(nil).Aggregate), varargs...)
}

// BulkWrite mocks base method
func (m *MockDbCollection) BulkWrite(arg0 context.Context, arg1 []mongo0.WriteModel, arg2 ...*options.BulkWriteOptions) (*mongo0.BulkWriteResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BulkWrite", varargs...)
	ret0, _ := ret[0].(*mongo0.BulkWriteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkWrite indicates an expected call of BulkWrite
func (mr *MockDbCollectionMockRecorder) BulkWrite(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkWrite", reflect.TypeOf((*MockDbCollection)(nil).BulkWrite), varargs...)
}

// CountDocuments mocks base method
func (m *MockDbCollection) CountDocuments(arg0 context.Context, arg1 interface{}, arg2 ...*options.CountOptions) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CountDocuments", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountDocuments indicates an expected call of CountDocuments
func (mr *MockDbCollectionMockRecorder) CountDocuments(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountDocuments", reflect.TypeOf((*MockDbCollection)(nil).CountDocuments), varargs...)
}

// DeleteMany mocks base method
func (m *MockDbCollection) DeleteMany(arg0 context.Context, arg1 interface{}, arg2 ...*options.DeleteOptions) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteMany", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMany indicates an expected call of DeleteMany
func (mr *MockDbCollectionMockRecorder) DeleteMany(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMany", reflect.TypeOf((*MockDbCollection)(nil).DeleteMany), varargs...)
}

// DeleteOne mocks base method
func (m *MockDbCollection) DeleteOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.DeleteOptions) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteOne", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteOne indicates an expected call of DeleteOne
func (mr *MockDbCollectionMockRecorder) DeleteOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOne", reflect.TypeOf((*MockDbCollection)(nil).DeleteOne), varargs...)
}

// Distinct mocks base method
func (m *MockDbCollection) Distinct(arg0 context.Context, arg1 string, arg2 interface{}, arg3 ...*options.DistinctOptions) ([]interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Distinct", varargs...)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Distinct indicates an expected call of Distinct
func (mr *MockDbCollectionMockRecorder) Distinct(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Distinct", reflect.TypeOf((*MockDbCollection)(nil).Distinct), varargs...)
}

// Drop mocks base method
func (m *MockDbCollection) Drop(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Drop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Drop indicates an expected call of Drop
func (mr *MockDbCollectionMockRecorder) Drop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Drop", reflect.TypeOf((*MockDbCollection)(nil).Drop), arg0)
}

// Find mocks base method
func (m *MockDbCollection) Find(arg0 context.Context, arg1 interface{}, arg2 ...*options.FindOptions) (mongo.Cursor, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].(mongo.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockDbCollectionMockRecorder) Find(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockDbCollection)(nil).Find), varargs...)
}

// FindOne mocks base method
func (m *MockDbCollection) FindOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.FindOneOptions) mongo.SingleResultHelper {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOne", varargs...)
	ret0, _ := ret[0].(mongo.SingleResultHelper)
	return ret0
}

// FindOne indicates an expected call of FindOne
func (mr *MockDbCollectionMockRecorder) FindOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockDbCollection)(nil).FindOne), varargs...)
}

// FindOneAndDelete mocks base method
func (m *MockDbCollection) FindOneAndDelete(arg0 context.Context, arg1 interface{}, arg2 ...*options.FindOneAndDeleteOptions) mongo.SingleResultHelper {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOneAndDelete", varargs...)
	ret0, _ := ret[0].(mongo.SingleResultHelper)
	return ret0
}

// FindOneAndDelete indicates an expected call of FindOneAndDelete
func (mr *MockDbCollectionMockRecorder) FindOneAndDelete(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneAndDelete", reflect.TypeOf((*MockDbCollection)(nil).FindOneAndDelete), varargs...)
}

// FindOneAndReplace mocks base method
func (m *MockDbCollection) FindOneAndReplace(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.FindOneAndReplaceOptions) mongo.SingleResultHelper {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOneAndReplace", varargs...)
	ret0, _ := ret[0].(mongo.SingleResultHelper)
	return ret0
}

// FindOneAndReplace indicates an expected call of FindOneAndReplace
func (mr *MockDbCollectionMockRecorder) FindOneAndReplace(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneAndReplace", reflect.TypeOf((*MockDbCollection)(nil).FindOneAndReplace), varargs...)
}

// FindOneAndUpdate mocks base method
func (m *MockDbCollection) FindOneAndUpdate(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.FindOneAndUpdateOptions) mongo.SingleResultHelper {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOneAndUpdate", varargs...)
	ret0, _ := ret[0].(mongo.SingleResultHelper)
	return ret0
}

// FindOneAndUpdate indicates an expected call of FindOneAndUpdate
func (mr *MockDbCollectionMockRecorder) FindOneAndUpdate(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneAndUpdate", reflect.TypeOf((*MockDbCollection)(nil).FindOneAndUpdate), varargs...)
}

// Indexes mocks base method
func (m *MockDbCollection) Indexes() mongo0.IndexView {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Indexes")
	ret0, _ := ret[0].(mongo0.IndexView)
	return ret0
}

// Indexes indicates an expected call of Indexes
func (mr *MockDbCollectionMockRecorder) Indexes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Indexes", reflect.TypeOf((*MockDbCollection)(nil).Indexes))
}

// InsertMany mocks base method
func (m *MockDbCollection) InsertMany(arg0 context.Context, arg1 []interface{}, arg2 ...*options.InsertManyOptions) ([]interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertMany", varargs...)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertMany indicates an expected call of InsertMany
func (mr *MockDbCollectionMockRecorder) InsertMany(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMany", reflect.TypeOf((*MockDbCollection)(nil).InsertMany), varargs...)
}

// InsertOne mocks base method
func (m *MockDbCollection) InsertOne(arg0 context.Context, arg1 interface{}, arg2 ...*options.InsertOneOptions) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertOne", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertOne indicates an expected call of InsertOne
func (mr *MockDbCollectionMockRecorder) InsertOne(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOne", reflect.TypeOf((*MockDbCollection)(nil).InsertOne), varargs...)
}

// ReplaceOne mocks base method
func (m *MockDbCollection) ReplaceOne(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.ReplaceOptions) (*mongo0.UpdateResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReplaceOne", varargs...)
	ret0, _ := ret[0].(*mongo0.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceOne indicates an expected call of ReplaceOne
func (mr *MockDbCollectionMockRecorder) ReplaceOne(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceOne", reflect.TypeOf((*MockDbCollection)(nil).ReplaceOne), varargs...)
}

// UpdateMany mocks base method
func (m *MockDbCollection) UpdateMany(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.UpdateOptions) (*mongo0.UpdateResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateMany", varargs...)
	ret0, _ := ret[0].(*mongo0.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMany indicates an expected call of UpdateMany
func (mr *MockDbCollectionMockRecorder) UpdateMany(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMany", reflect.TypeOf((*MockDbCollection)(nil).UpdateMany), varargs...)
}

// UpdateOne mocks base method
func (m *MockDbCollection) UpdateOne(arg0 context.Context, arg1, arg2 interface{}, arg3 ...*options.UpdateOptions) (*mongo0.UpdateResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateOne", varargs...)
	ret0, _ := ret[0].(*mongo0.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOne indicates an expected call of UpdateOne
func (mr *MockDbCollectionMockRecorder) UpdateOne(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOne", reflect.TypeOf((*MockDbCollection)(nil).UpdateOne), varargs...)
}

// MockDbClient is a mock of DbClient interface
type MockDbClient struct {
	ctrl     *gomock.Controller
	recorder *MockDbClientMockRecorder
}

// MockDbClientMockRecorder is the mock recorder for MockDbClient
type MockDbClientMockRecorder struct {
	mock *MockDbClient
}

// NewMockDbClient creates a new mock instance
func NewMockDbClient(ctrl *gomock.Controller) *MockDbClient {
	mock := &MockDbClient{ctrl: ctrl}
	mock.recorder = &MockDbClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDbClient) EXPECT() *MockDbClientMockRecorder {
	return m.recorder
}

// Collection mocks base method
func (m *MockDbClient) Collection(arg0 string) mongo.DbCollection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collection", arg0)
	ret0, _ := ret[0].(mongo.DbCollection)
	return ret0
}

// Collection indicates an expected call of Collection
func (mr *MockDbClientMockRecorder) Collection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collection", reflect.TypeOf((*MockDbClient)(nil).Collection), arg0)
}

// Disconnect mocks base method
func (m *MockDbClient) Disconnect(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect
func (mr *MockDbClientMockRecorder) Disconnect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockDbClient)(nil).Disconnect), arg0)
}

// SetRetry mocks base method
func (m *MockDbClient) SetRetry(arg0 int, arg1 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRetry", arg0, arg1)
}

// SetRetry indicates an expected call of SetRetry
func (mr *MockDbClientMockRecorder) SetRetry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRetry", reflect.TypeOf((*MockDbClient)(nil).SetRetry), arg0, arg1)
}

// MockSingleResultHelper is a mock of SingleResultHelper interface
type MockSingleResultHelper struct {
	ctrl     *gomock.Controller
	recorder *MockSingleResultHelperMockRecorder
}

// MockSingleResultHelperMockRecorder is the mock recorder for MockSingleResultHelper
type MockSingleResultHelperMockRecorder struct {
	mock *MockSingleResultHelper
}

// NewMockSingleResultHelper creates a new mock instance
func NewMockSingleResultHelper(ctrl *gomock.Controller) *MockSingleResultHelper {
	mock := &MockSingleResultHelper{ctrl: ctrl}
	mock.recorder = &MockSingleResultHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSingleResultHelper) EXPECT() *MockSingleResultHelperMockRecorder {
	return m.recorder
}

// Decode mocks base method
func (m *MockSingleResultHelper) Decode(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode
func (mr *MockSingleResultHelperMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockSingleResultHelper)(nil).Decode), arg0)
}

// DecodeBytes mocks base method
func (m *MockSingleResultHelper) DecodeBytes() (bson.Raw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeBytes")
	ret0, _ := ret[0].(bson.Raw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecodeBytes indicates an expected call of DecodeBytes
func (mr *MockSingleResultHelperMockRecorder) DecodeBytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeBytes", reflect.TypeOf((*MockSingleResultHelper)(nil).DecodeBytes))
}

// Err mocks base method
func (m *MockSingleResultHelper) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockSingleResultHelperMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockSingleResultHelper)(nil).Err))
}

// MockCursor is a mock of Cursor interface
type MockCursor struct {
	ctrl     *gomock.Controller
	recorder *MockCursorMockRecorder
}

// MockCursorMockRecorder is the mock recorder for MockCursor
type MockCursorMockRecorder struct {
	mock *MockCursor
}

// NewMockCursor creates a new mock instance
func NewMockCursor(ctrl *gomock.Controller) *MockCursor {
	mock := &MockCursor{ctrl: ctrl}
	mock.recorder = &MockCursorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCursor) EXPECT() *MockCursorMockRecorder {
	return m.recorder
}

// All mocks base method
func (m *MockCursor) All(arg0 context.Context, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// All indicates an expected call of All
func (mr *MockCursorMockRecorder) All(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockCursor)(nil).All), arg0, arg1)
}

// Close mocks base method
func (m *MockCursor) Close(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCursorMockRecorder) Close(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCursor)(nil).Close), arg0)
}

// Decode mocks base method
func (m *MockCursor) Decode(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode
func (mr *MockCursorMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockCursor)(nil).Decode), arg0)
}

// Err mocks base method
func (m *MockCursor) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockCursorMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockCursor)(nil).Err))
}

// Next mocks base method
func (m *MockCursor) Next(arg0 context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next
func (mr *MockCursorMockRecorder) Next(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockCursor)(nil).Next), arg0)
}
