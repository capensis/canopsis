package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/roleprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

const (
	randomBytesNumber = 16
)

type Provider interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
}

type provider struct {
	roleProvider       security.RoleProvider
	userProvider       security.UserProvider
	tokenService       apisecurity.TokenService
	oauth2Config       oauth2.Config
	config             security.OAuth2ProviderConfig
	store              sessions.Store
	maintenanceAdapter config.MaintenanceAdapter
	enforcer           security.Enforcer
	logger             zerolog.Logger

	// only for OpenID type of providers
	oidcVerifier *oidc.IDTokenVerifier
	provider     *oidc.Provider

	name   string
	source string

	maxResponseSize int64
}

func NewProvider(
	ctx context.Context,
	name string,
	roleValidator security.RoleProvider,
	config security.OAuth2ProviderConfig,
	store sessions.Store,
	userProvider security.UserProvider,
	maintenanceAdapter config.MaintenanceAdapter,
	enforcer security.Enforcer,
	tokenService apisecurity.TokenService,
	maxResponseSize int64,
) (Provider, error) {
	p := &provider{
		name:               name,
		roleProvider:       roleValidator,
		userProvider:       userProvider,
		maintenanceAdapter: maintenanceAdapter,
		tokenService:       tokenService,
		enforcer:           enforcer,
		store:              store,
		oauth2Config: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthURL,
				TokenURL: config.TokenURL,
			},
			Scopes: config.Scopes,
		},
		config:          config,
		source:          name,
		maxResponseSize: maxResponseSize,
	}

	if config.OpenID {
		var err error

		p.provider, err = oidc.NewProvider(ctx, config.Issuer)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s authentication provider: %w", name, err)
		}

		var claims oidcClaims
		if err := p.provider.Claims(&claims); err != nil {
			return nil, fmt.Errorf("failed to decode %s provider claims: %w", name, err)
		}

		if !claims.ValidateScopes(config.Scopes) {
			return nil, fmt.Errorf("scopes are not supported for %s provider", name)
		}

		p.oidcVerifier = p.provider.Verifier(&oidc.Config{ClientID: config.ClientID})
		p.oauth2Config.Endpoint = p.provider.Endpoint()
	}

	return p, nil
}

func (p *provider) Login(c *gin.Context) {
	request := loginRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	session, err := p.store.Get(c.Request, oauthSessionPrefix+p.name)
	if err != nil {
		panic(err)
	}

	options := make([]oauth2.AuthCodeOption, 0, 2)

	if p.config.PKCE {
		verifier := oauth2.GenerateVerifier()
		session.Values["pkce_verifier"] = verifier

		options = append(options, oauth2.S256ChallengeOption(verifier))
	}

	if p.config.OpenID {
		nonce, err := utils.RandBase64String(randomBytesNumber)
		if err != nil {
			panic(fmt.Errorf("failed to generate nonce: %w", err))
		}

		session.Values["nonce"] = nonce
		options = append(options, oidc.Nonce(nonce))
	}

	state, err := utils.RandBase64String(randomBytesNumber)
	if err != nil {
		panic(fmt.Errorf("failed to generate state: %w", err))
	}

	session.Values["state"] = state
	session.Values["redirect"] = request.Redirect
	session.Options.MaxAge = 300

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		panic(fmt.Errorf("failed to save session: %w", err))
	}

	c.Redirect(http.StatusPermanentRedirect, p.oauth2Config.AuthCodeURL(state, options...))
}

