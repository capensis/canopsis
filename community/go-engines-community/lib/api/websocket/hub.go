// Package websocket contains implementation of websocket.
package websocket

import (
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/keymutex"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// Hub interface is used to implement websocket room.
type Hub interface {
	// Subscribe creates listener connection and adds listener to room.
	Subscribe(w http.ResponseWriter, r *http.Request, room string) error
	// Send sends message to all listeners in room.
	Send(room string, msg interface{})
	// Stop closes all connections.
	Stop()
}

func NewHub(upgrader *websocket.Upgrader, logger zerolog.Logger) Hub {
	return &hub{
		upgrader: upgrader,
		roomsMx:  keymutex.New(),
		rooms:    make(map[string][]*websocket.Conn),
		logger:   logger,
	}
}

type hub struct {
	upgrader *websocket.Upgrader
	roomsMx  keymutex.KeyMutex
	rooms    map[string][]*websocket.Conn
	logger   zerolog.Logger
}

func (h *hub) Subscribe(w http.ResponseWriter, r *http.Request, room string) error {
	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	if len(h.rooms[room]) == 0 {
		h.rooms[room] = make([]*websocket.Conn, 0)
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	h.rooms[room] = append(h.rooms[room], conn)

	// Run goroutine to receive disconnect.
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err == nil {
				h.logger.Warn().
					Str("msg", string(msg)).
					Str("room", room).
					Str("addr", conn.RemoteAddr().String()).
					Msg("websocket shouldn't receive messages")
				continue
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				h.logger.
					Err(err).
					Str("room", room).
					Str("addr", conn.RemoteAddr().String()).
					Msg("connection closed unexpectedly")
			}

			h.removeConnection(room, conn)

			break
		}
	}()

	return nil
}

func (h *hub) Send(room string, msg interface{}) {
	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	for _, conn := range h.rooms[room] {
		err := conn.WriteJSON(msg)
		if err != nil {
			h.logger.Err(err).
				Str("room", room).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot write message to connection")
		}
	}
}

func (h *hub) Stop() {
	for room := range h.rooms {
		h.closeConnections(room)
	}
}

func (h *hub) closeConnections(room string) {
	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	for _, conn := range h.rooms[room] {
		err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
		if err != nil {
			h.logger.Err(err).
				Str("room", room).
				Str("addr", conn.RemoteAddr().String()).
				Msg("cannot close connection")
		}
	}

	h.rooms[room] = nil
}

func (h *hub) removeConnection(room string, conn *websocket.Conn) {
	h.roomsMx.Lock(room)
	defer h.roomsMx.Unlock(room)

	index := -1
	for i, v := range h.rooms[room] {
		if v == conn {
			index = i
			break
		}
	}

	if index >= 0 {
		h.rooms[room] = append(h.rooms[room][:index], h.rooms[room][index+1:]...)
	}
}
