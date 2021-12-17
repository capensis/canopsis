package websocket

import (
	"github.com/gin-gonic/gin"
)

type API interface {
	Handler(c *gin.Context)
}

func NewApi(hub Hub) API {
	return &api{
		hub: hub,
	}
}

type api struct {
	hub Hub
}

func (a *api) Handler(c *gin.Context) {
	err := a.hub.Connect(c.Writer, c.Request)
	if err != nil {
		panic(err)
	}
}
