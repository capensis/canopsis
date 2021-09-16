package mongoadapter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libmodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sort"
	"testing"
)

func TestAdapter_LoadPolicy_GivenRole_ShouldAddCRUDPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObjCursor := createMockCursor(ctrl, []libmodel.Rbac{
		{
			ID:         "testobj1",
			Name:       "testobjname1",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
		{
			ID:         "testobj2",
			Name:       "testobjname2",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
		{
			ID:         "testobj3",
			Name:       "testobjname3",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
		{
			ID:         "testobj4",
			Name:       "testobjname4",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
		{
			ID:         "testobj5",
			Name:       "testobjname5",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
		{
			ID:         "testobj6",
			Name:       "testobjname6",
			ObjectType: libmodel.LineObjectTypeCRUD,
		},
	})
	mockRoleCursor := createMockCursor(ctrl, []libmodel.Rbac{{
		ID:   "testrole",
		Name: "testrolename",
		PermConfigList: map[string]struct {
			Bitmask int `bson:"checksum"`
		}{
			"testobj1": {Bitmask: 15},
			"testobj2": {Bitmask: 8},
			"testobj3": {Bitmask: 4},
			"testobj4": {Bitmask: 2},
			"testobj5": {Bitmask: 1},
			"testobj6": {Bitmask: 0},
		},
	}})
	mockSubjCursor := createMockCursor(ctrl, nil)

	mockDbClient := createMockDbClient(ctrl, mockObjCursor, mockRoleCursor, mockSubjCursor)
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrolename", "testobjname1", "create"},
		{"testrolename", "testobjname1", "delete"},
		{"testrolename", "testobjname1", "read"},
		{"testrolename", "testobjname1", "update"},
		{"testrolename", "testobjname2", "create"},
		{"testrolename", "testobjname3", "read"},
		{"testrolename", "testobjname4", "update"},
		{"testrolename", "testobjname5", "delete"},
	}

	policy := m["p"]["p"].Policy
	sort.Slice(policy, sortPolicy(policy))

	if !reflect.DeepEqual(policy, expectedPolicy) {
		t.Errorf("expected policy %v but got %v", expectedPolicy, policy)
	}
}

func TestAdapter_LoadPolicy_GivenRole_ShouldAddRWPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObjCursor := createMockCursor(ctrl, []libmodel.Rbac{
		{
			ID:         "testobj1",
			Name:       "testobjname1",
			ObjectType: libmodel.LineObjectTypeRW,
		},
		{
			ID:         "testobj2",
			Name:       "testobjname2",
			ObjectType: libmodel.LineObjectTypeRW,
		},
		{
			ID:         "testobj3",
			Name:       "testobjname3",
			ObjectType: libmodel.LineObjectTypeRW,
		},
		{
			ID:         "testobj4",
			Name:       "testobjname4",
			ObjectType: libmodel.LineObjectTypeRW,
		},
		{
			ID:         "testobj5",
			Name:       "testobjname5",
			ObjectType: libmodel.LineObjectTypeRW,
		},
	})
	mockRoleCursor := createMockCursor(ctrl, []libmodel.Rbac{{
		ID:   "testrole",
		Name: "testrolename",
		PermConfigList: map[string]struct {
			Bitmask int `bson:"checksum"`
		}{
			"testobj1": {Bitmask: 7},
			"testobj2": {Bitmask: 4},
			"testobj3": {Bitmask: 2},
			"testobj4": {Bitmask: 1},
			"testobj5": {Bitmask: 0},
		},
	}})
	mockSubjCursor := createMockCursor(ctrl, nil)

	mockDbClient := createMockDbClient(ctrl, mockObjCursor, mockRoleCursor, mockSubjCursor)
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrolename", "testobjname1", "delete"},
		{"testrolename", "testobjname1", "read"},
		{"testrolename", "testobjname1", "update"},
		{"testrolename", "testobjname2", "read"},
		{"testrolename", "testobjname3", "update"},
		{"testrolename", "testobjname4", "delete"},
	}

	policy := m["p"]["p"].Policy
	sort.Slice(policy, sortPolicy(policy))

	if !reflect.DeepEqual(policy, expectedPolicy) {
		t.Errorf("expected policy %v but got %v", expectedPolicy, policy)
	}
}

func TestAdapter_LoadPolicy_GivenRole_ShouldAddCanPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObjCursor := createMockCursor(ctrl, []libmodel.Rbac{
		{
			ID:   "testobj1",
			Name: "testobjname1",
		},
		{
			ID:   "testobj2",
			Name: "testobjname2",
		},
	})
	mockRoleCursor := createMockCursor(ctrl, []libmodel.Rbac{{
		ID:   "testrole",
		Name: "testrolename",
		PermConfigList: map[string]struct {
			Bitmask int `bson:"checksum"`
		}{
			"testobj1": {Bitmask: 1},
			"testobj2": {Bitmask: 0},
		},
	}})
	mockSubjCursor := createMockCursor(ctrl, nil)

	mockDbClient := createMockDbClient(ctrl, mockObjCursor, mockRoleCursor, mockSubjCursor)
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrolename", "testobjname1", "can"},
	}
	policy := m["p"]["p"].Policy
	sort.Slice(policy, sortPolicy(policy))

	if !reflect.DeepEqual(policy, expectedPolicy) {
		t.Errorf("expected policy %v but got %v", expectedPolicy, policy)
	}
}

