package alarm_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/alarm"
	mock_websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/websocket"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestWatcher_StartWatch_GivenMultipleConnsWithTheSameRequest_ShouldCreateOneChangeStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId := "test-alarm"
	roomId := "test-room"
	connId1 := "test-conn-1"
	connId2 := "test-conn-2"
	connId3 := "test-conn-3"
	userId1 := "test-user-1"
	userId2 := "test-user-2"
	alarmResForUserId1 := &alarm.Alarm{ID: userId1}
	alarmResForUserId2 := &alarm.Alarm{ID: userId2}
	changeCh := make(chan struct{}, 1)
	done := make(chan struct{})
	mockChangeStream := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
	}) {
		changeEvent.DocumentKey.ID = alarmId
	})
	mockChangeStream.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		close(done)
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(mockChangeStream, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockStore.EXPECT().GetByID(gomock.Any(), gomock.Eq(alarmId), gomock.Eq(userId1), gomock.Any()).Return(alarmResForUserId1, nil)
	mockStore.EXPECT().GetByID(gomock.Any(), gomock.Eq(alarmId), gomock.Eq(userId2), gomock.Any()).Return(alarmResForUserId2, nil)
	mockHub := mock_websocket.NewMockHub(ctrl)
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId1}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForUserId1))
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId2, connId3}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForUserId2))

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	data := []string{alarmId}
	err := w.StartWatch(ctx, connId1, userId1, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatch(ctx, connId2, userId2, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatch(ctx, connId3, userId2, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	changeCh <- struct{}{}
	close(changeCh)

	select {
	case <-time.After(time.Second):
		t.Error("test is executing too long")
	case <-done:
	}
}

func TestWatcher_StartWatch_GivenMultipleConnsWithDiffRequest_ShouldCreateMultipleChangeStreams(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId1 := "test-alarm-1"
	alarmId2 := "test-alarm-2"
	roomId := "test-room"
	connId1 := "test-conn-1"
	connId2 := "test-conn-2"
	userId := "test-user-2"
	alarmResForConnId1 := &alarm.Alarm{ID: alarmId1}
	alarmResForConnId2 := &alarm.Alarm{ID: alarmId2}
	changeCh1 := make(chan struct{}, 1)
	changeCh2 := make(chan struct{}, 1)
	done := make(chan struct{}, 2)
	defer close(done)
	mockChangeStream1 := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream1.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh1:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream1.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
	}) {
		changeEvent.DocumentKey.ID = alarmId1
	})
	mockChangeStream1.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		done <- struct{}{}
	})
	mockChangeStream2 := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream2.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh2:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream2.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
	}) {
		changeEvent.DocumentKey.ID = alarmId2
	})
	mockChangeStream2.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		done <- struct{}{}
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(mockChangeStream1, nil)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(mockChangeStream2, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockStore.EXPECT().GetByID(gomock.Any(), gomock.Eq(alarmId1), gomock.Eq(userId), gomock.Any()).Return(alarmResForConnId1, nil)
	mockStore.EXPECT().GetByID(gomock.Any(), gomock.Eq(alarmId2), gomock.Eq(userId), gomock.Any()).Return(alarmResForConnId2, nil)
	mockHub := mock_websocket.NewMockHub(ctrl)
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId1}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForConnId1))
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId2}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForConnId2))

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	err := w.StartWatch(ctx, connId1, userId, roomId, []string{alarmId1})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatch(ctx, connId2, userId, roomId, []string{alarmId2})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	changeCh1 <- struct{}{}
	close(changeCh1)
	changeCh2 <- struct{}{}
	close(changeCh2)

	closedStreams := 0
	after := time.After(time.Second)
	for {
		select {
		case <-after:
			t.Error("test is executing too long")
		case <-done:
			closedStreams++
			if closedStreams == 2 {
				return
			}
		}
	}
}

