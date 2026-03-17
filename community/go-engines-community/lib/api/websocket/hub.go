package websocket

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/rs/zerolog"
)

// Lock ordering to prevent deadlocks: connsMu → roomsMu → userIDMu.
// Never acquire an outer lock while holding an inner one.

const (
	// connBuffSize caps the number of connections waiting to be registered.
	// Connect returns an error if the buffer is full rather than blocking.
	connBuffSize = 10
	// msgBuffSize is the outbound write buffer per connection.
	// If full, the connection is closed to prevent a slow client from blocking the hub.
	msgBuffSize = 256

	writeWait = 10 * time.Second
)

type Hub interface {
	Connect(w http.ResponseWriter, r *http.Request) error
	Run(ctx context.Context)
	SendMessage(ctx context.Context, payload any, f Filter)
	ConnectedUserIDs() []string
	ConnectedGroupRoomIDs(group string) []string
	ConnectionsInRoom(room string) []ConnectionInfo
	Connections() []ConnectionAuthInfo
	LeaveRoom(ctx context.Context, room string)
	SendMessageToUser(ctx context.Context, payload any, room, userID string) bool
}

type Authenticate func(ctx context.Context, token string) (userID string, err error)

func NewHub(
	upgrader Upgrader,
	roomRegistry RoomRegistry,
	authenticate Authenticate,
	configProvider config.ApiConfigProvider,
	checkAuthTokenInterval time.Duration,
	logger zerolog.Logger,
) Hub {
	return &hub{
		upgrader:               upgrader,
		roomRegistry:           roomRegistry,
		authenticate:           authenticate,
		configProvider:         configProvider,
		checkAuthTokenInterval: checkAuthTokenInterval,
		logger:                 logger,
		registerCh:             make(chan *connection, connBuffSize),
		conns:                  make(map[string]*connection),
	}
}

type hub struct {
	upgrader               Upgrader
	roomRegistry           RoomRegistry
	authenticate           Authenticate
	configProvider         config.ApiConfigProvider
	checkAuthTokenInterval time.Duration
	logger                 zerolog.Logger
	connsMu                sync.RWMutex
	conns                  map[string]*connection
	registerCh             chan *connection
}

func (h *hub) Connect(w http.ResponseWriter, r *http.Request) error {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	select {
	case h.registerCh <- newConnection(c, h.roomRegistry, h.authenticate, h.configProvider, h.logger):
		return nil
	default:
		_ = c.Close()

		return errors.New("too many connections")
	}
}

// Run processes new connections and periodically re-checks authentication.
// It blocks until ctx is cancelled, then closes all active connections with
// a short grace period so in-flight writes can complete.
func (h *hub) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := sync.WaitGroup{}
	// check auth for connections
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(h.checkAuthTokenInterval)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				h.connsMu.RLock()
				conns := make(map[string]*connection, len(h.conns))
				maps.Copy(conns, h.conns)
				h.connsMu.RUnlock()

				for _, c := range conns {
					c.checkAuth(ctx)
				}
			}
		}
	}()

	// register connections
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case c := <-h.registerCh:
			h.logger.Debug().Str("conn", c.id).Str("addr", c.conn.RemoteAddr().String()).Msg("connection registered")
			h.connsMu.Lock()
			h.conns[c.id] = c
			h.connsMu.Unlock()
			wg.Add(2)
			go func() {
				defer wg.Done()
				defer h.closeAndRmConn(ctx, c.id)
				c.runRead(ctx)
			}()
			go func() {
				defer wg.Done()
				defer h.closeConn(ctx, c.id)
				c.runWrite(ctx)
			}()
		}
	}

	// close all connections
	h.connsMu.Lock()
	conns := h.conns
	h.conns = make(map[string]*connection)
	h.connsMu.Unlock()
	closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	for _, c := range conns {
		go func() {
			err := c.close(closeCtx)
			if err != nil {
				h.logger.Err(err).Msg("cannot close connection")
			}
		}()
	}

	wg.Wait()
}

func (h *hub) SendMessage(ctx context.Context, payload any, f Filter) {
	msg := ServerMessage{
		Type:    ServerMessageInfo,
		Room:    f.room(),
		Payload: payload,
	}
	h.sendToConns(ctx, msg, f.match)
}

func (h *hub) SendMessageToUser(ctx context.Context, payload any, room, userID string) bool {
	msg := ServerMessage{
		Type:    ServerMessageInfo,
		Room:    room,
		Payload: payload,
	}
	f := ToUser(room, userID)

	return h.sendToConns(ctx, msg, f.match)
}

func (h *hub) ConnectedUserIDs() []string {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()

	seen := make(map[string]bool, len(h.conns))
	ids := make([]string, 0, len(h.conns))
	for _, c := range h.conns {
		if userID := c.getUserID(); userID != "" {
			if !seen[userID] {
				seen[userID] = true
				ids = append(ids, userID)
			}
		}
	}

	return ids
}

func (h *hub) ConnectionsInRoom(room string) []ConnectionInfo {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()

	conns := make([]ConnectionInfo, 0)
	for _, c := range h.conns {
		if c.inRoom(room) {
			conns = append(conns, ConnectionInfo{
				ID:     c.id,
				UserID: c.getUserID(),
			})
		}
	}

	return conns
}

func (h *hub) LeaveRoom(ctx context.Context, room string) {
	h.connsMu.RLock()
	conns := make([]*connection, 0, len(h.conns))
	for _, c := range h.conns {
		if c.inRoom(room) {
			conns = append(conns, c)
		}
	}

	h.connsMu.RUnlock()
	for _, c := range conns {
		c.leaveRoom(ctx, room)
	}
}

func (h *hub) Connections() []ConnectionAuthInfo {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()

	conns := make([]ConnectionAuthInfo, 0, len(h.conns))
	for _, c := range h.conns {
		userID, token := c.getAuth()
		conns = append(conns, ConnectionAuthInfo{
			ID:     c.id,
			UserID: userID,
			Token:  token,
		})
	}

	return conns
}

func (h *hub) ConnectedGroupRoomIDs(group string) []string {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()

	seen := make(map[string]bool)
	ids := make([]string, 0)
	for _, c := range h.conns {
		for _, roomID := range c.getRoomsForGroup(group) {
			if !seen[roomID] {
				seen[roomID] = true
				ids = append(ids, roomID)
			}
		}
	}

	return ids
}

func (h *hub) sendToConns(ctx context.Context, msg ServerMessage, filter func(c *connection) bool) bool {
	h.connsMu.RLock()
	conns := make([]*connection, 0)
	for _, c := range h.conns {
		if filter(c) {
			conns = append(conns, c)
		}
	}

	h.connsMu.RUnlock()
	for _, c := range conns {
		c.write(ctx, msg)
	}

	return len(conns) > 0
}

func (h *hub) closeAndRmConn(ctx context.Context, connID string) {
	h.connsMu.Lock()
	c, ok := h.conns[connID]
	if ok {
		delete(h.conns, connID)
	}

	h.connsMu.Unlock()

	if ok {
		err := c.close(ctx)
		if err != nil {
			h.logger.Err(err).Msg("cannot close connection")
		}
	}
}

func (h *hub) closeConn(ctx context.Context, connID string) {
	h.connsMu.RLock()
	c, ok := h.conns[connID]
	h.connsMu.RUnlock()

	if ok {
		err := c.close(ctx)
		if err != nil {
			h.logger.Err(err).Msg("cannot close connection")
		}
	}
}
