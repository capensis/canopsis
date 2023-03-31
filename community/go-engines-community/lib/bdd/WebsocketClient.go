package bdd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	libwebsocket "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/cucumber/godog"
	"github.com/gorilla/websocket"
	"github.com/kylelemons/godebug/pretty"
)

// WebsocketClient represents utility struct which implements Websocket steps to feature context.
type WebsocketClient struct {
	apiURL    string
	templater *Templater

	connsMx       sync.Mutex
	conns         map[int]*websocket.Conn
	lastConnIndex int
}

const wsWriteWait = 5 * time.Second

func NewWebsocketClient(
	apiURL string,
	templater *Templater,
) *WebsocketClient {
	return &WebsocketClient{
		apiURL:    apiURL,
		templater: templater,

		conns:         make(map[int]*websocket.Conn),
		lastConnIndex: -1,
	}
}

func (c *WebsocketClient) AfterScenario(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return ctx, nil
	}

	ctx = c.delConn(ctx)
	err = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(wsWriteWait))
	if err != nil {
		return ctx, err
	}
	err = conn.Close()
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c *WebsocketClient) IConnect(ctx context.Context) (context.Context, error) {
	if _, err := c.getConn(ctx); err == nil {
		return ctx, errors.New("websocket connection already exists")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.apiURL, nil)
	if err != nil {
		return ctx, err
	}

	return c.setConn(ctx, conn), nil
}

func (c *WebsocketClient) ISend(ctx context.Context, doc string) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	msg := libwebsocket.RMessage{}
	err = json.Unmarshal(b.Bytes(), &msg)
	if err != nil {
		return err
	}

	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	return conn.WriteJSON(msg)
}

func (c *WebsocketClient) IWaitMessage(ctx context.Context, doc string) error {
	return c.iWaitMessage(ctx, doc, func(receivedMsg, _ any) any {
		return receivedMsg
	})
}

func (c *WebsocketClient) IWaitMessageWhichContains(ctx context.Context, doc string) error {
	return c.iWaitMessage(ctx, doc, func(receivedMsg, expectedMsg any) any {
		return getPartialResponse(receivedMsg, expectedMsg)
	})
}

func (c *WebsocketClient) IWaitNextMessage(ctx context.Context, doc string) error {
	return c.iWaitNextMessage(ctx, doc, func(receivedMsg, _ any) any {
		return receivedMsg
	})
}

func (c *WebsocketClient) IWaitNextMessageWhichContains(ctx context.Context, doc string) error {
	return c.iWaitNextMessage(ctx, doc, func(receivedMsg, expectedMsg any) any {
		return getPartialResponse(receivedMsg, expectedMsg)
	})
}

func (c *WebsocketClient) IAuthenticate(ctx context.Context) error {
	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	token, ok := getApiAuthToken(ctx)
	if !ok {
		return errors.New("token not found")
	}

	err = conn.WriteJSON(libwebsocket.RMessage{
		Type:  libwebsocket.RMessageAuth,
		Token: token,
	})
	if err != nil {
		return err
	}

	return c.waitNextMessage(
		conn,
		map[string]any{
			"type": libwebsocket.WMessageAuthSuccess,
		},
		func(receivedMsg any, expectedMsg any) any {
			return getPartialResponse(receivedMsg, expectedMsg)
		},
	)
}

func (c *WebsocketClient) ISubscribeToRoom(ctx context.Context, room string) error {
	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	b, err := c.templater.Execute(ctx, room)
	if err != nil {
		return err
	}
	room = b.String()
	return conn.WriteJSON(libwebsocket.RMessage{
		Type: libwebsocket.RMessageJoin,
		Room: room,
	})
}

func (c *WebsocketClient) IWaitMessageFromRoom(ctx context.Context, room, doc string) error {
	return c.iWaitMessageFromRoom(ctx, room, doc, func(receivedMsg, _ any) any {
		return receivedMsg
	})
}

func (c *WebsocketClient) IWaitMessageFromRoomWhichContains(ctx context.Context, room, doc string) error {
	return c.iWaitMessageFromRoom(ctx, room, doc, func(receivedMsg, expectedMsg any) any {
		return getPartialResponse(receivedMsg, expectedMsg)
	})
}

func (c *WebsocketClient) getConn(ctx context.Context) (*websocket.Conn, error) {
	idx, ok := getWebsocketConn(ctx)
	if !ok {
		return nil, fmt.Errorf("websocket connection is not open")
	}

	c.connsMx.Lock()
	defer c.connsMx.Unlock()

	return c.conns[idx], nil
}

func (c *WebsocketClient) setConn(ctx context.Context, conn *websocket.Conn) context.Context {
	c.connsMx.Lock()
	defer c.connsMx.Unlock()
	c.lastConnIndex++
	index := c.lastConnIndex
	c.conns[index] = conn
	return setWebsocketConn(ctx, index)
}