func (p *provider) Callback(c *gin.Context) {
	session, err := p.store.Get(c.Request, oauthSessionPrefix+p.name)
	if err != nil {
		panic(err)
	}

	// expire auth session
	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		panic(err)
	}

	redirectRaw, ok := session.Values["redirect"]
	if !ok {
		panic(errors.New("redirect url not found"))
	}

	redirectStr, ok := redirectRaw.(string)
	if !ok {
		panic(errors.New("redirect url should be a string"))
	}

	redirectUrl, err := url.Parse(redirectStr)
	if err != nil {
		panic(fmt.Errorf("redirect string is not an URL: %w", err))
	}

	stateRaw, ok := session.Values["state"]
	if !ok {
		panic(errors.New("state not found"))
	}

	state, ok := stateRaw.(string)
	if !ok {
		panic(errors.New("state should be a string"))
	}

	if c.Query("state") != state {
		p.logger.Err(errors.New("state did not match")).Str("provider", p.name).Str("session_state", state).Str("query_state", c.Query("state")).Msg("oauth2 callback error")
		p.errorRedirect(c, redirectUrl, "state did not match")

		return
	}

	options := make([]oauth2.AuthCodeOption, 0, 1)

	if p.config.PKCE {
		verifierRaw, ok := session.Values["pkce_verifier"]
		if !ok {
			panic(errors.New("pkce_verifier not found"))
		}

		verifier, ok := verifierRaw.(string)
		if !ok {
			panic(errors.New("pkce_verifier should be a string"))
		}

		options = append(options, oauth2.VerifierOption(verifier))
	}

	if errorMsg, ok := c.GetQuery("error"); ok {
		if errorDesc, ok := c.GetQuery("error_description"); ok {
			errorMsg = errorDesc
		}

		p.errorRedirect(c, redirectUrl, errorMsg)
		return
	}

	code, ok := c.GetQuery("code")
	if !ok {
		p.logger.Err(errors.New("code is empty")).Str("provider", p.name).Msg("oauth2 callback error")
		p.errorRedirect(c, redirectUrl, "code is empty")

		return
	}

	oauth2Token, err := p.oauth2Config.Exchange(c, code, options...)
	if err != nil {
		panic(err)
	}

	var userID string
	var userInfo map[string]any

	if !p.config.OpenID {
		userID, userInfo, err = p.getUserInfoOAuth2(c, oauth2Token)
		if err != nil {
			panic(fmt.Errorf("failed to get user info: %w", err))
		}
	} else {
		nonceRaw, ok := session.Values["nonce"]
		if !ok {
			panic(errors.New("nonce not found"))
		}

		nonce, ok := nonceRaw.(string)
		if !ok {
			panic(errors.New("nonce should be a string"))
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			p.logger.Err(errors.New("id_token is not a string")).Str("provider", p.name).Interface("id_token", oauth2Token.Extra("id_token")).Msg("oauth2 callback error")
			p.errorRedirect(c, redirectUrl, "id_token is not a string")

			return
		}

		idToken, err := p.oidcVerifier.Verify(c, rawIDToken)
		if err != nil {
			p.logger.Err(fmt.Errorf("id_token is not valid: %w", err)).Str("provider", p.name).Str("id_token", rawIDToken).Msg("oauth2 callback error")
			p.errorRedirect(c, redirectUrl, "id_token is not valid")

			return
		}

		if idToken.Nonce != nonce {
			p.logger.Err(errors.New("nonce did not match")).Str("provider", p.name).Str("session_nonce", nonce).Str("token_nonce", idToken.Nonce).Msg("oauth2 callback error")
			p.errorRedirect(c, redirectUrl, "nonce did not match")

			return
		}

		userID, userInfo, err = p.getUserInfoOpenID(c, oauth2Token, idToken)
		if err != nil {
			panic(fmt.Errorf("failed to get user info: %w", err))
		}
	}

	if userID == "" {
		p.logger.Err(errors.New("userID cannot be empty")).Str("provider", p.name).Msg("oauth2 callback error")
		p.errorRedirect(c, redirectUrl, "userID cannot be empty")

		return
	}

	user, err := p.userProvider.FindByExternalSource(c, userID, p.source)
	if err != nil {
		panic(fmt.Errorf("failed to get user by external source: %w", err))
	}

	if user == nil {
		user, ok = p.createUser(c, redirectUrl, userID, userInfo)
		if !ok {
			return
		}
	}

	err = p.enforcer.LoadPolicy()
	if err != nil {
		panic(fmt.Errorf("reload enforcer error: %w", err))
	}

	maintenanceConf, err := p.maintenanceAdapter.GetConfig(c)
	if err != nil {
		panic(fmt.Errorf("failed to get maintenance config: %w", err))
	}

	if maintenanceConf.Enabled {
		ok, err = p.enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
		if err != nil {
			panic(fmt.Errorf("enforcer failed: %w", err))
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
			return
		}
	}

	var accessToken string
	if p.config.ExpirationInterval != "" {
		accessToken, err = p.tokenService.Create(c, *user, p.source)
	} else {
		accessToken, err = p.tokenService.CreateWithExpiration(c, *user, p.source, oauth2Token.Expiry)
	}

	if err != nil {
		panic(fmt.Errorf("failed to create token: %w", err))
	}

	query := redirectUrl.Query()
	query.Set("access_token", accessToken)
	redirectUrl.RawQuery = query.Encode()

	c.Redirect(http.StatusPermanentRedirect, redirectUrl.String())
}

