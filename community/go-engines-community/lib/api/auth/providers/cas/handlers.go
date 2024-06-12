package cas

import (
	"fmt"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
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
func CallbackHandler(p libsecurity.HttpProvider, enforcer libsecurity.Enforcer, tokenService security.TokenService, maintenanceAdapter config.MaintenanceAdapter) gin.HandlerFunc {
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

		maintenanceConf, err := maintenanceAdapter.GetConfig(c)
		if err != nil {
			panic(err)
		}

		if maintenanceConf.Enabled {
			ok, err = enforcer.Enforce(user.ID, security.PermMaintenance, model.PermissionCan)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
				return
			}
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
