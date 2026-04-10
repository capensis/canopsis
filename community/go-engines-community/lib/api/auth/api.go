package auth

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/wsconn"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	headerAuthorization = "Authorization"
	bearerPrefix        = "Bearer"
)

type API interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetLoggedUserCount(c *gin.Context)
	GetFileAccess(c *gin.Context)
}

func NewApi(
	tokenService apisecurity.TokenService,
	tokenProviders []security.TokenProvider,
	providers []security.Provider,
	websocketStore wsconn.Store,
	maintenanceAdapter config.MaintenanceAdapter,
	enforcer security.Enforcer,
	cookieName string,
	cookieMaxAge int,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		tokenService:       tokenService,
		tokenProviders:     tokenProviders,
		providers:          providers,
		websocketStore:     websocketStore,
		maintenanceAdapter: maintenanceAdapter,
		enforcer:           enforcer,
		errorResponder:     errorResponder,
		logger:             logger,

		cookieName:     cookieName,
		cookieMaxAge:   cookieMaxAge,
		cookieSameSite: http.SameSiteNoneMode,
		cookieSecure:   true, // must be always set with SameSite=None
	}
}

type api struct {
	tokenService       apisecurity.TokenService
	tokenProviders     []security.TokenProvider
	providers          []security.Provider
	websocketStore     wsconn.Store
	maintenanceAdapter config.MaintenanceAdapter
	enforcer           security.Enforcer
	errorResponder     httperror.Responder
	logger             zerolog.Logger

	cookieName     string
	cookieMaxAge   int
	cookieSecure   bool
	cookieSameSite http.SameSite
}

// Login
// @Param body body LoginRequest true "body"
// @Success 200 {object} LoginResponse
func (a *api) Login(c *gin.Context) {
	var request LoginRequest

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	var user *security.User
	var err error
	var provider string
	for _, p := range a.providers {
		user, err = p.Auth(c, request.Username, request.Password)
		if err != nil {
			a.logger.Err(err).Msg("Auth provider error")
		}

		if user != nil {
			provider = p.GetName()
			break
		}
	}

	if user == nil {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	maintenanceConf, err := a.maintenanceAdapter.GetConfig(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if maintenanceConf.Enabled {
		ok, err := a.enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.ErrMaintenance)

			return
		}
	}

	accessToken, err := a.tokenService.Create(c, *user, provider)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response := LoginResponse{AccessToken: accessToken}

	c.JSON(http.StatusOK, response)
}

func (a *api) Logout(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}
	ok, err := a.tokenService.Delete(c, tokenString)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	c.SetSameSite(a.cookieSameSite)
	c.SetCookie(a.cookieName, tokenString, -1, "", "", a.cookieSecure, false)
	c.Status(http.StatusNoContent)
}

// GetLoggedUserCount
// @Success 200 {object} LoggedUserCountResponse
func (a *api) GetLoggedUserCount(c *gin.Context) {
	count, err := a.websocketStore.CountActiveConnections(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	c.JSON(http.StatusOK, LoggedUserCountResponse{
		Count: count,
	})
}

func (a *api) GetFileAccess(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	var user *security.User
	var err error
	for _, provider := range a.tokenProviders {
		user, err = provider.Auth(c, tokenString)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}
		if user != nil {
			break
		}
	}

	if user == nil {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	c.SetSameSite(a.cookieSameSite)
	c.SetCookie(a.cookieName, tokenString, a.cookieMaxAge, "", "", a.cookieSecure, false)
	c.Status(http.StatusNoContent)
}

func getToken(c *gin.Context) string {
	header := c.GetHeader(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(header[len(bearerPrefix):])
}
