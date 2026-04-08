package websocket_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	mock_validation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/validation"
	mock_websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/websocket"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestHub_Connect_GivenStopRun_ShouldCloseConnection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		closeCalled := false
		readDone := make(chan struct{})
		mockConn.EXPECT().Close().Do(func() {
			closeCalled = true
			close(readDone)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			<-readDone

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		})
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		cancel()
		synctest.Wait()
		if !closeCalled {
			t.Fatalf("close not called after hub stopped")
		}
	})
}

func TestHub_Connect_GivenReadError_ShouldCloseConnection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		closeCalled := false
		readDone := make(chan struct{})
		mockConn.EXPECT().Close().Do(func() {
			closeCalled = true
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			<-readDone

			return -1, nil, errors.New("test error")
		})
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		close(readDone)
		synctest.Wait()
		if !closeCalled {
			t.Fatalf("close not called after read error")
		}
	})
}

func TestHub_Connect_GivenWriteError_ShouldCloseConnection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		closeCalled := false
		readDone := make(chan struct{})
		mockConn.EXPECT().Close().Do(func() {
			closeCalled = true
			close(readDone)
		})
		room := "test"
		joinMsg := &websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		encoder := json.NewEncoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			if joinMsg != nil {
				b, err := encoder.Encode(joinMsg)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				joinMsg = nil

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			<-readDone

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(nil, errors.New("test error"))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		err := roomRegistry.Register(room, websocket.RoomHandlers{})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, json.NewDecoder(), zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		hub.SendMessage(ctx, "test", websocket.ToRoom(room))
		synctest.Wait()
		if !closeCalled {
			t.Fatalf("close not called after write error")
		}
	})
}

func TestHub_Connect_GivenPingWriteError_ShouldCloseConnection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readDone := make(chan struct{})
		closeCalled := false
		mockConn.EXPECT().Close().Do(func() {
			closeCalled = true
			close(readDone)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			<-readDone

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		})
		mockConn.EXPECT().
			WriteControl(gomock.Eq(gorillawebsocket.PingMessage), gomock.Any(), gomock.Any()).
			Return(errors.New("test error"))
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "10s",
			},
		}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		time.Sleep(10 * time.Second)
		synctest.Wait()
		if !closeCalled {
			t.Fatalf("close not called after ping error")
		}
	})
}

