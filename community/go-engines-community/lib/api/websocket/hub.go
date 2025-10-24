// Package websocket contains implementation of websocket.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
	SendGroupRoom(group, id string, msg any)
	SendGroupRoomByConnections(connIds []string, group, id string, b any)
	// RegisterRoom adds room with permissions.
	RegisterRoom(room string, perms ...string) error
	RegisterGroup(group string, params GroupParameters, perms ...string) error
	GetGroupIds(group string) []string
	GetConnectedGroupIds(group string) []string
	CloseGroupRoom(group, id string) error
	CloseGroupRoomAndNotify(group, id string) error
	GetUserTokens() []string
	GetConnections() []UserConnection
}

type GroupParameters struct {
	CheckExists GroupCheckExists
	OnNotExist  GroupOnNotExist
	OnJoin      GroupOnJoin
	OnLeave     GroupOnLeave
}

func (p GroupParameters) IsZero() bool {
	return p.CheckExists == nil && p.OnJoin == nil && p.OnLeave == nil
}

type GroupCheckExists func(ctx context.Context, id string) (bool, error)
type GroupOnNotExist func(ctx context.Context, id string) (any, error)
type GroupOnJoin func(ctx context.Context, connId, userID, roomId string, data any) error
type GroupOnLeave func(connId, roomId string) error

func NewHub(
	ctx context.Context,
	upgrader Upgrader,
	authorizer Authorizer,
	pingInterval time.Duration,
	logger zerolog.Logger,
) Hub {
	return &hub{
		hubCtx:       ctx,
		upgrader:     upgrader,
		rooms:        make(map[string][]string),
		conns:        make(map[string]userConn),
		groups:       make(map[string]GroupParameters),
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
	Data  any    `json:"data"`
}

type WMessage struct {
	Type  int    `json:"type"`
	Room  string `json:"room,omitempty"`
	Msg   any    `json:"msg,omitempty"`
	Error int    `json:"error,omitempty"`
}

type UserConnection struct {
	ID     string
	UserID string
	Token  string
}

type hub struct {
	// hubCtx should be used for certain actions, which should be outside of user scope, e.g. group rooms.
	hubCtx context.Context //nolint:containedctx

	upgrader   Upgrader
	roomsMx    sync.RWMutex
	rooms      map[string][]string
	connsMx    sync.RWMutex
	conns      map[string]userConn
	groupsMx   sync.RWMutex
	groups     map[string]GroupParameters
	authorizer Authorizer
	// Send pings to peer with this period. Must be less than pongInterval.
	pingInterval time.Duration
	// Time allowed to read the next pong message from the peer.
	pongInterval time.Duration
	logger       zerolog.Logger
}

type userConn struct {
	userID, token string

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
	}
}

func (h *hub) SendGroupRoom(group, id string, b any) {
	h.Send(group+id, b)
}

func (h *hub) SendGroupRoomByConnections(connIds []string, group, id string, b any) {
	room := group + id
	closedConns := h.sendToRoomByConns(room, connIds, WMessage{
		Type: WMessageSuccess,
		Room: room,
		Msg:  b,
	})
	if len(closedConns) > 0 {
		h.closeConnections(closedConns...)
	}
}

func (h *hub) RegisterRoom(room string, perms ...string) error {
	defer h.logger.Debug().Str("room", room).Msg("register websocket room")
	return h.authorizer.AddRoom(room, perms)
}

func (h *hub) RegisterGroup(group string, params GroupParameters, perms ...string) error {
	defer h.logger.Debug().Str("group", group).Msg("register websocket group")
	if !params.IsZero() {
		h.groupsMx.Lock()
		h.groups[group] = params
		h.groupsMx.Unlock()
	}
	return h.authorizer.AddGroup(group, perms, params.CheckExists)
}

func (h *hub) GetGroupIds(group string) []string {
	return h.authorizer.GetGroupIds(group)
}

func (h *hub) GetConnectedGroupIds(group string) []string {
	ids := h.GetGroupIds(group)
	if len(ids) == 0 {
		return nil
	}

	h.roomsMx.RLock()
	defer h.roomsMx.RUnlock()

	k := 0
	for _, id := range ids {
		room := group + id
		if len(h.rooms[room]) > 0 {
			ids[k] = id
			k++
		}
	}

	return ids[:k]
}

func (h *hub) CloseGroupRoom(group, id string) error {
	room := group + id
	defer h.logger.Debug().Str("room", room).Msg("close websocket room")

	err := h.authorizer.RemoveGroupRoom(group, id)
	if err != nil {
		return err
	}

	onLeave := h.getOnLeaveByGroup(group)
	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()

	if onLeave != nil {
		for _, connId := range h.rooms[room] {
			err = onLeave(connId, id)
			if err != nil {
				h.logger.Err(err).Str("room", room).Msgf("cannot leave room")
			}
		}
	}

	delete(h.rooms, room)

	return nil
}

