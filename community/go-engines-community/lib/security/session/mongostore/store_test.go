package mongostore

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongoStore_New_GivenNoCookie_ShouldReturnNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "testsession"
	codecs := []byte("test")
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Times(0)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, codecs)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	session, err := store.New(req, name)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if session.ID != "" || !session.IsNew || len(session.Values) != 0 {
		t.Errorf("expected new session but got %v", session)
	}
}

func TestMongoStore_New_GivenCookie_ShouldReturnSessionFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionId := "5ed0e4be448c4dd1274c31e3"
	objectId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		t.Fatal("fail to create object id", err)
	}
	values := map[interface{}]interface{}{
		"test": "testvalue",
	}
	name := "testsession"
	codecs := []byte("test")
	encoded, err := securecookie.EncodeMulti(name, values, securecookie.CodecsFromPairs(codecs)...)
	if err != nil {
		t.Fatal("fail to encode session", err)
	}
	data := &sessionData{
		ID:      objectId,
		Data:    encoded,
		Expires: time.Now().Unix(),
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(mockCursor(ctrl, data), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, codecs)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	encoded, err = securecookie.EncodeMulti(name, sessionId, securecookie.CodecsFromPairs(codecs)...)
	if err != nil {
		t.Fatal("fail to encode session", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  name,
		Value: encoded,
	})
	session, err := store.New(req, name)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if session == nil || session.ID != sessionId {
		t.Errorf("expected session id=%v but got %v", sessionId, session)
	}

	if v, ok := session.Values["test"]; !ok || v != values["test"] {
		t.Errorf("expected session values: %v but got %v", values, session.Values)
	}
}

func TestMongoStore_Get_GivenNoCookie_ShouldReturnNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "testsession"
	codecs := []byte("test")
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Times(0)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, codecs)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	session, err := store.Get(req, name)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if session.ID != "" || !session.IsNew || len(session.Values) != 0 {
		t.Errorf("expected new session but got %v", session)
	}
}

func TestMongoStore_Get_GivenCookie_ShouldReturnSessionFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionId := "5ed0e4be448c4dd1274c31e3"
	objectId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		t.Fatal("fail to create object id", err)
	}
	values := map[interface{}]interface{}{
		"test": "testvalue",
	}
	name := "testsession"
	codecs := []byte("test")
	encoded, err := securecookie.EncodeMulti(name, values, securecookie.CodecsFromPairs(codecs)...)
	if err != nil {
		t.Fatal("fail to encode session", err)
	}
	data := &sessionData{
		ID:      objectId,
		Data:    encoded,
		Expires: time.Now().Unix(),
	}
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(mockCursor(ctrl, data), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, codecs)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	encoded, err = securecookie.EncodeMulti(name, sessionId, securecookie.CodecsFromPairs(codecs)...)
	if err != nil {
		t.Fatal("fail to encode session", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  name,
		Value: encoded,
	})
	session, err := store.Get(req, name)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if session == nil || session.ID != sessionId {
		t.Errorf("expected session id=%v but got %v", sessionId, session)
	}

	if v, ok := session.Values["test"]; !ok || v != values["test"] {
		t.Errorf("expected session values: %v but got %v", values, session.Values)
	}

	nextSession, nextErr := store.Get(req, name)

	if nextSession != session || !errors.Is(nextErr, err) {
		t.Errorf("expected second call return the same result")
	}
}

func TestMongoStore_Save_GivenNoCookie_ShouldCreateNewSessionAndSaveSessionToDBAndAddCookieToResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionId := "5ed0e4be448c4dd1274c31e3"
	objectId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		t.Fatal("fail to create object id", err)
	}

	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		InsertOne(gomock.Any(), gomock.Any()).
		Return(objectId, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, []byte("test"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	name := "testsession"
	session := sessions.NewSession(store, name)
	if session != nil {
		session.Values["test"] = "test"
		session.Options.MaxAge = 10
	}
	w := httptest.NewRecorder()
	err = store.Save(req, w, session)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if session == nil || session.ID != sessionId {
		t.Errorf("expected session id=%v but got %v", sessionId, session)
	}

	if len(w.Result().Cookies()) == 0 || w.Result().Cookies()[0].Name != name {
		t.Errorf("expected cookie in response but got nothing")
	}
}

func TestMongoStore_Save_GivenCookie_ShouldUpdateSessionInDBAndAddCookieToResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionId := "5ed0e4be448c4dd1274c31e3"
	name := "testsession"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&mongo.UpdateResult{MatchedCount: 1}, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, []byte("test"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	session := sessions.NewSession(store, name)
	session.ID = sessionId
	session.Values["test"] = "test"
	session.Options.MaxAge = 100
	w := httptest.NewRecorder()
	err := store.Save(req, w, session)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if len(w.Result().Cookies()) == 0 || w.Result().Cookies()[0].Name != name {
		t.Errorf("expected cookie in response but got nothing")
	}
}

func TestMongoStore_Save_GivenCookieAndZeroMaxAge_ShouldDeleteSessionFromDBAndDeleteCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionId := "5ed0e4be448c4dd1274c31e3"
	name := "testsession"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(1), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.SessionMongoCollection)).
		Return(mockDbCollection)

	store := NewStore(mockDbClient, []byte("test"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	session := sessions.NewSession(store, name)
	session.ID = sessionId
	session.Values["test"] = "test"
	session.Options.MaxAge = 0
	w := httptest.NewRecorder()
	err := store.Save(req, w, session)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if len(w.Result().Cookies()) == 0 || w.Result().Cookies()[0].Name != name {
		t.Errorf("expected cookie in response but got nothing")
	}
}

func mockCursor(ctrl *gomock.Controller, data *sessionData) libmongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)

	if data != nil {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true)
		mockCursor.
			EXPECT().
			Decode(gomock.Any()).
			Do(func(val interface{}) {
				if u, ok := val.(*sessionData); ok {
					*u = *data
				}
			}).
			Return(nil)
	} else {
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	}

	mockCursor.EXPECT().Close(gomock.Any())

	return mockCursor
}
