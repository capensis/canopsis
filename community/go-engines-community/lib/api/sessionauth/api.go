// Package sessionauth contains authentication by session.
// Deprecated : don't use session.
package sessionauth

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

type API interface {
	LogoutHandler() gin.HandlerFunc
	LoginHandler() gin.HandlerFunc
}

func NewApi(
	sessionStore libsession.Store,
	providers []security.Provider,
	maintenanceAdapter config.MaintenanceAdapter,
	enforcer security.Enforcer,
	logger zerolog.Logger,
) API {
	return &api{
		sessionStore:       sessionStore,
		providers:          providers,
		maintenanceAdapter: maintenanceAdapter,
		enforcer:           enforcer,
		logger:             logger,
	}
}

type api struct {
	sessionStore       libsession.Store
	providers          []security.Provider
	maintenanceAdapter config.MaintenanceAdapter
	enforcer           security.Enforcer
	logger             zerolog.Logger
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
			user, err = p.Auth(c, request.Username, request.Password)
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

		maintenanceConf, err := a.maintenanceAdapter.GetConfig(c)
		if err != nil {
			panic(err)
		}

		if maintenanceConf.Enabled {
			ok, err := a.enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
				return
			}
		}

		var response loginResponse
		response.AuthApiKey = user.AuthApiKey
		response.Contact.Name = user.Contact.Name
		response.Contact.Address = user.Contact.Address
		response.Name = user.Name
		response.Email = user.Email
		response.Roles = user.Roles

		session.Values["user"] = user.ID
		err = session.Save(c.Request, c.Writer)

		if err != nil {
			panic(err)
		}

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

		c.Next()
	}
}

func (a *api) getSession(c *gin.Context) *sessions.Session {
	session, err := a.sessionStore.Get(c.Request, security.SessionKey)

	if err == nil {
		return session
	}

	var securecookieError securecookie.Error
	if errors.As(err, &securecookieError) {
		// if securecookie decode failed (for example due changed key), then it's a new session
		session, err = a.sessionStore.New(c.Request, security.SessionKey)
	}
	if err != nil {
		panic(err)
	}
	return session
}
