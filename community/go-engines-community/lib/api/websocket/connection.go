package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// pongDeadlineRatio multiplied by the ping interval gives the read deadline.
// Using 2x allows one missed pong before the connection is considered dead.
const pongDeadlineRatio = 2

func newConnection(
	conn Conn,
	roomRegistry RoomRegistry,
	authenticate Authenticate,
	configProvider config.ApiConfigProvider,
	logger zerolog.Logger,
) *connection {
	return &connection{
		id:             utils.NewID(),
		conn:           conn,
		roomRegistry:   roomRegistry,
		authenticate:   authenticate,
		configProvider: configProvider,
		logger:         logger,
		writeCh:        make(chan ServerMessage, msgBuffSize),
		done:           make(chan struct{}),
		rooms:          make(map[string]bool),
	}
}

type connection struct {
	id             string
	conn           Conn
	roomRegistry   RoomRegistry
	authenticate   Authenticate
	configProvider config.ApiConfigProvider
	logger         zerolog.Logger
	closeOnce      sync.Once
	writeCh        chan ServerMessage
	done           chan struct{}
	roomsMx        sync.RWMutex
	rooms          map[string]bool
	userIDMx       sync.RWMutex
	userID         string
	userToken      string
}

// close uses sync.Once so it is safe to call from both runRead and runWrite
// goroutines — whichever exits first triggers cleanup, the other is a no-op.
func (c *connection) close(ctx context.Context) error {
	var err error
	c.closeOnce.Do(func() {
		err = c.conn.Close()
		close(c.done)

		c.roomsMx.Lock()
		rooms := c.rooms
		c.rooms = make(map[string]bool)
		c.roomsMx.Unlock()

		userID := c.getUserID()
		c.logger.Debug().Str("conn", c.id).Str("user", userID).Msg("connection closed")
		opts := LeaveOptions{
			ConnID: c.id,
		}
		for room := range rooms {
			h, ok := c.roomRegistry.Get(room)
			if ok && h.OnLeave != nil {
				if lerr := h.OnLeave(ctx, opts); lerr != nil {
					c.logger.Err(lerr).Str("room", room).Str("user", userID).Msg("cannot leave from room")
				}
			}
		}
	})

	return err
}

