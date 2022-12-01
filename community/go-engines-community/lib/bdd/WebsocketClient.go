package bdd

import (
	"context"
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
	if connIndex, ok := getWebsocketConn(ctx); ok {
		c.connsMx.Lock()
		defer c.connsMx.Unlock()
		conn := c.conns[connIndex]
		delete(c.conns, connIndex)
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(wsWriteWait))
		if err != nil {
			return ctx, err
		}
		err = conn.Close()
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (c *WebsocketClient) ISubscribeToRoom(ctx context.Context, room string) (context.Context, error) {
	connIndex, ok := getWebsocketConn(ctx)
	var conn *websocket.Conn
	if ok {
		c.connsMx.Lock()
		conn = c.conns[connIndex]
		c.connsMx.Unlock()
	} else {
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(c.apiURL, nil)
		if err != nil {
			return ctx, err
		}

		c.connsMx.Lock()
		c.lastConnIndex++
		connIndex = c.lastConnIndex
		c.conns[connIndex] = conn
		c.connsMx.Unlock()
		ctx = setWebsocketConn(ctx, connIndex)
		token, ok := getApiAuthToken(ctx)
		if ok {
			err = conn.WriteJSON(libwebsocket.RMessage{
				Type:  libwebsocket.RMessageAuth,
				Token: token,
			})
			if err != nil {
				return ctx, err
			}
		}
	}

	err := conn.WriteJSON(libwebsocket.RMessage{
		Type: libwebsocket.RMessageJoin,
		Room: room,
	})
	return ctx, err
}

func (c *WebsocketClient) IWaitMessageFromRoom(ctx context.Context, room, doc string) error {
	return c.waitMessageFromRoom(ctx, room, doc, func(receivedMsg interface{}, _ interface{}) interface{} {
		return receivedMsg
	})
}

func (c *WebsocketClient) IWaitMessageFromRoomWhichContains(ctx context.Context, room, doc string) error {
	return c.waitMessageFromRoom(ctx, room, doc, func(receivedMsg interface{}, expectedMsg interface{}) interface{} {
		return getPartialResponse(receivedMsg, expectedMsg)
	})
}

func (c *WebsocketClient) waitMessageFromRoom(
	ctx context.Context,
	room, doc string,
	transformReceivedMsg func(receivedMsg interface{}, expectedMsg interface{}) interface{},
) error {
	connIndex, ok := getWebsocketConn(ctx)
	if !ok {
		return fmt.Errorf("websocket connection is not open")
	}
	c.connsMx.Lock()
	conn := c.conns[connIndex]
	c.connsMx.Unlock()
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}
	expectedMsg, err := unmarshalJson(b.Bytes())
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	msgsMx := sync.Mutex{}
	caughtMsgs := make([]interface{}, 0)
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

			m, ok := msg.(map[string]interface{})
			if !ok {
				continue
			}
			t, ok := m["type"].(int64)
			if !ok {
				continue
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
	case err := <-errCh:
		return err
	case <-timer.C:
		msgsMx.Lock()
		defer msgsMx.Unlock()
		return fmt.Errorf("reached timeout: caught %d messages from %q room but none of them matches to expected message\n%s\n",
			len(caughtMsgs), room, pretty.Compare(caughtMsgs, []interface{}{expectedMsg}))
	}
}
