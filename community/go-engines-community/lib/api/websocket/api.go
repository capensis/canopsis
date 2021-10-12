package websocket

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
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
	userId := c.MustGet(auth.UserKey).(string)
	err := a.hub.Connect(userId, c.Writer, c.Request)
	if err != nil {
		panic(err)
	}
}
