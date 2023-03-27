// Package websocket contains implementation of websocket.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

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
	WMessageAuthSuccess
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
	Send(room string, msg any)
	// RegisterRoom adds room with permissions.
	RegisterRoom(room string, perms ...string) error
	// CloseRoom removes room.
	CloseRoom(room string) error
	CloseRoomAndNotify(room string) error
	// RoomHasConnection returns true if room listens at least one connection.
	RoomHasConnection(room string) bool
	// GetUsers returns users from all connections with their tokens.
	GetUsers() map[string][]string
	GetAuthConnectionsCount() int
}

func NewHub(
	upgrader Upgrader,
	authorizer Authorizer,
	pingInterval time.Duration,
	logger zerolog.Logger,
) Hub {
	return &hub{
		upgrader:     upgrader,
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
	Type  int    `json:"type"`
	Room  string `json:"room,omitempty"`
	Msg   any    `json:"msg,omitempty"`
	Error int    `json:"error,omitempty"`
}

type hub struct {
	upgrader   Upgrader
	roomsMx    sync.RWMutex
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
			updated := false
			closedConns := h.pingConnections()
			h.closeConnections(closedConns...)
			if len(closedConns) > 0 {
				updated = true
			}

			closedConns, connsToDisconnect := h.checkAuth(ctx)
			h.closeConnections(closedConns...)
			h.disconnectConnections(connsToDisconnect...)
			if len(closedConns) > 0 || len(connsToDisconnect) > 0 {
				updated = true
			}

			if updated {
				h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
			}
		}
	}

	h.stop()
}

func (h *hub) Connect(w http.ResponseWriter, r *http.Request) error {
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

	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	connId := utils.NewID()
	h.conns[connId] = userConn{
		conn: conn,
	}

	// Run goroutine to listen connection.
	go h.listen(connId, conn)

	return nil
}

func (h *hub) Send(room string, b any) {
	closedConns := h.sendToRoom(room, WMessage{
		Type: WMessageSuccess,
		Room: room,
		Msg:  b,
	})
	if len(closedConns) > 0 {
		h.closeConnections(closedConns...)
		h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
	}
}

func (h *hub) RegisterRoom(room string, perms ...string) error {
	defer func() {
		h.logger.Debug().Str("room", room).Msg("register websocket room")
	}()
	return h.authorizer.AddRoom(room, perms)
}

func (h *hub) CloseRoom(room string) error {
	defer func() {
		h.logger.Debug().Str("room", room).Msg("close websocket room")
	}()
	err := h.authorizer.RemoveRoom(room)
	if err != nil {
		return err
	}

	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()

	delete(h.rooms, room)

	return nil
}

func (h *hub) CloseRoomAndNotify(room string) error {
	closedConns := h.sendToRoom(room, WMessage{
		Type: WMessageCloseRoom,
		Room: room,
	})
	if len(closedConns) > 0 {
		h.closeConnections(closedConns...)
		h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
	}

	return h.CloseRoom(room)
}

func (h *hub) RoomHasConnection(room string) bool {
	h.roomsMx.RLock()
	defer h.roomsMx.RUnlock()

	return len(h.rooms[room]) > 0
}

func (h *hub) GetUsers() map[string][]string {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	users := make(map[string][]string)
	uniqueTokens := make(map[string]struct{}, len(h.conns))
	for _, conn := range h.conns {
		if conn.userId != "" {
			if _, ok := users[conn.userId]; !ok {
				users[conn.userId] = make([]string, 0, 1)
			}
			if _, ok := uniqueTokens[conn.token]; !ok {
				users[conn.userId] = append(users[conn.userId], conn.token)
				uniqueTokens[conn.token] = struct{}{}
			}
		}
	}

	return users
}

func (h *hub) GetAuthConnectionsCount() int {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	uniqueTokens := make(map[string]struct{})
	for _, conn := range h.conns {
		if conn.token != "" {
			if _, ok := uniqueTokens[conn.token]; !ok {
				uniqueTokens[conn.token] = struct{}{}
			}
		}
	}

	return len(uniqueTokens)
}

func (h *hub) join(connId, room string) (closed bool) {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	c := h.conns[connId]
	userId := c.userId
	conn := c.conn

	ok := h.authorizer.HasRoom(room)
	if !ok {
		err := conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Room:  room,
			Error: http.StatusNotFound,
		})
		if err != nil {
			closed = true
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
		return
	}

	ok, err := h.authorizer.Authorize(userId, room)
	if err != nil {
		h.logger.Err(err).Msg("cannot authorize user")

		err := conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Room:  room,
			Error: http.StatusInternalServerError,
		})
		if err != nil {
			closed = true
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
		return
	}

	if !ok {
		code := http.StatusForbidden
		if userId == "" {
			code = http.StatusUnauthorized
		}
		err := conn.WriteJSON(WMessage{
			Type:  WMessageFail,
			Room:  room,
			Error: code,
		})
		if err != nil {
			closed = true
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}
		return
	}

	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()

	for _, v := range h.rooms[room] {
		if v == connId {
			return
		}
	}

	if len(h.rooms[room]) == 0 {
		h.rooms[room] = make([]string, 0, 1)
	}

	h.rooms[room] = append(h.rooms[room], connId)
	return
}

