// Package websocket contains implementation of websocket.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/keymutex"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	RMessageClientPing = iota
	RMessageJoin
	RMessageLeave
	RMessageAuth
)

const (
	WMessageClientPong = iota
	WMessageSuccess
	WMessageFail
	WMessageCloseRoom
)

const (
	errAuthFailed          = "cannot authorize user"
	errUnknownRMessageType = "unknown message type"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// Hub interface is used to implement websocket room.
type Hub interface {
	// Start pings connections.
	Start(ctx context.Context)
	// Connect creates listener connection.
	Connect(w http.ResponseWriter, r *http.Request) error
	// Send sends message to all listeners in room.
	Send(room string, msg interface{})
	// RegisterRoom adds room with permissions.
	RegisterRoom(room string, perms ...string) error
	// CloseRoom removes room.
	CloseRoom(room string) error
	CloseRoomAndNotify(room string) error
	// RoomHasConnection returns true if room listens at least one connection.
	RoomHasConnection(room string) bool
	// GetUsers returns users from all connections with their tokens.
	GetUsers() map[string][]string
}

func NewHub(
	upgrader Upgrader,
	authorizer Authorizer,
	pingInterval time.Duration,
	logger zerolog.Logger,
) Hub {
	return &hub{
		upgrader:     upgrader,
		roomsMx:      keymutex.New(),
		rooms:        make(map[string][]string),
		conns:        make(map[string]userConn),
		authorizer:   authorizer,
		pingInterval: pingInterval,
		pongInterval: pingInterval * 10 / 9,
		logger:       logger,
	}
}

type RMessage struct {
	Type  int    `json:"type"`
	Room  string `json:"room"`
	Token string `json:"token"`
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
	conns      map[string]userConn
	authorizer Authorizer
	// Send pings to peer with this period. Must be less than pongInterval.
	pingInterval time.Duration
	// Time allowed to read the next pong message from the peer.
	pongInterval time.Duration
	logger       zerolog.Logger
}

type userConn struct {
	userId, token string

	conn Connection
}

func (h *hub) Start(ctx context.Context) {
	ticker := time.NewTicker(h.pingInterval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			closedConns := h.pingConnections()
			h.closeConnections(closedConns...)

			closedConns, connsToDisconnect := h.checkAuth(ctx)
			h.closeConnections(closedConns...)
			h.disconnectConnections(connsToDisconnect...)
		}
	}

	h.stop()
}

func (h *hub) Connect(w http.ResponseWriter, r *http.Request) error {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	err = conn.SetReadDeadline(time.Now().Add(h.pongInterval))
	if err != nil {
		h.logger.Err(err).
			Str("addr", conn.RemoteAddr().String()).
			Msg("cannot set read deadline")
	}
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(h.pongInterval))
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot set read deadline")
		}

		return nil
	})

	connId := utils.NewID()
	h.conns[connId] = userConn{
		conn: conn,
	}

	// Run goroutine to listen connection.
	go h.listen(connId, conn)

	return nil
}

func (h *hub) Send(room string, b interface{}) {
	closedConns := h.sendAndCheckConn(room, b)
	h.closeConnections(closedConns...)
}

func (h *hub) sendAndCheckConn(room string, b interface{}) []string {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
		h.connsMx.RUnlock()
	}()

	msg := WMessage{
		Type: WMessageSuccess,
		Room: room,
		Msg:  b,
	}
	closedConns := make([]string, 0)

	for _, connId := range h.rooms[room] {
		conn := h.conns[connId].conn
		err := conn.WriteJSON(msg)
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("room", room).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
	}

	return closedConns
}

func (h *hub) RegisterRoom(room string, perms ...string) error {
	return h.authorizer.AddRoom(room, perms)
}

func (h *hub) CloseRoom(room string) error {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
		h.connsMx.RUnlock()
	}()

	err := h.authorizer.RemoveRoom(room)
	if err != nil {
		return err
	}
	delete(h.rooms, room)

	return nil
}

func (h *hub) CloseRoomAndNotify(room string) error {
	closedConns, err := h.closeRoomAndCheckConn(room)
	if err != nil {
		return err
	}

	h.closeConnections(closedConns...)

	return nil
}

