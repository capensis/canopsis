package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
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
	tokenService token.Service,
	tokenStore token.Store,
	providers []security.Provider,
	sessionStore session.Store,
	websocketHub websocket.Hub,
	cookieName string,
	cookieMaxAge int,
	cookieSecure bool,
	logger zerolog.Logger,
) API {
	return &api{
		tokenService: tokenService,
		tokenStore:   tokenStore,
		providers:    providers,
		websocketHub: websocketHub,
		sessionStore: sessionStore,
		logger:       logger,

		cookieName:     cookieName,
		cookieMaxAge:   cookieMaxAge,
		cookieSecure:   cookieSecure,
		cookieSameSite: http.SameSiteNoneMode,
	}
}

type api struct {
	tokenService token.Service
	tokenStore   token.Store
	providers    []security.Provider
	websocketHub websocket.Hub
	logger       zerolog.Logger

	cookieName     string
	cookieMaxAge   int
	cookieSecure   bool
	cookieSameSite http.SameSite

	sessionStore session.Store
}

// Log in
// @Summary Log in
// @Description Log in
// @Tags auth
// @ID auth-login
// @Accept json
// @Produce json
// @Param body body loginRequest true "body"
// @Success 200 {object} loginResponse
// @Router /login [post]
func (a *api) Login(c *gin.Context) {
	var request loginRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	var user *security.User
	var err error
	var provider string
	for _, p := range a.providers {
		user, err = p.Auth(c.Request.Context(), request.Username, request.Password)
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

	accessToken, expiresAt, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	err = a.tokenStore.Save(c.Request.Context(), token.Token{
		ID:       accessToken,
		User:     user.ID,
		Provider: provider,
		Created:  types.CpsTime{Time: now},
		Expired:  types.CpsTime{Time: expiresAt},
	})
	if err != nil {
		panic(err)
	}

	a.sendWebsocketMessage(c.Request.Context())

	response := loginResponse{AccessToken: accessToken}

	c.JSON(http.StatusOK, response)
}

// Log out
// @Summary Log out
// @Description Log out
// @Tags auth
// @ID auth-logout
// @Accept json
// @Produce json
// @Security JWTAuth
// @Success 204
// @Router /logout [post]
func (a *api) Logout(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}
	ok, err := a.tokenStore.Delete(c.Request.Context(), tokenString)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	a.sendWebsocketMessage(c.Request.Context())

	c.SetSameSite(a.cookieSameSite)
	c.SetCookie(a.cookieName, tokenString, -1, "", "", a.cookieSecure, false)
	c.Status(http.StatusNoContent)
}

// Get logged user count
// @Summary Get logged user count
// @Description Get logged user count
// @Tags auth
// @ID auth-logged-user-count
// @Security ApiKeyAuth
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Router /logged-user-count [get]
func (a *api) GetLoggedUserCount(c *gin.Context) {
	count, err := a.fetchLoggedUserCount(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, loggedUserCountResponse{
		Count: count,
	})
}

// Get file access
// @Summary Get file access
// @Description Get file access
// @Tags auth
// @ID auth-get-file-access
// @Accept json
// @Produce json
// @Security JWTAuth
// @Success 204
// @Router /file-access [get]
func (a *api) GetFileAccess(c *gin.Context) {
	tokenString := getToken(c)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	ok, err := a.tokenStore.Exists(c.Request.Context(), tokenString)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	c.SetSameSite(a.cookieSameSite)
	c.SetCookie(a.cookieName, tokenString, a.cookieMaxAge, "", "", a.cookieSecure, false)
	c.Status(http.StatusNoContent)
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

	// todo : remove after session delete
	sessionCount, err := a.sessionStore.GetActiveSessionsCount(ctx)
	if err != nil {
		return 0, err
	}

	return count + sessionCount, nil
}

func getToken(c *gin.Context) string {
	header := c.GetHeader(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(header[len(bearerPrefix):])
}