func TestAdapter_LoadPolicy_GivenUser_ShouldAddRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObjCursor := createMockCursor(ctrl, nil)
	mockRoleCursor := createMockCursor(ctrl, []libmodel.Rbac{{
		ID:   "testrole",
		Name: "testrolename",
	}})
	mockSubjCursor := createMockCursor(ctrl, []libmodel.Rbac{{
		ID:   "testsubj",
		Role: "testrole",
	}})

	mockDbClient := createMockDbClient(ctrl, mockObjCursor, mockRoleCursor, mockSubjCursor)
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testsubj", "testrolename"},
	}
	policy := m["g"]["g"].Policy
	sort.Slice(policy, sortPolicy(policy))

	if !reflect.DeepEqual(policy, expectedPolicy) {
		t.Errorf("expected policy %v but got %v", expectedPolicy, policy)
	}
}

func createMockCursor(ctrl *gomock.Controller, models []libmodel.Rbac) mongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	if len(models) > 0 {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true).Times(len(models))
		calls := make([]*gomock.Call, len(models))
		for i := range models {
			model := models[i]
			calls[i] = mockCursor.EXPECT().Decode(gomock.Any()).Do(func(m *libmodel.Rbac) {
				*m = model
			}).Return(nil)
		}

		gomock.InOrder(calls...)
	}
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockCursor.EXPECT().Err().Return(nil)

	return mockCursor
}

func createMockDbClient(ctrl *gomock.Controller, objCursor, roleCursor, subjCursor mongo.Cursor) mongo.DbClient {
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.RightsMongoCollection)).Return(mockDbCollection)
	mockDbCollection.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, q bson.M, opts ...*options.FindOptions) (mongo.Cursor, error) {
			switch q["crecord_type"] {
			case libmodel.LineTypeObject:
				return objCursor, nil
			case libmodel.LineTypeRole:
				return roleCursor, nil
			case libmodel.LineTypeSubject:
				return subjCursor, nil
			}

			return nil, nil
		}).
		Times(3)

	return mockDbClient
}

func createCasbinModel() casbinmodel.Model {
	m := casbinmodel.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	return m
}

func sortPolicy(policy [][]string) func(i, j int) bool {
	return func(r, l int) bool {
		length := len(policy[r])
		llength := len(policy[l])
		if llength < length {
			length = llength
		}

		for i := 0; i < length; i++ {
			if policy[r][i] != policy[l][i] {
				return policy[r][i] < policy[l][i]
			}
		}

		return false
	}
}