func (h *hub) CloseGroupRoomAndNotify(group, id string) error {
	room := group + id
	closedConns := h.sendToRoom(room, WMessage{
		Type: WMessageCloseRoom,
		Room: room,
	})
	if len(closedConns) > 0 {
		h.closeConnections(closedConns...)
	}

	return h.CloseGroupRoom(group, id)
}

func (h *hub) GetUserTokens() []string {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	tokens := make([]string, 0)
	uniqueTokens := make(map[string]struct{}, len(h.conns))
	for _, conn := range h.conns {
		if conn.token != "" {
			if _, ok := uniqueTokens[conn.token]; !ok {
				tokens = append(tokens, conn.token)
				uniqueTokens[conn.token] = struct{}{}
			}
		}
	}

	return tokens
}

func (h *hub) GetConnections() []UserConnection {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	conns := make([]UserConnection, 0, len(h.conns))
	for connId, conn := range h.conns {
		if conn.userID != "" {
			conns = append(conns, UserConnection{
				ID:     connId,
				UserID: conn.userID,
				Token:  conn.token,
			})
		}
	}

	return conns
}

func (h *hub) join(ctx context.Context, connId, room string, data any) bool {
	h.connsMx.RLock()
	defer h.connsMx.RUnlock()

	c := h.conns[connId]
	userID := c.userID
	conn := c.conn
	closed := false
	granted, msg, err := h.authorizeOnJoin(ctx, userID, room)
	if err != nil {
		if errors.Is(err, ErrNotFoundRoom) || errors.Is(err, ErrNotFoundRoomInGroup) {
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

			return closed
		}

		h.logger.Err(err).Msg("cannot authorize user")
		err = conn.WriteJSON(WMessage{
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

		return closed
	}

	if msg != nil {
		err = conn.WriteJSON(WMessage{
			Type: WMessageSuccess,
			Room: room,
			Msg:  msg,
		})
		if err != nil {
			closed = true
			h.logger.Err(err).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection, connection will be closed")
		}

		return closed
	}

	if !granted {
		code := http.StatusForbidden
		if userID == "" {
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

		return closed
	}

	onJoin, id := h.getOnJoin(room)
	h.roomsMx.Lock()
	defer h.roomsMx.Unlock()
	for _, v := range h.rooms[room] {
		if v == connId {
			return false
		}
	}

	if onJoin != nil {
		err := onJoin(h.hubCtx, connId, userID, id, data) //nolint:contextcheck
		if err != nil {
			err = conn.WriteJSON(WMessage{
				Type:  WMessageFail,
				Room:  room,
				Error: http.StatusInternalServerError,
			})
			if err != nil {
				closed = true
				h.logger.Err(err).
					Str("addr", conn.RemoteAddr().String()).
					Msg("cannot join to room, connection will be closed")
			}

			return closed
		}
	}

	h.rooms[room] = append(h.rooms[room], connId)
	return false
}

func (h *hub) authorizeOnJoin(ctx context.Context, userID, room string) (bool, any, error) {
	ok, err := h.authorizer.Authorize(ctx, userID, room)
	if err == nil {
		return ok, nil, nil
	}

	if errors.Is(err, ErrNotFoundRoom) {
		return false, nil, err
	}

	if errors.Is(err, ErrNotFoundRoomInGroup) {
		onNotExist, id := h.getOnNotExist(room)
		if onNotExist == nil {
			return false, nil, err
		}

		msg, notExistErr := onNotExist(ctx, id)
		if notExistErr != nil {
			return false, nil, notExistErr
		}

		if msg == nil {
			return false, nil, err
		}

		return false, msg, nil
	}

	return false, nil, err
}

func (h *hub) leave(connId, room string) {
	h.roomsMx.Lock()
	index := -1
	for i, v := range h.rooms[room] {
		if v == connId {
			index = i
			break
		}
	}

	if index >= 0 {
		h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
	}

	h.roomsMx.Unlock()
	onLeave, id := h.getOnLeave(room)
	if onLeave != nil {
		err := onLeave(connId, id)
		if err != nil {
			h.logger.Err(err).Str("room", room).Msg("cannot leave room")
		}
	}
}

func (h *hub) stop() {
	h.connsMx.Lock()
	h.roomsMx.Lock()
	defer func() {
		h.roomsMx.Unlock()
		h.connsMx.Unlock()
	}()

	for room, connIds := range h.rooms {
		onLeave, id := h.getOnLeave(room)
		if onLeave != nil {
			for _, connId := range connIds {
				err := onLeave(connId, id)
				if err != nil {
					h.logger.Err(err).Str("room", room).Msgf("cannot leave room")
				}
			}
		}
	}

	h.rooms = make(map[string][]string)

	for _, c := range h.conns {
		conn := c.conn
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
		if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
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
			if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
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
		if c.userID == "" {
			checked[connId] = true
			continue
		}

		conn := c.conn
		userID, err := h.authorizer.Authenticate(ctx, c.token)
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

		if userID == "" {
			connsToDisconnect = append(connsToDisconnect, connId)
			h.logger.Error().
				Str("addr", conn.RemoteAddr().String()).
				Str("user", c.userID).
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
		closed := h.checkRoomAuth(ctx, room, checked)
		if len(closed) > 0 {
			closedConns = append(closedConns, closed...)
		}
	}

	return closedConns, connsToDisconnect
}

func (h *hub) checkRoomAuth(ctx context.Context, room string, checked map[string]bool) []string {
	roomConns := h.rooms[room]
	authRoomConns := make([]string, 0, len(roomConns))
	closedConns := make([]string, 0)
	onLeave, id := h.getOnLeave(room)

	for _, connId := range roomConns {
		if !checked[connId] {
			continue
		}

		c := h.conns[connId]
		conn := c.conn
		userID := c.userID
		ok, err := h.authorizer.Authorize(ctx, userID, room)
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

			if onLeave != nil {
				err = onLeave(connId, id)
				if err != nil {
					h.logger.Err(err).Str("room", room).Msgf("cannot leave room")
				}
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

			if onLeave != nil {
				err = onLeave(connId, id)
				if err != nil {
					h.logger.Err(err).Str("room", room).Msgf("cannot leave room")
				}
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
				if closeErr.Code != websocket.CloseNormalClosure && closeErr.Code != websocket.CloseGoingAway {
					h.logger.
						Warn().
						Err(err).
						Str("addr", conn.RemoteAddr().String()).
						Msg("connection closed unexpectedly")
				}
				h.removeConnections(connId)
			} else if errors.Is(err, websocket.ErrCloseSent) {
				h.logger.Warn().Err(err).Str("conn", connId).Str("addr", conn.RemoteAddr().String()).Msg("connection closed")
				h.removeConnections(connId)
			} else {
				h.logger.Err(err).Str("conn", connId).Str("addr", conn.RemoteAddr().String()).Msg("cannot read message from connection, connection will be closed")
				h.disconnectConnections(connId)
			}

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

			userID, err := h.authorizer.Authenticate(ctx, msg.Token)
			if err != nil {
				h.logger.Err(err).Msg("authentication failed")
			}
			if userID == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusUnauthorized,
				})
			} else {
				h.setConnAuth(connId, userID, msg.Token)
				closed = h.sendToConn(connId, WMessage{Type: WMessageAuthSuccess})
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

			closed = h.join(ctx, connId, msg.Room, msg.Data)
		case RMessageLeave:
			if msg.Room == "" {
				closed = h.sendToConn(connId, WMessage{
					Type:  WMessageFail,
					Error: http.StatusBadRequest,
					Msg:   "room is missing",
				})
				continue
			}
			h.leave(connId, msg.Room)
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

func (h *hub) setConnAuth(connId, userID, token string) {
	h.connsMx.Lock()
	defer h.connsMx.Unlock()

	h.conns[connId] = userConn{
		userID: userID,
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
			Msg("sent to websocket conn")
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

func (h *hub) sendToRoomByConns(room string, connIds []string, msg WMessage) []string {
	h.connsMx.RLock()
	h.roomsMx.RLock()
	defer func() {
		h.roomsMx.RUnlock()
		h.connsMx.RUnlock()
	}()

	toSend := make(map[string]struct{}, len(connIds))
	for _, id := range connIds {
		toSend[id] = struct{}{}
	}

	closedConns := make([]string, 0)
	count := 0
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
		if _, ok := toSend[connId]; !ok {
			continue
		}

		count++
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
		onLeave, id := h.getOnLeave(room)
		for _, connId := range conns {
			found := false
			for _, toRemove := range connIds {
				if connId == toRemove {
					found = true
					break
				}
			}

			if found {
				if onLeave != nil {
					err := onLeave(connId, id)
					if err != nil {
						h.logger.Err(err).Str("room", room).Msgf("cannot leave room")
					}
				}

				continue
			}

			filteredConns = append(filteredConns, connId)
		}

		h.rooms[room] = filteredConns
	}
}

func (h *hub) getOnNotExist(room string) (GroupOnNotExist, string) {
	h.groupsMx.RLock()
	defer h.groupsMx.RUnlock()
	for group, params := range h.groups {
		if strings.HasPrefix(room, group) {
			return params.OnNotExist, room[len(group):]
		}
	}

	return nil, ""
}

func (h *hub) getOnJoin(room string) (GroupOnJoin, string) {
	h.groupsMx.RLock()
	defer h.groupsMx.RUnlock()
	for group, params := range h.groups {
		if strings.HasPrefix(room, group) {
			return params.OnJoin, room[len(group):]
		}
	}

	return nil, ""
}

func (h *hub) getOnLeave(room string) (GroupOnLeave, string) {
	h.groupsMx.RLock()
	defer h.groupsMx.RUnlock()
	for group, params := range h.groups {
		if strings.HasPrefix(room, group) {
			return params.OnLeave, room[len(group):]
		}
	}

	return nil, ""
}

func (h *hub) getOnLeaveByGroup(group string) GroupOnLeave {
	h.groupsMx.RLock()
	defer h.groupsMx.RUnlock()
	return h.groups[group].OnLeave
}
