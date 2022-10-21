package websocket

//go:generate mockgen -destination=../../../mocks/lib/api/websocket/websocket.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket Upgrader,Connection,Authorizer,Hub

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error)
}

func NewUpgrader(u websocket.Upgrader) Upgrader {
	return &upgrader{Upgrader: u}
}

type Connection interface {
	WriteControl(messageType int, data []byte, deadline time.Time) error
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	Close() error
	SetReadDeadline(t time.Time) error
	SetPongHandler(h func(string) error)
	RemoteAddr() net.Addr
}

type upgrader struct {
	websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error) {
	conn, err := u.Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return &connection{conn: conn}, nil
}

type connection struct {
	conn *websocket.Conn

	readMx  sync.Mutex
	writeMx sync.Mutex
}

func (c *connection) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return c.conn.WriteControl(messageType, data, deadline)
}

func (c *connection) WriteJSON(v interface{}) error {
	c.writeMx.Lock()
	defer c.writeMx.Unlock()

	return c.conn.WriteJSON(v)
}

func (c *connection) ReadJSON(v interface{}) error {
	c.readMx.Lock()
	defer c.readMx.Unlock()

	return c.conn.ReadJSON(v)
}

func (c *connection) Close() error {
	return c.conn.Close()
}

func (c *connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connection) SetPongHandler(h func(string) error) {
	c.readMx.Lock()
	defer c.readMx.Unlock()

	c.conn.SetPongHandler(h)
}

func (c *connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
