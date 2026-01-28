package websocket

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type API interface {
	Handler(c *gin.Context)
}

func NewApi(hub Hub, errorResponder response.ErrorResponder) API {
	return &api{
		hub:            hub,
		errorResponder: errorResponder,
	}
}

type api struct {
	hub            Hub
	errorResponder response.ErrorResponder
}

func (a *api) Handler(c *gin.Context) {
	err := a.hub.Connect(c.Writer, c.Request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
}
