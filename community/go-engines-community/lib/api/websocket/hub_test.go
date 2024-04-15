package websocket_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	libwebsocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	mock_websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/websocket"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const waitTimeout = time.Second

func TestHub_Send_GivenNotJoinedToRoomConnection_ShouldNotSendMessageToConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	msgBody := map[string]string{"test": "msg"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		hub.Send(room, msgBody)
		cancel()
		return &websocket.CloseError{Code: websocket.CloseNormalClosure}
	})
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Send_GivenJoinedToRoomConnection_ShouldSendMessageToConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	msgBody := map[string]string{"test": "msg"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		default:
			hub.Send(room, msgBody)
			cancel()
			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(2)
	mockConnection.EXPECT().WriteJSON(gomock.Eq(libwebsocket.WMessage{
		Type: libwebsocket.WMessageSuccess,
		Room: room,
		Msg:  msgBody,
	})).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Send_GivenLeftRoomConnection_ShouldNotSendMessageToConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	msgBody := map[string]string{"test": "msg"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		case 2:
			hub.Send(room, msgBody)
			msg.Type = libwebsocket.RMessageLeave
			return nil
		default:
			hub.Send(room, msgBody)
			cancel()
			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(3)
	mockConnection.EXPECT().WriteJSON(gomock.Eq(libwebsocket.WMessage{
		Type: libwebsocket.WMessageSuccess,
		Room: room,
		Msg:  msgBody,
	})).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Send_GivenDisconnectedConnection_ShouldNotSendMessageToConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	msgBody := map[string]string{"test": "msg"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		default:
			hub.Send(room, msgBody)

			go func() {
				time.Sleep(waitTimeout / 2)
				hub.Send(room, msgBody)
				cancel()
			}()

			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(2)
	mockConnection.EXPECT().WriteJSON(gomock.Eq(libwebsocket.WMessage{
		Type: libwebsocket.WMessageSuccess,
		Room: room,
		Msg:  msgBody,
	})).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Send_GivenErrOnWriteMessage_ShouldCloseConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	msgBody := map[string]string{"test": "msg"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		default:
			hub.Send(room, msgBody)
			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(2)
	mockConnection.EXPECT().WriteJSON(gomock.Any()).Return(errors.New("test error"))
	mockConnection.EXPECT().RemoteAddr().Return(&net.TCPAddr{})
	mockConnection.EXPECT().Close().Do(func() {
		cancel()
	}).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Send_GivenErrOnWriteError_ShouldCloseConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().ReadJSON(gomock.Any()).Do(func(msg *libwebsocket.RMessage) {
		msg.Room = room
		msg.Type = libwebsocket.RMessageJoin
	}).Return(nil)
	mockConnection.EXPECT().WriteJSON(gomock.Any()).Return(errors.New("test error"))
	mockConnection.EXPECT().RemoteAddr().Return(&net.TCPAddr{})
	mockConnection.EXPECT().Close().Do(func() {
		cancel()
	}).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Connect_GivenUnauthUser_ShouldNotJoinToRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Eq(""), gomock.Eq(room)).Return(false, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		default:
			cancel()
			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(2)
	mockConnection.EXPECT().WriteJSON(gomock.Eq(libwebsocket.WMessage{
		Type:  libwebsocket.WMessageFail,
		Room:  room,
		Error: http.StatusUnauthorized,
	})).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func TestHub_Connect_GivenInvalidRMessage_ShouldSendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	room := "test-room"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test-ws", nil)
	mockConnection := mock_websocket.NewMockConnection(ctrl)
	mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
	mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConnection, nil)
	mockAuthorizer := mock_websocket.NewMockAuthorizer(ctrl)
	mockAuthorizer.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	hub := libwebsocket.NewHub(ctx, mockUpgrader, mockAuthorizer, time.Hour, zerolog.Nop())
	mockConnection.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
	mockConnection.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
	readTime := 0
	mockConnection.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *libwebsocket.RMessage) error {
		readTime++
		msg.Room = room
		switch readTime {
		case 1:
			msg.Type = libwebsocket.RMessageJoin
			return nil
		case 2:
			return &json.SyntaxError{}
		default:
			cancel()
			return &websocket.CloseError{Code: websocket.CloseNormalClosure}
		}
	}).Times(3)
	mockConnection.EXPECT().WriteJSON(gomock.Eq(libwebsocket.WMessage{
		Type:  libwebsocket.WMessageFail,
		Error: http.StatusBadRequest,
		Msg:   "cannot parse JSON",
	})).Return(nil)
	mockConnection.EXPECT().WriteControl(gomock.Eq(websocket.CloseMessage), gomock.Any(), gomock.Any()).MaxTimes(1)

	go func() {
		hub.Start(ctx)
		close(done)
	}()

	err := hub.Connect(w, r)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	waitDone(t, done)
}

func waitDone(t *testing.T, done <-chan struct{}) {
	select {
	case <-time.After(waitTimeout):
		t.Error("timeout expired")
	case <-done:
	}
}
