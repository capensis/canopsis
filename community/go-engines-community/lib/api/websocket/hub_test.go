package websocket_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	mock_websocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/websocket"
	gorillawebsocket "github.com/gorilla/websocket"
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(v any) error {
			<-readDone

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		})
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(v any) error {
			<-readDone

			return errors.New("test error")
		})
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			if joinMsg != nil {
				*msg = *joinMsg
				joinMsg = nil

				return nil
			}

			<-readDone

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Any()).Return(errors.New("test error"))
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(v any) error {
			<-readDone

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
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
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "10s",
			},
		}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(v any) error {
			<-readDone

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
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
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "10s",
			},
		}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
			Type: websocket.ClientMessageLeave,
			Room: room,
		}
		synctest.Wait()
		if !onLeaveCalled {
			t.Fatalf("onLeave not called after leave message")
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "testuser", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		if !authorizeCalled {
			t.Fatalf("authorize not called after join message")
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
		readCh := make(chan websocket.ClientMessage, 1)
		mockConn.EXPECT().Close().Do(func() {
			close(readCh)
		})
		room := "test"
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}))
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusForbidden,
			Payload: http.StatusText(http.StatusForbidden),
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "testuser", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "1h",
			},
		}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, 10*time.Second, zerolog.Nop())
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
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		time.Sleep(10 * time.Second)
		synctest.Wait()
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

		// conn1 — joins the room and receives the message
		readCh1 := make(chan websocket.ClientMessage, 1)
		mockConn1 := mock_websocket.NewMockConn(ctrl)
		mockConn1.EXPECT().Close().Do(func() { close(readCh1) })
		mockConn1.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh1 {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn1.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: payload,
		}))
		mockConn1.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn1.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn1.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn1.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn1.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		// conn2 — does not join the room; no WriteJSON expected
		readCh2 := make(chan websocket.ClientMessage)
		mockConn2 := mock_websocket.NewMockConn(ctrl)
		mockConn2.EXPECT().Close().Do(func() { close(readCh2) })
		mockConn2.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh2 {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		hub.SendMessage(ctx, payload, websocket.ToRoom(room))
		synctest.Wait()

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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}))
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: payload,
		}))
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return userID, nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		readCh <- websocket.ClientMessage{
			Type: websocket.ClientMessageJoin,
			Room: room,
		}
		synctest.Wait()
		sent := hub.SendMessageToUser(ctx, payload, room, userID)
		synctest.Wait()
		if !sent {
			t.Fatal("expected message to be sent to user")
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageCloseRoom,
			Room: room,
		}))
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		hub.LeaveRoom(ctx, room)
		synctest.Wait()
		if !onLeaveCalled {
			t.Fatal("onLeave not called after hub LeaveRoom")
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}))
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "testuser", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", errors.New("invalid token")
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageAuthSuccess,
		}))
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Error:   http.StatusUnauthorized,
			Payload: http.StatusText(http.StatusUnauthorized),
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Times(2)
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		callCount := 0
		authenticate := func(ctx context.Context, token string) (string, error) {
			callCount++
			if callCount == 1 {
				return "testuser", nil
			}

			return "", errors.New("token expired")
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{
			API: config.SectionApi{
				WebsocketPingInterval: "1h",
			},
		}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, 10*time.Second, zerolog.Nop())
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
		time.Sleep(10 * time.Second)
		synctest.Wait()

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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageError,
			Room:    room,
			Error:   http.StatusNotFound,
			Payload: http.StatusText(http.StatusNotFound),
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type:    websocket.ServerMessageInfo,
			Room:    room,
			Payload: joinErrPayload,
		}))
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageCloseRoom,
			Room: room,
		}))
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(2)
		mockConn.EXPECT().WriteJSON(gomock.Eq(websocket.ServerMessage{
			Type: websocket.ServerMessageClientPong,
		}))
		mockConn.EXPECT().SetWriteDeadline(gomock.Any())
		mockConn.EXPECT().SetReadDeadline(gomock.Any()).AnyTimes()
		mockConn.EXPECT().SetPongHandler(gomock.Any()).AnyTimes()
		mockConn.EXPECT().RemoteAddr().Return(&net.IPAddr{}).AnyTimes()
		mockConn.EXPECT().WriteControl(gomock.Eq(gorillawebsocket.CloseMessage), gomock.Any(), gomock.Any())

		mockUpgrader := mock_websocket.NewMockUpgrader(ctrl)
		mockUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockConn, nil)

		roomRegistry := websocket.NewRoomRegistry()
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		mockConn.EXPECT().ReadJSON(gomock.Any()).DoAndReturn(func(msg *websocket.ClientMessage) error {
			for v := range readCh {
				*msg = v

				return nil
			}

			return &gorillawebsocket.CloseError{Code: gorillawebsocket.CloseNormalClosure}
		}).Times(3)
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

		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())
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
		authenticate := func(ctx context.Context, token string) (string, error) {
			return "", nil
		}
		configProvider := config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop())

		// hub.Run is NOT called so registerCh never drains
		hub := websocket.NewHub(mockUpgrader, roomRegistry, authenticate, configProvider, time.Hour, zerolog.Nop())

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