func (h *hub) leave(connId, room string) (closed bool) {
	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()

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
	h.roomsMx.Lock()
	defer func() {
		h.roomsMx.Unlock()
		h.connsMx.Unlock()
	}()

	h.rooms = make(map[string][]string)

	for _, c := range h.conns {
		conn := c.conn
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
		if err != nil {
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot close connection")
		}
	}

	h.conns = make(map[string]userConn)
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
		if err != nil {
			h.logger.Err(err).Msg("cannot authorize user")
			err = conn.WriteJSON(WMessage{
				Type:  WMessageFail,
				Error: http.StatusInternalServerError,
			})
			if err != nil {
				closedConns = append(closedConns, connId)
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot write message to connection, connection will be closed")
			}

			continue
		}

		if userId == "" {
			connsToDisconnect = append(connsToDisconnect, connId)
			h.logger.Error().
				Str("addr", conn.RemoteAddr().String()).
				Str("user", c.userId).
				Msg("cannot found user, connection will be closed")

			err = conn.WriteJSON(WMessage{
				Type:  WMessageFail,
				Error: http.StatusUnauthorized,
			})
			if err != nil {
				closedConns = append(closedConns, connId)
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot write message to connection, connection will be closed")
			}

			continue
		}

		checked[connId] = true
	}

	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()
	for room := range h.rooms {
		closed := h.checkRoomAuth(room, checked)
		if len(closed) > 0 {
			closedConns = append(closedConns, closed...)
		}
	}

	return closedConns, connsToDisconnect
}

func (h *hub) checkRoomAuth(room string, checked map[string]bool) []string {
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
		if err != nil {
			h.logger.Err(err).Msg("cannot authorize user")

			err = conn.WriteJSON(WMessage{
				Type:  WMessageFail,
				Room:  room,
				Error: http.StatusInternalServerError,
			})
			if err != nil {
				closedConns = append(closedConns, connId)
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot write message to connection, connection will be closed")
			}

			continue
		}

		if !ok {
			err = conn.WriteJSON(WMessage{
				Type:  WMessageFail,
				Room:  room,
				Error: http.StatusForbidden,
			})
			if err != nil {
				closedConns = append(closedConns, connId)
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot write message to connection, connection will be closed")
			}
			continue
		}

		authRoomConns = append(authRoomConns, connId)
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
			h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
			return
		}

		msg := RMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			unmarshalErr := &json.UnmarshalTypeError{}
			syntaxErr := &json.SyntaxError{}
			if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusBadRequest,
					Msg:   "cannot parse JSON",
				})
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

			h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
			return
		}

		h.logger.Debug().
			Str("conn", connId).
			Int("type", msg.Type).
			Str("room", msg.Room).
			Msg("read websocket conn")

		switch msg.Type {
		case RMessageClientPing:
			closed = h.sendToConn(connId, WMessage{Type: WMessageClientPong})
		case RMessageAuth:
			if msg.Token == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusBadRequest,
					Msg:   "token is missing",
				})
				continue
			}

			userId, err := h.authorizer.Authenticate(ctx, msg.Token)
			if err != nil {
				h.logger.Err(err).Msg("authentication failed")
			}
			if userId == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusUnauthorized,
				})
			} else {
				h.setConnAuth(connId, userId, msg.Token)
				closed = h.sendToConn(connId, WMessage{Type: WMessageAuthSuccess})
				h.Send(RoomLoggedUserCount, h.GetAuthConnectionsCount())
			}
		case RMessageJoin:
			if msg.Room == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusBadRequest,
					Msg:   "room is missing",
				})
				continue
			}
			closed = h.join(connId, msg.Room)
		case RMessageLeave:
			if msg.Room == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusBadRequest,
					Msg:   "room is missing",
				})
				continue
			}
			closed = h.leave(connId, msg.Room)
		default:
			closed = h.sendToConn(connId, WMessage{
				Type:  WMessageFail,
				Room:  msg.Room,
				Error: http.StatusBadRequest,
				Msg:   "unknown message type",
			})
		}
	}
}

func (h *hub) setConnAuth(connId, userId, token string) {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	h.conns[connId] = userConn{
		userId: userId,
		token:  token,
		conn:   h.conns[connId].conn,
	}
}

func (h *hub) sendToConn(connId string, msg WMessage) (closed bool) {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	defer func() {
		h.logger.Debug().
			Int("type", msg.Type).
			Str("conn", connId).
			Bool("closed", closed).
			Msg("send to webhook conn")
	}()

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

func (h *hub) sendToRoom(room string, msg WMessage) []string {
	h.connsMx.RLock()
	h.roomsMx.RLock()
	defer func() {
		h.roomsMx.RUnlock()
		h.connsMx.RUnlock()
	}()

	closedConns := make([]string, 0)
	count := len(h.rooms[room])
	defer func() {
		if count > 0 {
			h.logger.Debug().
				Str("room", room).
				Int("type", msg.Type).
				Int("conns", count).
				Int("closed", len(closedConns)).
				Msg("send to webhook room")
		}
	}()
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

func (h *hub) removeConnsFromRooms(connIds []string) {
	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()

	for room, conns := range h.rooms {
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
}