func TestHub_Connect_ShouldSendPingMessage(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readDone := make(chan struct{})
		mockConn.EXPECT().Close().Do(func() {
			close(readDone)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			<-readDone

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		})
		pingCalled := false
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.PingMessage), gomock.Any(), gomock.Any()).Do(func(_ int, _ []byte, _ time.Time) {
			pingCalled = true
		})
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "10s",
			},
		}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		time.Sleep(10 * time.Second)
		synctest.Wait()
		if !pingCalled {
			t.Fatalf("ping not called after timeout")
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_ShouldCallOnJoinHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedJoinedMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onJoinCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnJoin: func(ctx context.Context, opts websocket.JoinOptions) error {
				onJoinCalled = true

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if !onJoinCalled {
			t.Fatalf("onJoin not called after join message")
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_ShouldCallOnLeaveHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		expectedJoinedMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		expectedLeftMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageLeft,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(2)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onLeaveCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnLeave: func(ctx context.Context, opts websocket.LeaveOptions) error {
				onLeaveCalled = true

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageLeave,
			Room: room,
		}
		synctest.Wait()
		if !onLeaveCalled {
			t.Fatalf("onLeave not called after leave message")
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedLeftMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenStopRun_ShouldCallOnLeaveHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onLeaveCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnLeave: func(ctx context.Context, opts websocket.LeaveOptions) error {
				onLeaveCalled = true

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		cancel()
		synctest.Wait()
		if !onLeaveCalled {
			t.Fatalf("onLeave not called after leave message")
		}
	})
}

func TestHub_Connect_ShouldCallAuthorizeHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		expectedAuthMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}
		expectedJoinedMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(2)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authorizeCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			Authorize: func(ctx context.Context, userID string) (bool, error) {
				authorizeCalled = true

				return true, nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{ID: "testuser"}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()

		readCh <- websocket.ClientMessage{
			Type:  websocket.ClientMessageAuth,
			Token: "test",
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedAuthMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if !authorizeCalled {
			t.Fatalf("authorize not called after join message")
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_ShouldPeriodicallyCheckAuthorization(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		expectedAuthMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}
		expectedJoinedMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		expectedForbiddenMessage := websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusForbidden,
			Payload: http.StatusText(http.StatusForbidden),
		}
		expectedCloseRoomMessage := websocket.ServerMessage{
			Type: websocket.ServerMessageCloseRoom,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 2)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(4)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(4)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authorizeCalledCount := 0
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			Authorize: func(ctx context.Context, userID string) (bool, error) {
				authorizeCalledCount++

				return authorizeCalledCount == 1, nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{ID: "testuser"}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "1h",
			},
		}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, 10*time.Second, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()

		readCh <- websocket.ClientMessage{
			Type:  websocket.ClientMessageAuth,
			Token: "test",
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedAuthMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		time.Sleep(10 * time.Second)
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedForbiddenMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedCloseRoomMessage); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_SendMessage_ShouldDeliverToRoomConnections(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		room := "test"
		payload := "hello"

		encoder := json.NewEncoder()
		decoder := json.NewDecoder()

		// conn1 — joins the room and receives the message
		readCh1 := make(chan websocket.ClientMessage, 1)
		mockConn1 := mock_websocket.NewMockConn(ctrl)
		mockConn1.EXPECT().Close().Do(func() { close(readCh1) })
		mockConn1.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh1 {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedJoinedMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		expectedInfoMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: payload,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer1 := newServerMessageWriter(writeCh, decoder)
		mockConn1.EXPECT().NextWriter(gomock.Any()).Return(writer1, nil).Times(2)
		mockConn1.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn1.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn1.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn1.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn1.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		// conn2 — does not join the room; no WriteJSON expected
		readCh2 := make(chan websocket.ClientMessage)
		mockConn2 := mock_websocket.NewMockConn(ctrl)
		mockConn2.EXPECT().Close().Do(func() { close(readCh2) })
		mockConn2.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh2 {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		})
		mockConn2.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn2.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn2.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn2.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn1, nil)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn2, nil)

		roomRegistry := websocket.NewRoomRegistry()
		err := roomRegistry.Register(room, websocket.RoomHandlers{})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect conn1: %v", err)
		}

		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect conn2: %v", err)
		}

		synctest.Wait()
		readCh1 <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		hub.SendMessage(ctx, payload, websocket.ToRoom(room))
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedInfoMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_SendMessageToUser_ShouldDeliverToUser(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		userID := "testuser"
		payload := "hello"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		expectedAuthMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}
		expectedJoinedMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		expectedInfoMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: payload,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(3)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		err := roomRegistry.Register(room, websocket.RoomHandlers{})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{ID: userID}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type:  websocket.ClientMessageAuth,
			Token: "testtoken",
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedAuthMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		sent := hub.SendMessageToUser(ctx, payload, room, userID)
		synctest.Wait()
		if !sent {
			t.Fatal("expected message to be sent to user")
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedInfoMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_LeaveRoom_ShouldSendCloseRoomAndCallOnLeave(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedJoinedMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		expectedCloseRoomMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageCloseRoom,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(2)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onLeaveCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnLeave: func(ctx context.Context, opts websocket.LeaveOptions) error {
				onLeaveCalled = true

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		hub.LeaveRoom(ctx, room)
		synctest.Wait()
		if !onLeaveCalled {
			t.Fatal("onLeave not called after hub LeaveRoom")
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedCloseRoomMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenJoinMessageWithoutAuth_ShouldReturn401(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedUnauthMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			Authorize: func(ctx context.Context, userID string) (bool, error) {
				return true, nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{ID: "testuser"}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		// send join without authenticating first
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedUnauthMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenAuthMessageWithInvalidToken_ShouldReturn401(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedUnauthMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, errors.New("invalid token")
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type:  websocket.ClientMessageAuth,
			Token: "badtoken",
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedUnauthMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenCheckAuthWithExpiredToken_ShouldReturn401(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedAuthMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}
		expectedUnauthMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(2)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		callCount := 0
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			callCount++
			if callCount == 1 {
				return websocket.User{ID: "testuser"}, nil
			}

			return websocket.User{}, errors.New("token expired")
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "1h",
			},
		}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, 10*time.Second, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type:  websocket.ClientMessageAuth,
			Token: "testtoken",
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedAuthMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		time.Sleep(10 * time.Second)
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedUnauthMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenJoinUnregisteredRoom_ShouldReturn404(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "unknown-room"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedNotFoundMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusNotFound,
			Payload: http.StatusText(http.StatusNotFound),
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedNotFoundMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenOnJoinReturnsJoinError_ShouldSendInfoMessage(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		joinErrPayload := map[string]string{"msg": "room is full"}
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedInfoMsg := websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: joinErrPayload,
		}
		expectedCloseMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageCloseRoom,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 2)
		defer close(writeCh)
		writer := newServerMessageWriter(writeCh, decoder)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(writer, nil).Times(2)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnJoin: func(ctx context.Context, opts websocket.JoinOptions) error {
				return websocket.NewJoinError(errors.New("room is full"), joinErrPayload)
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedInfoMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedCloseMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenDoubleJoin_ShouldCallOnJoinOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		expectedJoinedMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageJoined,
			Room: room,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onJoinCount := 0
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnJoin: func(ctx context.Context, opts websocket.JoinOptions) error {
				onJoinCount++

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedJoinedMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if onJoinCount != 1 {
			t.Fatalf("expected onJoin called once, got %d", onJoinCount)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenClientPingMessage_ShouldSendClientPong(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		expectedPongMsg := websocket.ServerMessage{
			Type: websocket.ServerMessageClientPong,
		}
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err := hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageClientPing,
		}
		synctest.Wait()
		if diff := pretty.Compare(readServerMessageCh(writeCh), expectedPongMsg); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenInfoMessage_ShouldCallOnMessageHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mockConn := mock_websocket.NewMockConn(ctrl)
		readCh := make(chan websocket.ClientMessage, 2)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		encoder := json.NewEncoder()
		decoder := json.NewDecoder()
		mockConn.EXPECT().NextReader().DoAndReturn(func() (int, io.Reader, error) {
			for v := range readCh {
				b, err := encoder.Encode(v)
				if err != nil {
					t.Fatalf("cannot encode join message: %v", err)
				}

				return gorillawebsocket.TextMessage, bytes.NewReader(b), nil
			}

			return -1, nil, &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		writeCh := make(chan websocket.ServerMessage, 1)
		defer close(writeCh)
		mockConn.EXPECT().NextWriter(gomock.Any()).Return(newServerMessageWriter(writeCh, decoder), nil)
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		onMessageCalled := false
		err := roomRegistry.Register(room, websocket.RoomHandlers{
			OnMessage: func(ctx context.Context, opts websocket.MessageOptions) error {
				onMessageCalled = true

				return nil
			},
		})
		if err != nil {
			t.Fatalf("cannot register room: %v", err)
		}

		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())
		go func() {
			hub.Run(ctx)
		}()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		err = hub.Connect(w, r)
		if err != nil {
			t.Fatalf("cannot connect: %v", err)
		}

		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageInfo,
			Room: room,
		}
		synctest.Wait()
		if !onMessageCalled {
			t.Fatal("onMessage not called after info message")
		}

		cancel()
		synctest.Wait()
	})
}

func TestHub_Connect_GivenFullBuffer_ShouldReturnError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConn := mock_websocket.NewMockConn(ctrl)
		mockConn.EXPECT().Close()

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil).Times(11)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (websocket.User, error) {
			return websocket.User{}, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		mockErrTrans := mock_validation.NewMockErrorTranslator(ctrl)

		encoder := json.NewEncoder()
		decoder := json.NewDecoder()

		// hub.Run is NOT called so registerCh never drains
		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, mockErrTrans,
			encoder, decoder, zerolog.Nop())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		for i := range 10 {
			err := hub.Connect(w, r)
			if err != nil {
				t.Fatalf("expected no error on connect %d: %v", i, err)
			}
		}

		err := hub.Connect(w, r)
		if err == nil {
			t.Fatal("expected error when registerCh buffer is full")
		}
	})
}

func newServerMessageWriter(ch chan<- websocket.ServerMessage, decoder encoding.Decoder) *serverMessageWriter {
	return &serverMessageWriter{
		ch:      ch,
		decoder: decoder,
	}
}

type serverMessageWriter struct {
	ch      chan<- websocket.ServerMessage
	decoder encoding.Decoder
}

func (w *serverMessageWriter) Write(b []byte) (int, error) {
	received := websocket.ServerMessage{}
	err := w.decoder.Decode(b, &received)
	if err != nil {
		return 0, err
	}

	select {
	case w.ch <- received:
	default:
		return 0, fmt.Errorf("write channel full")
	}

	return len(b), nil
}

func (w *serverMessageWriter) Close() error {
	return nil
}

func readServerMessageCh(ch <-chan websocket.ServerMessage) websocket.ServerMessage {
	select {
	case msg := <-ch:
		return msg
	default:
		return websocket.ServerMessage{}
	}
}