func (h *hub) closeRoomAndCheckConn(room string) ([]string, error) {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
		h.connsMx.RUnlock()
	}()

	msg := WMessage{
		Type: WMessageCloseRoom,
		Room: room,
	}

	closedConns := make([]string, 0)
	for _, connId := range h.rooms[room] {
		conn := h.conns[connId].conn
		err := conn.WriteJSON(msg)
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("room", room).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
	}

	err := h.authorizer.RemoveRoom(room)
	if err != nil {
		return nil, err
	}
	delete(h.rooms, room)

	return closedConns, nil
}

func (h *hub) RoomHasConnection(room string) bool {
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
	}()

	return len(h.rooms[room]) > 0
}

func (h *hub) join(connId, room string) (closed bool) {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
		h.connsMx.RUnlock()
	}()

	c := h.conns[connId]
	userId := c.userId
	conn := c.conn

	for _, v := range h.rooms[room] {
		if v == connId {
			return
		}
	}

	ok, err := h.authorizer.Authorize(userId, room)
	if err != nil {
		h.logger.Err(err).Msg(errAuthFailed)
		return
	}

	if !ok {
		err := conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Room:  room,
			Error: errAuthFailed,
		})
		if err != nil {
			closed = true
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
		return
	}

	if len(h.rooms[room]) == 0 {
		h.rooms[room] = make([]string, 0)
	}

	h.rooms[room] = append(h.rooms[room], connId)
	return
}

func (h *hub) leave(connId, room string) (closed bool) {
	h.connsMx.RLock()
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
		h.connsMx.RUnlock()
	}()

	index := -1
	for i, v := range h.rooms[room] {
		if v == connId {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
	return
}

func (h *hub) stop() {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	for room := range h.rooms {
		h.cleanRoom(room)
	}

	for _, c := range h.conns {
		conn := c.conn
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot close connection")
		}
	}

	h.conns = nil
}

func (h *hub) cleanRoom(room string) {
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
	}()
	h.rooms[room] = nil
}

func (h *hub) removeConnections(connIds ...string) {
	if len(connIds) == 0 {
		return
	}

	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	h.removeConnsFromRooms(connIds)

	for _, connId := range connIds {
		delete(h.conns, connId)
	}
}

func (h *hub) closeConnections(connIds ...string) {
	if len(connIds) == 0 {
		return
	}

	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	h.removeConnsFromRooms(connIds)

	for _, connId := range connIds {
		if c, ok := h.conns[connId]; ok {
			conn := c.conn
			err := conn.Close()
			if err != nil {
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("connection close failed")
			}

			delete(h.conns, connId)
		}
	}
}

func (h *hub) disconnectConnections(connIds ...string) {
	if len(connIds) == 0 {
		return
	}

	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	h.removeConnsFromRooms(connIds)

	for _, connId := range connIds {
		if c, ok := h.conns[connId]; ok {
			conn := c.conn
			err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
			if err != nil {
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("connection close failed")
			}

			delete(h.conns, connId)
		}
	}
}

func (h *hub) pingConnections() []string {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	closedConns := make([]string, 0)
	for connId, c := range h.conns {
		conn := c.conn
		err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot ping connection, connection will be closed")
		}
	}

	return closedConns
}

func (h *hub) checkAuth(ctx context.Context) ([]string, []string) {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	closedConns := make([]string, 0)
	connsToDisconnect := make([]string, 0)
	checked := make(map[string]bool, len(h.conns))

	for connId, c := range h.conns {
		if c.userId == "" {
			checked[connId] = true
			continue
		}

		conn := c.conn
		userId, err := h.authorizer.Authenticate(ctx, c.token)
		if err == nil {
			if userId != "" {
				checked[connId] = true
				continue
			}

			connsToDisconnect = append(connsToDisconnect, connId)
			h.logger.Error().
				Str("addr", conn.RemoteAddr().String()).
				Str("user", c.userId).
				Msg("cannot found user, connection will be closed")

			continue
		}

		h.logger.Err(err).Msg(errAuthFailed)
		err = conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Error: errAuthFailed,
		})
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
	}

	for room := range h.rooms {
		closed := h.checkRoomAuth(room, checked)
		if len(closed) > 0 {
			closedConns = append(closedConns, closed...)
		}
	}

	return closedConns, connsToDisconnect
}

