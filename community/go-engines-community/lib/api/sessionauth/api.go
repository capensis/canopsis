// Package sessionauth contains authentication by session.
// Deprecated : don't use session.
package sessionauth

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

type API interface {
	LogoutHandler() gin.HandlerFunc
	LoginHandler() gin.HandlerFunc
}

type TokenStore interface {
	Count(ctx context.Context) (int64, error)
}

func NewApi(
	sessionStore libsession.Store,
	providers []security.Provider,
	websocketHub websocket.Hub,
	tokenStore TokenStore,
	logger zerolog.Logger,
) API {
	return &api{
		sessionStore: sessionStore,
		providers:    providers,
		websocketHub: websocketHub,
		tokenStore:   tokenStore,
		logger:       logger,
	}
}

type api struct {
	sessionStore libsession.Store
	providers    []security.Provider
	websocketHub websocket.Hub
	tokenStore   TokenStore
	logger       zerolog.Logger
}

// LoginHandler authenticates user and starts sessions.
func (a *api) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := a.getSession(c)
		var request loginRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		var user *security.User
		var err error
		for _, p := range a.providers {
			user, err = p.Auth(c.Request.Context(), request.Username, request.Password)
			if err != nil {
				a.logger.Err(err).Msg("Auth provider error")
			}

			if user != nil {
				break
			}
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		var response loginResponse
		response.AuthApiKey = user.AuthApiKey
		response.Contact.Name = user.Contact.Name
		response.Contact.Address = user.Contact.Address
		response.Name = user.Name
		response.Email = user.Email
		response.Role = user.Role

		session.Values["user"] = user.ID
		err = session.Save(c.Request, c.Writer)

		if err != nil {
			panic(err)
		}

		a.sendWebsocketMessage(c)

		c.JSON(http.StatusOK, response)
	}
}

// LogoutHandler deletes session.
func (a *api) LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := a.getSession(c)
		session.Options.MaxAge = -1
		err := session.Save(c.Request, c.Writer)

		if err != nil {
			panic(err)
		}

		a.sendWebsocketMessage(c)
		c.Next()
	}
}

func (a *api) getSession(c *gin.Context) *sessions.Session {
	session, err := a.sessionStore.Get(c.Request, security.SessionKey)

	if err != nil {
		panic(err)
	}

	return session
}

func (a *api) sendWebsocketMessage(ctx context.Context) {
	count, err := a.fetchLoggedUserCount(ctx)
	if err != nil {
		panic(err)
	}
	a.websocketHub.Send(websocket.RoomLoggedUserCount, count)
}

func (a *api) fetchLoggedUserCount(ctx context.Context) (int64, error) {
	count, err := a.tokenStore.Count(ctx)
	if err != nil {
		return 0, err
	}

	sessionCount, err := a.sessionStore.GetActiveSessionsCount(ctx)
	if err != nil {
		return 0, err
	}

	return count + sessionCount, nil
}
