package externaldata_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/mock/gomock"
)

func TestSyncMongoCollections_GivenCollections_ShouldAdd(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noTag := externaldata.ColumnTagNoTag

	collNames := []string{
		"test_coll_1",
		"test_coll_2",
		"test_coll_3",
	}
	collNamesToIgnore := []string{
		"test_coll_2",
	}
	tablesToCreate := []externaldata.Table{
		{
			Name: "test_coll_1",
			ColumnConfigs: []externaldata.ColumnConfig{
				{
					Name: "test_field_1",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
				{
					Name: "test_field_2",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
			},
		},
		{
			Name: "test_coll_3",
			ColumnConfigs: []externaldata.ColumnConfig{
				{
					Name: "test_field_3",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
				{
					Name: "test_field_4",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
				{
					Name: "test_field_5",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
			},
		},
	}
	mockClient := mock_mongo.NewMockDbClient(ctrl)
	mockExdataTableCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockClient.EXPECT().Collection(gomock.Eq(mongo.ExternalDataTableCollection)).Return(mockExdataTableCollection)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, nil), nil)
	mockExdataTableCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongodriver.UpdateResult{}, nil)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, collNamesToIgnore), nil)
	mockExdataTableCollection.EXPECT().InsertMany(gomock.Any(), gomock.Any()).Return(nil, nil).Do(checkInsertedTables(t, tablesToCreate))
	newMockCollections(ctrl, mockClient, tablesToCreate, "test")
	err := externaldata.SyncMongoCollections(ctx, mockClient, collNames, nil, zerolog.Nop())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncMongoCollections_GivenEmptyCollections_ShouldNotCreateExdata(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noTag := externaldata.ColumnTagNoTag

	collNames := []string{
		"test_coll_1",
		"test_coll_2",
	}
	tablesToIgnore := []externaldata.Table{
		{
			Name:          "test_coll_1",
			ColumnConfigs: []externaldata.ColumnConfig{},
		},
		{
			Name: "test_coll_2",
			ColumnConfigs: []externaldata.ColumnConfig{
				{
					Name: "test_field_3",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
				{
					Name: "test_field_4",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
				{
					Name: "test_field_5",
					Type: externaldata.ColumnTypeString,
					Tag:  &noTag,
				},
			},
		},
	}
	mockClient := mock_mongo.NewMockDbClient(ctrl)
	mockExdataTableCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockClient.EXPECT().Collection(gomock.Eq(mongo.ExternalDataTableCollection)).Return(mockExdataTableCollection)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, nil), nil)
	mockExdataTableCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongodriver.UpdateResult{}, nil)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, nil), nil)
	newMockCollections(ctrl, mockClient, tablesToIgnore, 1)
	err := externaldata.SyncMongoCollections(ctx, mockClient, collNames, nil, zerolog.Nop())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncMongoCollections_GivenMissingCollections_ShouldDeleteUnlinked(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collNames := []string{
		"test_coll_1",
	}
	collNamesToIgnore := []string{
		"test_coll_1",
	}
	collNamesToDelete := []string{
		"test_coll_2",
		"test_coll_3",
	}
	collNamesToBlockDelete := []string{
		"test_coll_4",
		"test_coll_5",
		"test_coll_6",
	}
	missingCollNames := []string{
		"test_coll_2",
		"test_coll_3",
		"test_coll_4",
		"test_coll_5",
		"test_coll_6",
	}
	refCollNames := []string{
		"test_ref_coll_1",
		"test_ref_coll_2",
	}
	mockClient := mock_mongo.NewMockDbClient(ctrl)
	mockRefCollection1 := mock_mongo.NewMockDbCollection(ctrl)
	mockRefCollection1.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(newMockCursorRef(ctrl, collNamesToBlockDelete[:1]), nil)
	mockRefCollection2 := mock_mongo.NewMockDbCollection(ctrl)
	mockRefCollection2.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(newMockCursorRef(ctrl, collNamesToBlockDelete[1:2]), nil)
	mockWidgetCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockWidgetCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(newMockCursorRef(ctrl, collNamesToBlockDelete[2:]), nil)
	mockClient.EXPECT().Collection(gomock.Eq(refCollNames[0])).Return(mockRefCollection1)
	mockClient.EXPECT().Collection(gomock.Eq(refCollNames[1])).Return(mockRefCollection2)
	mockClient.EXPECT().Collection(gomock.Eq(mongo.WidgetMongoCollection)).Return(mockWidgetCollection)
	mockExdataTableCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockClient.EXPECT().Collection(gomock.Eq(mongo.ExternalDataTableCollection)).Return(mockExdataTableCollection)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, missingCollNames), nil)
	mockExdataTableCollection.EXPECT().DeleteMany(gomock.Any(), gomock.Eq(bson.M{"_id": bson.M{"$in": collNamesToDelete}}))
	mockExdataTableCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Eq(bson.M{"_id": bson.M{"$in": collNamesToBlockDelete}}), gomock.Any()).Return(&mongodriver.UpdateResult{}, nil)
	mockExdataTableCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongodriver.UpdateResult{}, nil)
	mockExdataTableCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(newMockCursorTables(ctrl, collNamesToIgnore), nil)
	err := externaldata.SyncMongoCollections(ctx, mockClient, collNames, refCollNames, zerolog.Nop())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newMockCursorTables(ctrl *gomock.Controller, collNames []string) mongo.Cursor {
	mockTableCursor := mock_mongo.NewMockCursor(ctrl)
	mockTableCursor.EXPECT().Next(gomock.Any()).Return(true).Times(len(collNames))
	for i := range collNames {
		name := collNames[i]
		mockTableCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v *externaldata.Table) error {
			v.ID = name
			v.Name = name

			return nil
		})
	}

	mockTableCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockTableCursor.EXPECT().Err().Return(nil)
	mockTableCursor.EXPECT().Close(gomock.Any()).Return(nil)

	return mockTableCursor
}