func (h *hub) checkRoomAuth(room string, checked map[string]bool) []string {
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
	}()

	roomConns := h.rooms[room]
	authRoomConns := make([]string, 0, len(roomConns))
	closedConns := make([]string, 0)

	for _, connId := range roomConns {
		if !checked[connId] {
			continue
		}

		c := h.conns[connId]
		conn := c.conn
		userId := c.userId
		ok, err := h.authorizer.Authorize(userId, room)
		if err == nil && ok {
			authRoomConns = append(authRoomConns, connId)
			continue
		}

		if err != nil {
			h.logger.Err(err).Msg(errAuthFailed)
		}

		err = conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Room:  room,
			Error: errAuthFailed,
		})
		if err != nil {
			closedConns = append(closedConns, connId)
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
	}

	h.rooms[room] = authRoomConns

	return closedConns
}

func (h *hub) listen(connId string, conn Connection) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	closed := false

	for {
		if closed {
			h.closeConnections(connId)
			return
		}

		msg := RMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			unmarshalErr := &json.UnmarshalTypeError{}
			syntaxErr := &json.SyntaxError{}
			if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
				closed = h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "invalid message"})
				continue
			}

			closeErr := &websocket.CloseError{}
			if errors.As(err, &closeErr) {
				if closeErr.Code != websocket.CloseNormalClosure {
					h.logger.
						Err(err).
						Str("addr", conn.RemoteAddr().String()).
						Msg("connection closed unexpectedly")
				}

				h.removeConnections(connId)
			} else {
				h.logger.
					Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot read message from connection, connection will be closed")
				h.disconnectConnections(connId)
			}

			return
		}

		switch msg.Type {
		case RMessageClientPing:
			closed = h.sendToConn(connId, WMessage{Type: WMessageClientPong})
		case RMessageAuth:
			if msg.Token == "" {
				closed = h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "token is missing"})
				continue
			}

			userId, err := h.authorizer.Authenticate(ctx, msg.Token)
			if err != nil {
				h.logger.Err(err).Msg("authentication failed")
			}
			if userId == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: "authentication failed",
				})
			} else {
				h.setConnAuth(connId, userId, msg.Token)
			}
		case RMessageJoin:
			if msg.Room == "" {
				closed = h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "room is missing"})
				continue
			}
			closed = h.join(connId, msg.Room)
		case RMessageLeave:
			if msg.Room == "" {
				closed = h.sendToConn(connId, WMessage{Type: WMessageFail, Error: "room is missing"})
				continue
			}
			closed = h.leave(connId, msg.Room)
		default:
			closed = h.sendToConn(connId, WMessage{
				Type:  WMessageFail,
				Room:  msg.Room,
				Error: errUnknownRMessageType,
			})
		}
	}
}

func (h *hub) setConnAuth(connId, userId, token string) {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	h.conns[connId] = userConn{
		userId: userId,
		token:  token,
		conn:   h.conns[connId].conn,
	}
}

func (h *hub) sendToConn(connId string, msg interface{}) (closed bool) {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	conn := h.conns[connId].conn
	err := conn.WriteJSON(msg)
	if err != nil {
		closed = true
		h.logger.Err(err).
			Str("addr", conn.RemoteAddr().String()).
			Msg("cannot write message to connection, connection will be closed")
	}

	return
}

func (h *hub) removeConnsFromRooms(connIds []string) {
	for room := range h.rooms {
		h.removeConnsFromRoom(room, connIds)
	}
}

func (h *hub) removeConnsFromRoom(room string, connIds []string) {
	h.roomsMx.Lock(room)
	defer func() {
		if err := h.roomsMx.Unlock(room); err != nil {
			h.logger.Err(err).Msg("roomsMx unlock")
		}
	}()

	conns := h.rooms[room]
	filteredConns := make([]string, 0, len(conns))

	for _, v := range conns {
		found := false
		for _, connId := range connIds {
			if v == connId {
				found = true
				break
			}
		}
		if !found {
			filteredConns = append(filteredConns, v)
		}
	}

	h.rooms[room] = filteredConns
}

func (h *hub) GetUsers() map[string][]string {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	users := make(map[string][]string, 0)
	for _, conn := range h.conns {
		if conn.userId != "" {
			if _, ok := users[conn.userId]; !ok {
				users[conn.userId] = make([]string, 0, 1)
			}
			users[conn.userId] = append(users[conn.userId], conn.token)
		}
	}

	return users
}