func (c *connection) runRead(ctx context.Context) {
	pongInterval := pongDeadlineRatio * c.configProvider.Get().WebsocketPingInterval
	err := c.conn.SetReadDeadline(time.Now().Add(pongInterval))
	if err != nil {
		c.logger.Err(err).Str("addr", c.conn.RemoteAddr().String()).Msg("cannot set read deadline")
	}

	c.conn.SetPongHandler(func(string) error {
		pongInterval := pongDeadlineRatio * c.configProvider.Get().WebsocketPingInterval
		err := c.conn.SetReadDeadline(time.Now().Add(pongInterval))
		if err != nil {
			c.logger.Err(err).Str("addr", c.conn.RemoteAddr().String()).Msg("cannot set read deadline")
		}

		return nil
	})

	for {
		msg := ClientMessage{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			unmarshalErr := &json.UnmarshalTypeError{}
			syntaxErr := &json.SyntaxError{}
			if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
				c.writeError(ctx, "", http.StatusBadRequest, "cannot parse JSON")
				continue
			}

			if _, ok := errors.AsType[*websocket.CloseError](err); ok || errors.Is(err, websocket.ErrCloseSent) {
				c.logger.Debug().Err(err).
					Str("addr", c.conn.RemoteAddr().String()).
					Msg("connection closed")

				return
			}

			c.logger.Err(err).
				Str("addr", c.conn.RemoteAddr().String()).
				Msg("cannot read message from connection, connection will be closed")

			return
		}

		c.logger.Debug().Int("type", msg.Type).Str("room", msg.Room).Interface("payload", msg.Payload).Msg("received message")
		switch msg.Type {
		case ClientMessageClientPing:
			c.write(ctx, ServerMessage{Type: ServerMessageClientPong})
		case ClientMessageJoin:
			c.join(ctx, msg)
		case ClientMessageLeave:
			c.leave(ctx, msg)
		case ClientMessageAuth:
			c.auth(ctx, msg)
		case ClientMessageInfo:
			c.msg(ctx, msg)
		default:
			c.writeError(ctx, msg.Room, http.StatusBadRequest, map[string]string{
				"type": "invalid value",
			})
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (c *connection) runWrite(ctx context.Context) {
	pingInterval := c.configProvider.Get().WebsocketPingInterval
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	var err error
	defer func() {
		// CloseMessage is sent only on a clean exit. If err != nil the underlying
		// connection is already broken, so sending CloseMessage is skipped.
		if err == nil {
			_ = c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case msg, ok := <-c.writeCh:
			if !ok {
				return
			}

			err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.logger.Err(err).Str("addr", c.conn.RemoteAddr().String()).Msg("cannot set write deadline")
			}

			c.logger.Debug().Int("type", msg.Type).Str("room", msg.Room).Interface("payload", msg.Payload).Msg("writing message")
			err = c.conn.WriteJSON(msg)
			if err != nil {
				c.logger.Err(err).
					Str("addr", c.conn.RemoteAddr().String()).
					Msg("cannot write message to connection, connection will be closed")

				return
			}
		case <-t.C:
			err = c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
			if err != nil {
				c.logger.Err(err).
					Str("addr", c.conn.RemoteAddr().String()).
					Msg("cannot ping connection, connection will be closed")

				return
			}

			newInterval := c.configProvider.Get().WebsocketPingInterval
			if pingInterval != newInterval {
				pingInterval = newInterval
				t.Reset(pingInterval)
			}
		}
	}
}

func (c *connection) write(ctx context.Context, msg ServerMessage) {
	select {
	case <-ctx.Done():
	case <-c.done:
	case c.writeCh <- msg:
	default:
		c.logger.Warn().Str("conn", c.id).Str("user", c.getUserID()).Msg("write buffer full, closing connection")
		err := c.close(ctx)
		if err != nil {
			c.logger.Err(err).Msg("cannot close connection")
		}
	}
}

func (c *connection) writeError(ctx context.Context, room string, status int, payload any) {
	if payload == nil {
		payload = http.StatusText(status)
	}

	c.write(ctx, ServerMessage{
		Type:    ServerMessageError,
		Room:    room,
		Error:   status,
		Payload: payload,
	})
}

func (c *connection) join(ctx context.Context, msg ClientMessage) {
	if msg.Room == "" {
		c.writeError(ctx, "", http.StatusBadRequest, map[string]string{
			"room": "is missing",
		})

		return
	}

	h, ok := c.roomRegistry.Get(msg.Room)
	if !ok {
		c.writeError(ctx, msg.Room, http.StatusNotFound, nil)

		return
	}

	userID := c.getUserID()
	if h.Authorize != nil {
		if userID == "" {
			c.writeError(ctx, msg.Room, http.StatusUnauthorized, nil)

			return
		}

		granted, err := h.Authorize(ctx, userID)
		if err != nil {
			c.logger.Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot authorize user")
			c.writeError(ctx, msg.Room, http.StatusInternalServerError, nil)

			return
		}

		if !granted {
			c.writeError(ctx, msg.Room, http.StatusForbidden, nil)

			return
		}
	}

	c.roomsMx.Lock()
	exists := c.rooms[msg.Room]
	c.rooms[msg.Room] = true
	c.roomsMx.Unlock()
	if exists {
		return
	}

	if h.OnJoin != nil {
		err := h.OnJoin(ctx, JoinOptions{
			ConnID:  c.id,
			UserID:  userID,
			Payload: msg.Payload,
		})
		if err != nil {
			c.roomsMx.Lock()
			delete(c.rooms, msg.Room)
			c.roomsMx.Unlock()

			if jerr, ok := errors.AsType[*JoinError](err); ok {
				c.logger.Debug().Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot join to room")
				c.write(ctx, ServerMessage{
					Type:    ServerMessageInfo,
					Room:    msg.Room,
					Payload: jerr.Payload,
				})
				c.write(ctx, ServerMessage{
					Type: ServerMessageCloseRoom,
					Room: msg.Room,
				})

				return
			}

			if verr, ok := errors.AsType[*ValidationError](err); ok {
				c.logger.Debug().Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot join to room")
				c.writeError(ctx, msg.Room, http.StatusBadRequest, verr.Payload)

				return
			}

			if errors.Is(err, ErrRoomNotFound) {
				c.writeError(ctx, msg.Room, http.StatusNotFound, nil)

				return
			}

			c.logger.Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot join to room")
			c.writeError(ctx, msg.Room, http.StatusInternalServerError, nil)

			return
		}
	}

	c.logger.Debug().Str("conn", c.id).Str("room", msg.Room).Str("user", userID).Msg("joined room")
}

func (c *connection) leave(ctx context.Context, msg ClientMessage) {
	room := msg.Room
	if room == "" {
		c.writeError(ctx, "", http.StatusBadRequest, map[string]string{
			"room": "is missing",
		})

		return
	}

	c.roomsMx.Lock()
	exists := c.rooms[room]
	if exists {
		delete(c.rooms, room)
	}

	c.roomsMx.Unlock()
	if !exists {
		return
	}

	userID := c.getUserID()
	h, ok := c.roomRegistry.Get(room)
	if ok && h.OnLeave != nil {
		err := h.OnLeave(ctx, LeaveOptions{
			ConnID: c.id,
		})
		if err != nil {
			c.logger.Err(err).Str("room", room).Str("user", userID).Msg("cannot leave from room")
			c.writeError(ctx, room, http.StatusInternalServerError, nil)

			return
		}
	}

	c.logger.Debug().Str("conn", c.id).Str("room", room).Str("user", userID).Msg("left room")
}

func (c *connection) auth(ctx context.Context, msg ClientMessage) {
	if msg.Token == "" {
		c.writeError(ctx, "", http.StatusBadRequest, map[string]string{
			"token": "is missing",
		})

		return
	}

	userID, err := c.authenticate(ctx, msg.Token)
	if err != nil {
		c.logger.Debug().Err(err).Msg("cannot authenticate token")
		c.writeError(ctx, "", http.StatusUnauthorized, nil)

		return
	}

	changed := c.setAuth(userID, msg.Token)

	c.logger.Debug().Str("conn", c.id).Str("user", userID).Msg("user authenticated")
	c.write(ctx, ServerMessage{
		Type: ServerMessageAuthSuccess,
	})

	if changed {
		c.checkRooms(ctx)
	}
}

func (c *connection) msg(ctx context.Context, msg ClientMessage) {
	if msg.Room == "" {
		c.writeError(ctx, "", http.StatusBadRequest, map[string]string{
			"room": "is missing",
		})

		return
	}

	if !c.inRoom(msg.Room) {
		c.writeError(ctx, msg.Room, http.StatusForbidden, nil)

		return
	}

	h, ok := c.roomRegistry.Get(msg.Room)
	if !ok {
		c.writeError(ctx, msg.Room, http.StatusNotFound, nil)

		return
	}

	if h.OnMessage == nil {
		c.writeError(ctx, msg.Room, http.StatusMethodNotAllowed, nil)

		return
	}

	userID := c.getUserID()
	err := h.OnMessage(ctx, MessageOptions{
		ConnID:  c.id,
		UserID:  userID,
		Payload: msg.Payload,
	})
	if err != nil {
		if verr, ok := errors.AsType[*ValidationError](err); ok {
			c.logger.Debug().Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot handle message from room")
			c.writeError(ctx, msg.Room, http.StatusBadRequest, verr.Payload)

			return
		}

		if cerr, ok := errors.AsType[*CloseRoomError](err); ok {
			c.logger.Debug().Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot handle message from room")
			c.writeError(ctx, msg.Room, http.StatusGone, cerr.Payload)
			c.leaveRoom(ctx, msg.Room)

			return
		}

		c.logger.Err(err).Str("room", msg.Room).Str("user", userID).Msg("cannot handle message from room")
		c.writeError(ctx, msg.Room, http.StatusInternalServerError, nil)

		return
	}
}

func (c *connection) checkAuth(ctx context.Context) {
	c.userIDMx.RLock()
	token := c.userToken
	c.userIDMx.RUnlock()

	if token == "" {
		return
	}

	_, err := c.authenticate(ctx, token)
	if err != nil {
		c.setAuth("", "")

		c.logger.Debug().Err(err).Msg("cannot authenticate token")
		c.writeError(ctx, "", http.StatusUnauthorized, nil)
	}

	// checkRooms is called even on auth failure: after clearing credentials above,
	// it evicts the connection from any rooms that require authorization.
	c.checkRooms(ctx)
}

func (c *connection) checkRooms(ctx context.Context) {
	userID := c.getUserID()
	removedRooms := make([]string, 0)

	c.roomsMx.RLock()
	rooms := make(map[string]bool, len(c.rooms))
	maps.Copy(rooms, c.rooms)
	c.roomsMx.RUnlock()

	for room := range rooms {
		h, ok := c.roomRegistry.Get(room)
		if ok && h.Authorize != nil {
			var granted bool
			var err error
			if userID == "" {
				granted = false
			} else if granted, err = h.Authorize(ctx, userID); err != nil {
				c.logger.Err(err).Str("room", room).Str("user", userID).Msg("cannot authorize user")
				continue
			}

			if !granted {
				removedRooms = append(removedRooms, room)
			}
		}
	}

	c.roomsMx.Lock()
	for _, room := range removedRooms {
		delete(c.rooms, room)
	}

	c.roomsMx.Unlock()

	for _, room := range removedRooms {
		h, ok := c.roomRegistry.Get(room)
		if ok && h.OnLeave != nil {
			err := h.OnLeave(ctx, LeaveOptions{
				ConnID: c.id,
			})
			if err != nil {
				c.logger.Err(err).Str("room", room).Str("user", userID).Msg("cannot leave from room")
			}
		}

		c.writeError(ctx, room, http.StatusForbidden, nil)
	}
}

func (c *connection) leaveRoom(ctx context.Context, room string) {
	c.roomsMx.Lock()
	exists := c.rooms[room]
	if exists {
		delete(c.rooms, room)
	}

	c.roomsMx.Unlock()
	if !exists {
		return
	}

	h, ok := c.roomRegistry.Get(room)
	if ok && h.OnLeave != nil {
		userID := c.getUserID()
		err := h.OnLeave(ctx, LeaveOptions{
			ConnID: c.id,
		})
		if err != nil {
			c.logger.Err(err).Str("room", room).Str("user", userID).Msg("cannot leave from room")
		}
	}

	c.write(ctx, ServerMessage{
		Type: ServerMessageCloseRoom,
		Room: room,
	})
}

func (c *connection) inRoom(room string) bool {
	c.roomsMx.RLock()
	defer c.roomsMx.RUnlock()

	return c.rooms[room]
}

func (c *connection) getUserID() string {
	c.userIDMx.RLock()
	defer c.userIDMx.RUnlock()

	return c.userID
}

func (c *connection) getAuth() (userID, token string) {
	c.userIDMx.RLock()
	defer c.userIDMx.RUnlock()

	return c.userID, c.userToken
}

func (c *connection) setAuth(userID, token string) bool {
	c.userIDMx.Lock()
	defer c.userIDMx.Unlock()

	changes := c.userID != userID || c.userToken != token

	c.userID = userID
	c.userToken = token

	return changes
}

func (c *connection) getRoomsForGroup(group string) []string {
	prefix := group + groupDelimiter

	c.roomsMx.RLock()
	defer c.roomsMx.RUnlock()

	ids := make([]string, 0)
	for room := range c.rooms {
		if id, ok := strings.CutPrefix(room, prefix); ok {
			ids = append(ids, id)
		}
	}

	return ids
}
