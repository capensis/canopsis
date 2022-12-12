package auth

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
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
	websocketHub websocket.Hub,
	cookieName string,
	cookieMaxAge int,
	logger zerolog.Logger,
) API {
	return &api{
		tokenService:   tokenService,
		tokenProviders: tokenProviders,
		providers:      providers,
		websocketHub:   websocketHub,
		logger:         logger,

		cookieName:     cookieName,
		cookieMaxAge:   cookieMaxAge,
		cookieSameSite: http.SameSiteNoneMode,
		cookieSecure:   true, // must be always set with SameSite=None
	}
}

type api struct {
	tokenService   apisecurity.TokenService
	tokenProviders []security.TokenProvider
	providers      []security.Provider
	websocketHub   websocket.Hub
	logger         zerolog.Logger

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

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	accessToken, err := a.tokenService.Create(c, *user, provider)
	if err != nil {
		panic(err)
	}

	response := LoginResponse{AccessToken: accessToken}

	c.JSON(http.StatusOK, response)
}

func (a *api) Logout(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}
	ok, err := a.tokenService.Delete(c, tokenString)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	c.SetSameSite(a.cookieSameSite)
	c.SetCookie(a.cookieName, tokenString, -1, "", "", a.cookieSecure, false)
	c.Status(http.StatusNoContent)
}

// GetLoggedUserCount
// @Success 200 {object} LoggedUserCountResponse
func (a *api) GetLoggedUserCount(c *gin.Context) {
	c.JSON(http.StatusOK, LoggedUserCountResponse{
		Count: int64(a.websocketHub.GetAuthConnectionsCount()),
	})
}

func (a *api) GetFileAccess(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	var user *security.User
	var err error
	for _, provider := range a.tokenProviders {
		user, err = provider.Auth(c, tokenString)
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
