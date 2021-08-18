package auth

import (
	"net/http"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"github.com/gin-gonic/gin"
)

const (
	headerAuthorization = "Authorization"
	bearerPrefix        = "Bearer"
)

type API interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetLoggedUserCount(c *gin.Context)
}

func NewApi(
	tokenService token.Service,
	tokenStore token.Store,
	providers []security.Provider,
	sessionStore session.Store,
) API {
	return &api{
		tokenService: tokenService,
		tokenStore:   tokenStore,
		providers:    providers,
		sessionStore: sessionStore,
	}
}

type api struct {
	tokenService token.Service
	tokenStore   token.Store
	providers    []security.Provider

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
			panic(err)
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

	err = a.tokenStore.Save(c.Request.Context(), token.Token{
		ID:       accessToken,
		User:     user.ID,
		Provider: provider,
		Created:  types.CpsTime{Time: time.Now()},
		Expired:  types.CpsTime{Time: expiresAt},
	})
	if err != nil {
		panic(err)
	}

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
// @Success 204
// @Router /logout [post]
func (a *api) Logout(c *gin.Context) {
	header := c.GetHeader(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	tokenString := strings.TrimSpace(header[len(bearerPrefix):])
	ok, err := a.tokenStore.Delete(c.Request.Context(), tokenString)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

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
	count, err := a.tokenStore.Count(c.Request.Context())
	if err != nil {
		panic(err)
	}

	// todo : remove after session delete
	sessionCount, err := a.sessionStore.GetActiveSessionsCount(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, loggedUserCountResponse{
		Count: count + sessionCount,
	})
}
