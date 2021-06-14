package userprovider

import (
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongoProvider_FindByUsername_GivenID_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testname"
	m := &model.Rbac{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	expectedUser := security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          username,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, m), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByUsername(username)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user == nil || *user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestMongoProvider_FindByUsername_GivenID_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testname"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          username,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, nil), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByUsername(username)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestMongoProvider_FindByAuthApiKey_GivenID_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApiKey := "testkey"
	m := &model.Rbac{
		ID:             "testid",
		AuthApiKey:     userApiKey,
		HashedPassword: "testhash",
	}
	expectedUser := security.User{
		ID:             "testid",
		AuthApiKey:     userApiKey,
		HashedPassword: "testhash",
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"authkey":      userApiKey,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, m), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByAuthApiKey(userApiKey)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user == nil || *user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestMongoProvider_FindByAuthApiKey_GivenID_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userApiKey := "testkey"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"authkey":      userApiKey,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, nil), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByAuthApiKey(userApiKey)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestMongoProvider_FindByID_GivenID_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "testid"
	m := &model.Rbac{
		ID:             userID,
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	expectedUser := security.User{
		ID:             userID,
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          userID,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, m), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByID(userID)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user == nil || *user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestMongoProvider_FindByID_GivenID_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "testid"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          userID,
	}
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Eq(filter)).
		Return(mockUserCursor(ctrl, nil), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	user, err := p.FindByID(userID)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestMongoProvider_Save_GivenUser_ShouldUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "testid"
	expectedModel := model.Rbac{
		ID:             userID,
		Type:           model.LineTypeSubject,
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	user := security.User{
		ID:             userID,
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	filter := bson.M{"_id": userID}
	mockDbCollection.
		EXPECT().
		UpdateOne(
			gomock.Any(),
			gomock.Eq(filter),
			gomock.Eq(bson.M{"$set": expectedModel}),
			gomock.Eq(options.Update().SetUpsert(true)),
		).
		Return(&mongo.UpdateResult{
			MatchedCount:  1,
			ModifiedCount: 1,
			UpsertedCount: 0,
			UpsertedID:    nil,
		}, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)

	p := NewMongoProvider(mockDbClient)
	err := p.Save(&user)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func mockUserCursor(ctrl *gomock.Controller, user *model.Rbac) libmongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)

	if user != nil {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true)
		mockCursor.
			EXPECT().
			Decode(gomock.Any()).
			Do(func(val interface{}) {
				if u, ok := val.(*model.Rbac); ok {
					*u = *user
				}
			}).
			Return(nil)
	} else {
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	}

	mockCursor.EXPECT().Close(gomock.Any())

	return mockCursor
}
