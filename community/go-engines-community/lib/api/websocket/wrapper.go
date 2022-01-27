package websocket

//go:generate mockgen -destination=../../../mocks/lib/api/websocket/websocket.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket Upgrader,Connection,Authorizer,Hub

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"time"
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

	return &connection{Conn: conn}, nil
}

type connection struct {
	*websocket.Conn
}