func newMockCollections(
	ctrl *gomock.Controller,
	mockClient *mock_mongo.MockDbClient,
	tables []externaldata.Table,
	fieldVal any,
) {
	count := 5
	for _, table := range tables {
		doc := make(bson.D, len(table.ColumnConfigs)+1)
		doc[0] = bson.E{Key: externaldata.IDColumnName, Value: "test"}
		for i, s := range table.ColumnConfigs {
			doc[i+1] = bson.E{Key: s.Name, Value: fieldVal}
		}

		mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
		mockClient.EXPECT().Collection(gomock.Eq(table.Name)).Return(mockDbCollection)
		mockCursor := mock_mongo.NewMockCursor(ctrl)
		mockCursor.EXPECT().Next(gomock.Any()).Return(true).Times(count)
		mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *bson.D) {
			*v = doc
		}).Times(count)
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
		mockCursor.EXPECT().Err().Return(nil)
		mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
		mockDbCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCursor, nil)
	}
}

func checkInsertedTables(t *testing.T, expected []externaldata.Table) func(context.Context, []any, ...*options.InsertManyOptions) {
	return func(_ context.Context, docs []any, _ ...*options.InsertManyOptions) {
		res := make([]externaldata.Table, len(docs))
		for i, doc := range docs {
			if table, ok := doc.(externaldata.Table); ok {
				res[i] = externaldata.Table{
					Name:          table.Name,
					ColumnConfigs: table.ColumnConfigs,
				}
			} else {
				t.Fatalf("unknown doc: %T %+v", doc, doc)
			}
		}

		cmp := func(l, f externaldata.Table) int {
			return strings.Compare(l.Name, f.Name)
		}

		slices.SortFunc(res, cmp)
		slices.SortFunc(expected, cmp)
		if diff := pretty.Compare(expected, res); diff != "" {
			t.Fatal("unexpected result: " + diff)
		}
	}
}

func newMockCursorRef(ctrl *gomock.Controller, collIDs []string) mongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockCursor.EXPECT().Next(gomock.Any()).Return(true).Times(len(collIDs))
	for i := range collIDs {
		id := collIDs[i]
		mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v *struct {
			ID string `bson:"_id"`
		}) error {
			v.ID = id

			return nil
		})
	}

	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockCursor.EXPECT().Err().Return(nil)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)

	return mockCursor
}
