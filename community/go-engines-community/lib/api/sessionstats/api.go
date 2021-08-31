package sessionstats

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/stats"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

type API interface {
	StartHandler() gin.HandlerFunc
	PingHandler() gin.HandlerFunc
	ChangePathHandler() gin.HandlerFunc
	ListHandler() gin.HandlerFunc
}

func NewApi(sessionStore sessions.Store, manager stats.Manager) API {
	return &api{
		sessionStore: sessionStore,
		manager:      manager,
	}
}

type api struct {
	sessionStore sessions.Store
	manager      stats.Manager
}

func (a *api) StartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, userID := a.getSession(c)
		s, err := a.manager.Ping(
			c.Request.Context(),
			stats.SessionData{
				SessionID: session.ID,
				UserID:    userID,
			},
			stats.PathData{},
		)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, s)
	}
}

func (a *api) PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r pingRequest

		if err := c.ShouldBind(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
			return
		}

		session, userID := a.getSession(c)
		viewID := r.Path[0]
		tabID := ""
		if len(r.Path) > 1 {
			tabID = r.Path[1]
		}

		s, err := a.manager.Ping(
			c.Request.Context(),
			stats.SessionData{
				SessionID: session.ID,
				UserID:    userID,
			},
			stats.PathData{
				ViewID:  viewID,
				TabID:   tabID,
				Visible: r.Visible != nil && *r.Visible,
			},
		)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, s)
	}
}

func (a *api) ChangePathHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r changePathRequest

		if err := c.ShouldBind(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
			return
		}

		session, userID := a.getSession(c)
		viewID := r.Path[0]
		tabID := ""
		if len(r.Path) > 1 {
			tabID = r.Path[1]
		}

		s, err := a.manager.Ping(
			c.Request.Context(),
			stats.SessionData{
				SessionID: session.ID,
				UserID:    userID,
			},
			stats.PathData{
				ViewID:  viewID,
				TabID:   tabID,
				Visible: true,
			},
		)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, s)
	}
}

func (a *api) ListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r listRequest

		if err := c.ShouldBind(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
			return
		}

		s, err := a.manager.Find(c.Request.Context(), stats.Filter{
			IsActive:      r.IsActive,
			Usernames:     r.Usernames,
			StartedAfter:  r.StartedAfter,
			StoppedBefore: r.StoppedBefore,
		})
		if err != nil {
			panic(err)
		}

		res := listResponse{
			Stats: s,
		}
		c.JSON(http.StatusOK, res)
	}
}

func (a *api) getSession(c *gin.Context) (*sessions.Session, string) {
	session, err := a.sessionStore.Get(c.Request, security.SessionKey)
	if err != nil {
		panic(err)
	}

	userID, ok := session.Values["user"].(string)
	if !ok {
		panic(errors.New("session is empty"))
	}

	return session, userID
}
