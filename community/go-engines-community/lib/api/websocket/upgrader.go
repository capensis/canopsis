package websocket

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/api/websocket/websocket.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket Upgrader,Conn,Hub

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error)
}

type Conn interface {
	SetReadDeadline(t time.Time) error
	ReadJSON(v any) error

	SetPongHandler(h func(appData string) error)

	SetWriteDeadline(t time.Time) error
	WriteControl(messageType int, data []byte, deadline time.Time) error
	WriteJSON(v interface{}) error

	RemoteAddr() net.Addr

	Close() error
}

func NewUpgrader(gorillaUpgrader *websocket.Upgrader) Upgrader {
	return &upgrader{
		gorillaUpgrader: gorillaUpgrader,
	}
}

type upgrader struct {
	gorillaUpgrader *websocket.Upgrader
}

func (u *upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error) {
	return u.gorillaUpgrader.Upgrade(w, r, responseHeader)
}
