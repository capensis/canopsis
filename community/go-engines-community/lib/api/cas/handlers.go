package cas

import (
	"fmt"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// LoginHandler redirects to CAS login url and saves referer url to service url.
func LoginHandler(config libsecurity.CasConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		casUrl, err := url.Parse(config.LoginUrl)
		if err != nil {
			panic(err)
		}

		service := fmt.Sprintf("%s?redirect=%s&service=%s",
			request.Service, request.Redirect, request.Service)
		q := casUrl.Query()
		q.Set("service", service)
		casUrl.RawQuery = q.Encode()

		c.Redirect(http.StatusPermanentRedirect, casUrl.String())
	}
}

// CallbackHandler validates CAS ticket, creates access token and redirects to referer url.
func CallbackHandler(p libsecurity.HttpProvider, enforcer libsecurity.Enforcer, tokenService security.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		user, err, ok := p.Auth(c.Request)
		if err != nil {
			panic(err)
		}

		if !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		err = enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		accessToken, err := tokenService.Create(c, *user, libsecurity.AuthMethodCas)
		if err != nil {
			panic(err)
		}

		redirectUrl, err := url.Parse(request.Redirect)
		if err != nil {
			panic(fmt.Errorf("parse redirect url error: %w", err))
		}

		q := redirectUrl.Query()
		q.Set("access_token", accessToken)
		redirectUrl.RawQuery = q.Encode()

		c.Redirect(http.StatusPermanentRedirect, redirectUrl.String())
	}
}

// SessionLoginHandler redirects to CAS login url and saves referer url to service url.
func SessionLoginHandler(config libsecurity.CasConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		casUrl, err := url.Parse(config.LoginUrl)
		if err != nil {
			panic(err)
		}

		service := fmt.Sprintf("%s?redirect=%s&service=%s",
			request.Service, request.Redirect, request.Service)
		q := casUrl.Query()
		q.Set("service", service)
		casUrl.RawQuery = q.Encode()

		c.Redirect(http.StatusPermanentRedirect, casUrl.String())
	}
}

// SessionCallbackHandler validates CAS ticket, inits session and redirects to referer url.
func SessionCallbackHandler(p libsecurity.HttpProvider, enforcer libsecurity.Enforcer, sessionStore libsession.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		user, err, ok := p.Auth(c.Request)
		if err != nil {
			panic(err)
		}

		if !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		session := getSession(c, sessionStore)
		session.Values["user"] = user.ID
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			panic(err)
		}

		err = enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		c.Redirect(http.StatusPermanentRedirect, request.Redirect)
	}
}

func getSession(c *gin.Context, sessionStore libsession.Store) *sessions.Session {
	session, err := sessionStore.Get(c.Request, libsecurity.SessionKey)
	if err != nil {
		panic(err)
	}

	return session
}
