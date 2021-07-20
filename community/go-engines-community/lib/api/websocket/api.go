package websocket

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Handler(c *gin.Context)
}

func NewApi(hub Hub, rooms []string) API {
	return &api{
		hub:   hub,
		rooms: rooms,
	}
}

type api struct {
	rooms []string
	hub   Hub
}

func (a *api) Handler(c *gin.Context) {
	room := c.Param("room")
	found := false
	for _, v := range a.rooms {
		if v == room {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err := a.hub.Subscribe(c.Writer, c.Request, room)
	if err != nil {
		panic(err)
	}
}