func (c *WebsocketClient) delConn(ctx context.Context) context.Context {
	idx, ok := getWebsocketConn(ctx)
	if !ok {
		return ctx
	}

	c.connsMx.Lock()
	defer c.connsMx.Unlock()
	delete(c.conns, idx)

	return setWebsocketConn(ctx, -1)
}

func (c *WebsocketClient) iWaitMessage(
	ctx context.Context,
	doc string,
	transformReceivedMsg func(receivedMsg any, expectedMsg any) any,
) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedMsg, err := unmarshalJson(b.Bytes())
	if err != nil {
		return err
	}

	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	msgsMx := sync.Mutex{}
	caughtMsgs := make([]any, 0)
	go func() {
		defer close(errCh)
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			receivedMsg, err := unmarshalJson(b)
			if err != nil {
				errCh <- err
				return
			}

			diff := pretty.Compare(transformReceivedMsg(receivedMsg, expectedMsg), expectedMsg)
			if diff != "" {
				msgsMx.Lock()
				caughtMsgs = append(caughtMsgs, receivedMsg)
				msgsMx.Unlock()
				continue
			}

			return
		}
	}()

	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()
	select {
	case err, ok := <-errCh:
		if !ok {
			return nil
		}

		return err
	case <-timer.C:
		msgsMx.Lock()
		defer msgsMx.Unlock()
		return fmt.Errorf("reached timeout: caught %d messages but none of them matches to expected message\n%s\n",
			len(caughtMsgs), pretty.Compare(caughtMsgs, []any{expectedMsg}))
	}
}

func (c *WebsocketClient) iWaitNextMessage(
	ctx context.Context,
	doc string,
	transformReceivedMsg func(receivedMsg any, expectedMsg any) any,
) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedMsg, err := unmarshalJson(b.Bytes())
	if err != nil {
		return err
	}

	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	return c.waitNextMessage(conn, expectedMsg, transformReceivedMsg)
}

func (c *WebsocketClient) waitNextMessage(
	conn *websocket.Conn,
	expectedMsg any,
	transformReceivedMsg func(receivedMsg any, expectedMsg any) any,
) error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		_, b, err := conn.ReadMessage()
		if err != nil {
			errCh <- err
			return
		}
		receivedMsg, err := unmarshalJson(b)
		if err != nil {
			errCh <- err
			return
		}

		diff := pretty.Compare(transformReceivedMsg(receivedMsg, expectedMsg), expectedMsg)
		if diff != "" {
			errCh <- fmt.Errorf("unexpected message: %s", diff)
			return
		}
	}()

	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()
	select {
	case err, ok := <-errCh:
		if !ok {
			return nil
		}

		return err
	case <-timer.C:
		return errors.New("reached timeout")
	}
}

func (c *WebsocketClient) iWaitMessageFromRoom(
	ctx context.Context,
	room, doc string,
	transformReceivedMsg func(receivedMsg any, expectedMsg any) any,
) error {
	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}
	expectedMsg, err := unmarshalJson(b.Bytes())
	if err != nil {
		return err
	}
	b, err = c.templater.Execute(ctx, room)
	if err != nil {
		return err
	}
	room = b.String()
	errCh := make(chan error, 1)
	msgsMx := sync.Mutex{}
	caughtMsgs := make([]any, 0)
	go func() {
		defer close(errCh)
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			msg, err := unmarshalJson(b)
			if err != nil {
				errCh <- err
				return
			}

			m, ok := msg.(map[string]any)
			if !ok {
				errCh <- fmt.Errorf("unknown message type %+v", m)
				return
			}
			t, ok := m["type"].(int64)
			if !ok {
				continue
			}

			if t == libwebsocket.WMessageFail {
				errCh <- fmt.Errorf("received fail message: %+v", m)
				return
			}

			receivedRoom, ok := m["room"].(string)
			if !ok {
				continue
			}
			receivedMsg, ok := m["msg"]
			if !ok {
				continue
			}

			if t == libwebsocket.WMessageSuccess && receivedRoom == room {
				diff := pretty.Compare(transformReceivedMsg(receivedMsg, expectedMsg), expectedMsg)
				if diff == "" {
					return
				}

				msgsMx.Lock()
				caughtMsgs = append(caughtMsgs, receivedMsg)
				msgsMx.Unlock()
			}
		}
	}()

	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()
	select {
	case err, ok := <-errCh:
		if !ok {
			return nil
		}

		return err
	case <-timer.C:
		msgsMx.Lock()
		defer msgsMx.Unlock()
		return fmt.Errorf("reached timeout: caught %d messages from %q room but none of them matches to expected message\n%s\n",
			len(caughtMsgs), room, pretty.Compare(caughtMsgs, []any{expectedMsg}))
	}
}
