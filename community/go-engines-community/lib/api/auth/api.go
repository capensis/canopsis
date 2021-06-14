package auth

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/security"
	libsession "git.canopsis.net/canopsis/go-engines/lib/security/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

type API interface {
	// LogoutHandler deletes session.
	LogoutHandler() gin.HandlerFunc
	// LoginHandler authenticates user and starts sessions.
	LoginHandler() gin.HandlerFunc
	// GetSessionsCount returns active sessions count.
	GetSessionsCount() gin.HandlerFunc
}

func NewApi(
	sessionStore libsession.Store,
	providers []security.Provider,
) API {
	return &api{
		sessionStore: sessionStore,
		providers:    providers,
	}
}

type api struct {
	sessionStore libsession.Store
	providers    []security.Provider
}

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
			user, err = p.Auth(request.Username, request.Password)
			if err != nil {
				panic(err)
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

		c.JSON(http.StatusOK, response)
	}
}

func (a *api) LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := a.getSession(c)
		session.Options.MaxAge = -1
		err := session.Save(c.Request, c.Writer)

		if err != nil {
			panic(err)
		}

		c.Next()
	}
}

// Get counts of active sessions
// @Summary Get counts of active sessions
// @Description Get counts of active sessions
// @Tags auth
// @ID auth-get-session-counts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 200 {object} sessionsCountResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /sessions-count [get]
func (a *api) GetSessionsCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := a.sessionStore.GetActiveSessionsCount()
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, sessionsCountResponse{Count: count})
	}
}

func (a *api) getSession(c *gin.Context) *sessions.Session {
	session, err := a.sessionStore.Get(c.Request, security.SessionKey)

	if err != nil {
		panic(err)
	}

	return session
}
