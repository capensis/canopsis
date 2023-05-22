package mongoadapter

import (
	"reflect"
	"sort"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libmodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/golang/mock/gomock"
)

func TestAdapter_LoadPolicy_GivenRole_ShouldAddCRUDPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRoleCursor := createMockCursor(ctrl, []role{{
		ID: "testrole",
		Permissions: map[string]permission{
			"testobj1": {Bitmask: 15, Type: libmodel.ObjectTypeCRUD},
			"testobj2": {Bitmask: 8, Type: libmodel.ObjectTypeCRUD},
			"testobj3": {Bitmask: 4, Type: libmodel.ObjectTypeCRUD},
			"testobj4": {Bitmask: 2, Type: libmodel.ObjectTypeCRUD},
			"testobj5": {Bitmask: 1, Type: libmodel.ObjectTypeCRUD},
			"testobj6": {Bitmask: 0, Type: libmodel.ObjectTypeCRUD},
		},
	}})
	mockRoleDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockRoleDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockRoleCursor, nil)
	mockSubjCursor := createMockCursor[user](ctrl, nil)
	mockSubjDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockSubjDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockSubjCursor, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.RoleCollection:
			return mockRoleDbCollection
		case mongo.UserCollection:
			return mockSubjDbCollection
		default:
			return nil
		}
	}).AnyTimes()

	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrole", "testobj1", "create"},
		{"testrole", "testobj1", "delete"},
		{"testrole", "testobj1", "read"},
		{"testrole", "testobj1", "update"},
		{"testrole", "testobj2", "create"},
		{"testrole", "testobj3", "read"},
		{"testrole", "testobj4", "update"},
		{"testrole", "testobj5", "delete"},
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
	mockRoleCursor := createMockCursor(ctrl, []role{{
		ID: "testrole",
		Permissions: map[string]permission{
			"testobj1": {Bitmask: 7, Type: libmodel.ObjectTypeRW},
			"testobj2": {Bitmask: 4, Type: libmodel.ObjectTypeRW},
			"testobj3": {Bitmask: 2, Type: libmodel.ObjectTypeRW},
			"testobj4": {Bitmask: 1, Type: libmodel.ObjectTypeRW},
			"testobj5": {Bitmask: 0, Type: libmodel.ObjectTypeRW},
		},
	}})
	mockRoleDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockRoleDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockRoleCursor, nil)
	mockSubjCursor := createMockCursor[user](ctrl, nil)
	mockSubjDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockSubjDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockSubjCursor, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.RoleCollection:
			return mockRoleDbCollection
		case mongo.UserCollection:
			return mockSubjDbCollection
		default:
			return nil
		}
	}).AnyTimes()
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrole", "testobj1", "delete"},
		{"testrole", "testobj1", "read"},
		{"testrole", "testobj1", "update"},
		{"testrole", "testobj2", "read"},
		{"testrole", "testobj3", "update"},
		{"testrole", "testobj4", "delete"},
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
	mockRoleCursor := createMockCursor(ctrl, []role{{
		ID: "testrole",
		Permissions: map[string]permission{
			"testobj1": {Bitmask: 1},
			"testobj2": {Bitmask: 0},
		},
	}})
	mockRoleDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockRoleDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockRoleCursor, nil)
	mockSubjCursor := createMockCursor[user](ctrl, nil)
	mockSubjDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockSubjDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockSubjCursor, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.RoleCollection:
			return mockRoleDbCollection
		case mongo.UserCollection:
			return mockSubjDbCollection
		default:
			return nil
		}
	}).AnyTimes()
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testrole", "testobj1", "can"},
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
	mockRoleCursor := createMockCursor(ctrl, []role{{
		ID: "testrole",
	}})
	mockSubjCursor := createMockCursor(ctrl, []user{{
		ID:   "testsubj",
		Role: "testrole",
	}})
	mockRoleDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockRoleDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockRoleCursor, nil)
	mockSubjDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockSubjDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockSubjCursor, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.RoleCollection:
			return mockRoleDbCollection
		case mongo.UserCollection:
			return mockSubjDbCollection
		default:
			return nil
		}
	}).AnyTimes()
	adapter := NewAdapter(mockDbClient)
	m := createCasbinModel()
	err := adapter.LoadPolicy(m)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedPolicy := [][]string{
		{"testsubj", "testrole"},
	}
	policy := m["g"]["g"].Policy
	sort.Slice(policy, sortPolicy(policy))

	if !reflect.DeepEqual(policy, expectedPolicy) {
		t.Errorf("expected policy %v but got %v", expectedPolicy, policy)
	}
}

func createMockCursor[T any](ctrl *gomock.Controller, models []T) mongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	if len(models) > 0 {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true).Times(len(models))
		calls := make([]*gomock.Call, len(models))
		for i := range models {
			model := models[i]
			calls[i] = mockCursor.EXPECT().Decode(gomock.Any()).Do(func(m *T) {
				*m = model
			}).Return(nil)
		}

		gomock.InOrder(calls...)
	}
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)

	return mockCursor
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