func TestWatcher_StartWatchDetails_GivenMultipleConnsWithTheSameRequest_ShouldCreateOneChangeStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId := "test-alarm"
	roomId := "test-room"
	connId1 := "test-conn-1"
	connId2 := "test-conn-2"
	connId3 := "test-conn-3"
	userId1 := "test-user-1"
	userId2 := "test-user-2"
	alarmResForUserId1 := &alarm.Details{}
	alarmResForUserId1.Entity.ID = userId1
	alarmResForUserId2 := &alarm.Details{}
	alarmResForUserId2.Entity.ID = userId2
	changeCh := make(chan struct{}, 1)
	done := make(chan struct{})
	mockChangeStream := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
		FullDocument types.Alarm `bson:"fullDocument"`
	}) {
		changeEvent.DocumentKey.ID = alarmId
	})
	mockChangeStream.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		close(done)
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockChangeStream, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockStore.EXPECT().GetDetails(gomock.Any(), gomock.Any(), gomock.Eq(userId1)).Return(alarmResForUserId1, nil)
	mockStore.EXPECT().GetDetails(gomock.Any(), gomock.Any(), gomock.Eq(userId2)).Return(alarmResForUserId2, nil)
	mockHub := mock_websocket.NewMockHub(ctrl)
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId1}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForUserId1))
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId2, connId3}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForUserId2))

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	data := []alarm.DetailsRequest{{ID: alarmId}}
	err := w.StartWatchDetails(ctx, connId1, userId1, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatchDetails(ctx, connId2, userId2, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatchDetails(ctx, connId3, userId2, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	changeCh <- struct{}{}
	close(changeCh)

	select {
	case <-time.After(time.Second):
		t.Error("test is executing too long")
	case <-done:
	}
}

func TestWatcher_StartWatchDetails_GivenMultipleConnsWithDiffRequest_ShouldCreateMultipleChangeStreams(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId1 := "test-alarm-1"
	alarmId2 := "test-alarm-2"
	roomId := "test-room"
	connId1 := "test-conn-1"
	connId2 := "test-conn-2"
	userId := "test-user-2"
	alarmResForConnId1 := &alarm.Details{}
	alarmResForConnId1.Entity.ID = alarmId1
	alarmResForConnId2 := &alarm.Details{}
	alarmResForConnId2.Entity.ID = alarmId2
	changeCh1 := make(chan struct{}, 1)
	changeCh2 := make(chan struct{}, 1)
	done := make(chan struct{}, 2)
	defer close(done)
	mockChangeStream1 := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream1.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh1:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream1.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
		FullDocument types.Alarm `bson:"fullDocument"`
	}) {
		changeEvent.DocumentKey.ID = alarmId1
	})
	mockChangeStream1.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		done <- struct{}{}
	})
	mockChangeStream2 := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream2.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		select {
		case <-ctx.Done():
			return false
		case _, ok := <-changeCh2:
			if ok {
				return true
			}

			return false
		}
	}).Times(2)
	mockChangeStream2.EXPECT().Decode(gomock.Any()).Do(func(changeEvent *struct {
		DocumentKey struct {
			ID string `bson:"_id"`
		} `bson:"documentKey"`
		FullDocument types.Alarm `bson:"fullDocument"`
	}) {
		changeEvent.DocumentKey.ID = alarmId2
	})
	mockChangeStream2.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		done <- struct{}{}
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockChangeStream1, nil)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockChangeStream2, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockStore.EXPECT().GetDetails(gomock.Any(), gomock.Eq(alarm.DetailsRequest{ID: alarmId1}), gomock.Eq(userId)).Return(alarmResForConnId1, nil)
	mockStore.EXPECT().GetDetails(gomock.Any(), gomock.Eq(alarm.DetailsRequest{ID: alarmId2}), gomock.Eq(userId)).Return(alarmResForConnId2, nil)
	mockHub := mock_websocket.NewMockHub(ctrl)
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId1}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForConnId1))
	mockHub.EXPECT().SendGroupRoomByConnections(gomock.Eq([]string{connId2}), gomock.Any(),
		gomock.Eq(roomId), gomock.Eq(alarmResForConnId2))

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	err := w.StartWatchDetails(ctx, connId1, userId, roomId, []alarm.DetailsRequest{{ID: alarmId1}})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StartWatchDetails(ctx, connId2, userId, roomId, []alarm.DetailsRequest{{ID: alarmId2}})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	changeCh1 <- struct{}{}
	close(changeCh1)
	changeCh2 <- struct{}{}
	close(changeCh2)

	closedStreams := 0
	after := time.After(time.Second)
	for {
		select {
		case <-after:
			t.Error("test is executing too long")
		case <-done:
			closedStreams++
			if closedStreams == 2 {
				return
			}
		}
	}
}

func TestWatcher_StopWatch_GivenStartWatch_ShouldCloseChangeStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId := "test-alarm"
	roomId := "test-room"
	connId := "test-conn"
	userId := "test-user"
	changeCh := make(chan struct{})
	defer close(changeCh)
	done := make(chan struct{})
	mockChangeStream := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream.EXPECT().Next(gomock.Any()).DoAndReturn(func(nextCtx context.Context) bool {
		select {
		case <-nextCtx.Done():
			return false
		case <-changeCh:
			return false
		}
	})
	mockChangeStream.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		close(done)
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any()).Return(mockChangeStream, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockHub := mock_websocket.NewMockHub(ctrl)

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	data := []string{alarmId}
	err := w.StartWatch(ctx, connId, userId, roomId, data)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StopWatch(connId, roomId)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	select {
	case <-time.After(time.Second):
		t.Error("test is executing too long")
	case <-done:
	}
}

func TestWatcher_StopWatch_GivenStartWatchDetails_ShouldCloseChangeStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	alarmId := "test-alarm"
	roomId := "test-room"
	connId := "test-conn"
	userId := "test-user"
	changeCh := make(chan struct{})
	defer close(changeCh)
	done := make(chan struct{})
	mockChangeStream := mock_mongo.NewMockChangeStream(ctrl)
	mockChangeStream.EXPECT().Next(gomock.Any()).DoAndReturn(func(nextCtx context.Context) bool {
		select {
		case <-nextCtx.Done():
			return false
		case <-changeCh:
			return false
		}
	})
	mockChangeStream.EXPECT().Close(gomock.Any()).Do(func(_ context.Context) {
		close(done)
	})
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockChangeStream, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.AlarmMongoCollection)).Return(mockDbCollection)
	mockStore := mock_alarm.NewMockStore(ctrl)
	mockHub := mock_websocket.NewMockHub(ctrl)

	w := alarm.NewWatcher(mockDbClient, mockHub, mockStore, json.NewEncoder(), json.NewDecoder(), zerolog.Nop())
	err := w.StartWatchDetails(ctx, connId, userId, roomId, []alarm.DetailsRequest{{ID: alarmId}})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = w.StopWatch(connId, roomId)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	select {
	case <-time.After(time.Second):
		t.Error("test is executing too long")
	case <-done:
	}
}