func (p *provider) createUser(c *gin.Context, redirectUrl *url.URL, subj string, userInfo map[string]any) (*security.User, bool) {
	roles, err := p.roleProvider.GetValidRoleIDs(c, p.getAssocArrayAttribute(userInfo, "role", []string{}), p.config.DefaultRole)
	if err != nil {
		roleNotFoundError := roleprovider.ProviderError{}
		if errors.As(err, &roleNotFoundError) {
			p.logger.Err(roleNotFoundError).Msg("user registration failed")
			p.errorRedirect(c, redirectUrl, roleNotFoundError.Error())

			return nil, false
		}

		panic(err)
	}

	user := &security.User{
		Name:       p.getAssocAttribute(userInfo, "name", subj),
		Roles:      roles,
		IsEnabled:  true,
		ExternalID: subj,
		Source:     p.source,
		Firstname:  p.getAssocAttribute(userInfo, "firstname", ""),
		Lastname:   p.getAssocAttribute(userInfo, "lastname", ""),
		Email:      p.getAssocAttribute(userInfo, "email", ""),
	}

	err = p.userProvider.Save(c, user)
	if err != nil {
		p.logger.Err(err).Msg("user registration failed")
		panic(fmt.Errorf("cannot save user: %w", err))
	}

	return user, true
}

func (p *provider) getAssocAttribute(userInfo map[string]any, name, defaultValue string) string {
	if path, ok := p.config.AttributesMap[name]; ok {
		if valueRaw, ok := userInfo[path]; ok {
			if value, ok := valueRaw.(string); ok {
				return value
			}
		}
	}

	return defaultValue
}

func (p *provider) getAssocArrayAttribute(userInfo map[string]any, name string, defaultValue []string) []string {
	if path, ok := p.config.AttributesMap[name]; ok {
		if valueRaw, ok := userInfo[path]; ok {
			switch value := valueRaw.(type) {
			case []any:
				if strSlice, ok := utils.InterfaceSliceToStringSlice(value); ok && len(strSlice) > 0 {
					return strSlice
				}
			case []string:
				if len(value) != 0 {
					return value
				}
			case string:
				if value != "" {
					return []string{value}
				}
			}
		}
	}

	return defaultValue
}

func (p *provider) getUserInfoOAuth2(c context.Context, token *oauth2.Token) (string, map[string]any, error) {
	resp, err := p.oauth2Config.Client(c, token).Get(p.config.UserURL) //nolint:noctx
	if err != nil {
		return "", nil, errors.New("failed to get user request")
	}

	defer resp.Body.Close()

	flatResp, err := libhttp.FlattenResponse(nil, resp, p.maxResponseSize)
	if err != nil {
		return "", nil, errors.New("failed to flatten response")
	}

	userInfo := make(map[string]any, len(flatResp))
	for k, v := range flatResp {
		userInfo["user."+k] = v
	}

	userIDRaw, ok := userInfo[p.config.UserID]
	if !ok {
		return "", userInfo, errors.New("can't find user_id")
	}

	var userID string

	switch userIDRaw := userIDRaw.(type) {
	case string:
		userID = userIDRaw
	case int:
		userID = strconv.Itoa(userIDRaw)
	default:
		return "", userInfo, errors.New("user_id should be string or int")
	}

	return userID, userInfo, nil
}

func (p *provider) getUserInfoOpenID(c context.Context, token *oauth2.Token, idToken *oidc.IDToken) (string, map[string]any, error) {
	userInfo := make(map[string]any)
	claims := make(map[string]any)

	err := idToken.Claims(&claims)
	if err != nil {
		return "", userInfo, fmt.Errorf("failed to decode token claims: %w", err)
	}

	for k, v := range claims {
		userInfo["token."+k] = v
	}

	userInfoResp, err := p.provider.UserInfo(c, oauth2.StaticTokenSource(token))
	if err != nil {
		return "", userInfo, fmt.Errorf("failed to get userinfo: %w", err)
	}

	err = userInfoResp.Claims(&claims)
	if err != nil {
		return "", userInfo, fmt.Errorf("failed to decode user info claims: %w", err)
	}

	for k, v := range claims {
		userInfo["user."+k] = v
	}

	return idToken.Subject, userInfo, nil
}

func (p *provider) errorRedirect(c *gin.Context, relayUrl *url.URL, errorMessage string) {
	query := relayUrl.Query()
	query.Set("errorMessage", errorMessage)
	relayUrl.RawQuery = query.Encode()

	c.Redirect(http.StatusPermanentRedirect, relayUrl.String())
}
