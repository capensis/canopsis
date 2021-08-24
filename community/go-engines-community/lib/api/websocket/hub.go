// Package websocket contains implementation of websocket.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/keymutex"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	RMessageJoin = iota
	RMessageLeave
)

const (
	WMessageSuccess = iota
	WMessageFail
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Hub interface is used to implement websocket room.
type Hub interface {
	// Start pings connections.
	Start(ctx context.Context)
	// Connect creates listener connection.
	Connect(userId string, w http.ResponseWriter, r *http.Request) error
	// Send sends message to all listeners in room.
	Send(room string, msg interface{})
}

func NewHub(upgrader Upgrader, authorizer Authorizer, logger zerolog.Logger) Hub {
	return &hub{
		upgrader:   upgrader,
		roomsMx:    keymutex.New(),
		rooms:      make(map[string][]string),
		conns:      make(map[string]Connection),
		authorizer: authorizer,
		logger:     logger,
	}
}

type RMessage struct {
	Type int    `json:"type"`
	Room string `json:"room"`
}

type WMessage struct {
	Type  int         `json:"type"`
	Room  string      `json:"room,omitempty"`
	Msg   interface{} `json:"msg,omitempty"`
	Error string      `json:"error,omitempty"`
}

type hub struct {
	upgrader   Upgrader
	roomsMx    keymutex.KeyMutex
	rooms      map[string][]string
	connsMx    sync.RWMutex
	conns      map[string]Connection
	authorizer Authorizer
	logger     zerolog.Logger
}

func (h *hub) Start(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			h.pingConnections()
		}
	}

	h.closeConnections()
}

func (h *hub) Connect(userId string, w http.ResponseWriter, r *http.Request) error {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		h.logger.Err(err).
			Str("addr", conn.RemoteAddr().String()).
			Msg("cannot set read deadline")
	}
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot set read deadline")
		}

		return nil
	})

	connId := utils.NewID()
	h.conns[connId] = conn

	// Run goroutine to listen connection.
	go h.listen(connId, userId, conn)

	return nil
}

func (h *hub) Send(room string, b interface{}) {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)

	msg := WMessage{
		Type: WMessageSuccess,
		Room: room,
		Msg:  b,
	}
	closedConns := make([]string, 0)

	for _, connId := range h.rooms[room] {
		conn := h.conns[connId]
		err := conn.WriteJSON(msg)
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("room", room).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
	}

	h.roomsMx.Unlock(room)
	h.connsMx.RUnlock()

	for _, connId := range closedConns {
		h.closeConnection(connId)
	}
}

func (h *hub) join(connId, room string) error {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	if _, ok := h.conns[connId]; !ok {
		return fmt.Errorf("connection not found")
	}

	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	for _, v := range h.rooms[room] {
		if v == connId {
			return fmt.Errorf("connection has already joined to room")
		}
	}

	if len(h.rooms[room]) == 0 {
		h.rooms[room] = make([]string, 0)
	}

	h.rooms[room] = append(h.rooms[room], connId)

	return nil
}

func (h *hub) leave(connId, room string) error {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	if _, ok := h.conns[connId]; !ok {
		return fmt.Errorf("connection not found")
	}

	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	index := -1
	for i, v := range h.rooms[room] {
		if v == connId {
			index = i
			break
		}
	}

	if index < 0 {
		return fmt.Errorf("connection hasn't joined to room")
	}

	h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
	return nil
}

func (h *hub) closeConnections() {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	for room := range h.rooms {
		h.roomsMx.Lock(room)
		h.rooms[room] = nil
		h.roomsMx.Unlock(room)
	}

	for _, conn := range h.conns {
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot close connection")
		}
	}

	h.conns = nil
}

func (h *hub) removeConnection(connId string) {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	for room, conns := range h.rooms {
		h.roomsMx.Lock(room)

		index := -1
		for i, v := range conns {
			if v == connId {
				index = i
				break
			}
		}

		if index >= 0 {
			h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
		}

		h.roomsMx.Unlock(room)
	}

	delete(h.conns, connId)
}

func (h *hub) closeConnection(connId string) {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	for room, conns := range h.rooms {
		h.roomsMx.Lock(room)

		index := -1
		for i, v := range conns {
			if v == connId {
				index = i
				break
			}
		}

		if index >= 0 {
			h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
		}

		h.roomsMx.Unlock(room)
	}

	if conn, ok := h.conns[connId]; ok {
		err := conn.Close()
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("connection close failed")
		}

		delete(h.conns, connId)
	}
}

func (h *hub) pingConnections() {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	conns := make(map[string]Connection, len(h.conns))
	for id, conn := range h.conns {
		err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
		if err == nil {
			conns[id] = conn
		} else {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot ping connection, connection will be closed")
			err = conn.Close()
			if err != nil {
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("connection close failed")
			}
		}
	}

	h.conns = conns
}

func (h *hub) listen(connId, userId string, conn Connection) {
	for {
		msg := RMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			syntaxErr := &json.SyntaxError{}
			if errors.As(err, &syntaxErr) {
				if !h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "invalid message"}) {
					return
				}

				continue
			}

			closeErr := &websocket.CloseError{}
			if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseNormalClosure {
				h.logger.
					Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("connection closed unexpectedly")
			}

			h.removeConnection(connId)
			return
		}

		if msg.Room == "" {
			if !h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "room is missing"}) {
				return
			}

			continue
		}

		errMsg := WMessage{
			Type: WMessageFail,
			Room: msg.Room,
		}

		switch msg.Type {
		case RMessageJoin:
			ok, err := h.authorizer.Auth(userId, msg.Room)
			if err != nil {
				h.logger.Err(err).Msg("cannot authorize user")
			}

			if err != nil || !ok {
				errMsg.Error = "cannot authorize user"
				if !h.sendToConn(connId, errMsg) {
					return
				}
				continue
			}

			err = h.join(connId, msg.Room)
			if err != nil {
				errMsg.Error = err.Error()
				if !h.sendToConn(connId, errMsg) {
					return
				}
			}
		case RMessageLeave:
			err := h.leave(connId, msg.Room)
			if err != nil {
				errMsg.Error = err.Error()
				if !h.sendToConn(connId, errMsg) {
					return
				}
			}
		default:
			errMsg.Error = "unknown message type"
			if !h.sendToConn(connId, errMsg) {
				return
			}
		}
	}
}

func (h *hub) sendToConn(connId string, msg interface{}) bool {
	h.connsMx.RLock()
	conn := h.conns[connId]
	err := conn.WriteJSON(msg)
	if err != nil {
		h.logger.Err(err).
			Str("addr", conn.RemoteAddr().String()).
			Msg("cannot write message to connection, connection will be closed")
		h.connsMx.RUnlock()
		h.closeConnection(connId)
		return false
	}

	h.connsMx.RUnlock()
	return true
}
