package cas

import (
	"fmt"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

// LoginHandler redirects to CAS login url and saves referer url to service url.
func LoginHandler(config libsecurity.CasConfig, errorResponder httperror.Responder) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := validation.Bind(c, &request); err != nil {
			errorResponder.Respond(c, err)

			return
		}

		casUrl, err := url.Parse(config.LoginUrl)
		if err != nil {
			errorResponder.Respond(c, err)

			return
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
func CallbackHandler(
	p libsecurity.HttpProvider,
	enforcer libsecurity.Enforcer,
	tokenService security.TokenService,
	maintenanceAdapter config.MaintenanceAdapter,
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := validation.Bind(c, &request); err != nil {
			errorResponder.Respond(c, err)

			return
		}

		redirectUrl, err := url.Parse(request.Redirect)
		if err != nil {
			errorResponder.Respond(c, fmt.Errorf("parse redirect url error: %w", err))

			return
		}

		q := redirectUrl.Query()

		user, err, ok := p.Auth(c.Request)
		if err != nil {
			q.Set("errorMessage", err.Error())
			redirectUrl.RawQuery = q.Encode()

			c.Redirect(http.StatusPermanentRedirect, redirectUrl.String())
		}

		if !ok || user == nil {
			errorResponder.Respond(c, httperror.ErrUnauthorized)

			return
		}

		err = enforcer.LoadPolicy()
		if err != nil {
			errorResponder.Respond(c, fmt.Errorf("reload enforcer error: %w", err))

			return
		}

		maintenanceConf, err := maintenanceAdapter.GetConfig(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		if maintenanceConf.Enabled {
			ok, err = enforcer.Enforce(user.ID, security.PermMaintenance, model.PermissionCan)
			if err != nil {
				errorResponder.Respond(c, err)

				return
			}

			if !ok {
				errorResponder.Respond(c, httperror.ErrMaintenance)

				return
			}
		}

		accessToken, err := tokenService.Create(c, *user, libsecurity.AuthMethodCas)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		q.Set("access_token", accessToken)
		redirectUrl.RawQuery = q.Encode()

		c.Redirect(http.StatusPermanentRedirect, redirectUrl.String())
	}
}
