package api

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/saml"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/configprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/httpprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/provider"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/tokenprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"os"
	"time"
)

const JwtSecretEnv = "CPS_JWT_SECRET"

// Security is used to init auth methods by config.
type Security interface {
	// GetHttpAuthProviders creates http providers which authenticates each API request.
	GetHttpAuthProviders() []libsecurity.HttpProvider
	// GetAuthProviders creates providers which are used in auth API request.
	GetAuthProviders() []libsecurity.Provider
	// RegisterCallbackRoutes registers callback routes for auth methods.
	RegisterCallbackRoutes(router gin.IRouter, client mongo.DbClient)
	// GetAuthMiddleware returns corresponding config auth middlewares.
	GetAuthMiddleware() []gin.HandlerFunc
	// GetFileAuthMiddleware returns auth middlewares for files.
	GetFileAuthMiddleware() []gin.HandlerFunc
	GetSessionStore() libsession.Store
	GetConfig() libsecurity.Config
	GetPasswordEncoder() password.Encoder
	GetTokenService() token.Service
	GetTokenStore() token.Store
	GetTokenProvider() libsecurity.TokenProvider
	GetCookieOptions() CookieOptions
}

type CookieOptions struct {
	FileAccessName string
	MaxAge         int
	Secure         bool
}

type security struct {
	Config       *libsecurity.Config
	DbClient     mongo.DbClient
	SessionStore libsession.Store
	enforcer     libsecurity.Enforcer
	Logger       zerolog.Logger

	apiConfigProvider config.ApiConfigProvider

	cookieOptions CookieOptions
}

// NewSecurity creates new security.
func NewSecurity(
	config *libsecurity.Config,
	dbClient mongo.DbClient,
	sessionStore libsession.Store,
	enforcer libsecurity.Enforcer,
	apiConfigProvider config.ApiConfigProvider,
	cookieOptions CookieOptions,
	logger zerolog.Logger,
) Security {
	return &security{
		Config:       config,
		DbClient:     dbClient,
		SessionStore: sessionStore,
		enforcer:     enforcer,
		Logger:       logger,

		cookieOptions: cookieOptions,

		apiConfigProvider: apiConfigProvider,
	}
}

func (s *security) GetHttpAuthProviders() []libsecurity.HttpProvider {
	res := make([]libsecurity.HttpProvider, 0)

	for _, v := range s.Config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodBasic:
			baseProvider := s.newBaseAuthProvider()
			res = append(res, httpprovider.NewBasicProvider(baseProvider))
			res = append(res, httpprovider.NewBearerProvider(s.GetTokenProvider()))
		case libsecurity.AuthMethodApiKey:
			res = append(res, httpprovider.NewApikeyProvider(s.newUserProvider()))
		case libsecurity.AuthMethodLdap:
			ldapProvider := s.newLdapAuthProvider()
			res = append(res, httpprovider.NewQueryBasicProvider(ldapProvider))
		}
	}

	return res
}

func (s *security) GetAuthProviders() []libsecurity.Provider {
	res := make([]libsecurity.Provider, 0)

	for _, v := range s.Config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodBasic:
			res = append(res, s.newBaseAuthProvider())
		case libsecurity.AuthMethodLdap:
			res = append(res, s.newLdapAuthProvider())
		}
	}

	return res
}

func (s *security) RegisterCallbackRoutes(router gin.IRouter, client mongo.DbClient) {
	for _, v := range s.Config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodCas:
			p := httpprovider.NewCasProvider(
				http.DefaultClient,
				s.newConfigProvider(),
				s.newUserProvider(),
			)
			router.GET("/cas/login", s.casSessionLoginHandler())
			router.GET("/cas/loggedin", s.casSessionCallbackHandler(p))
			router.GET("/api/v4/cas/login", s.casLoginHandler())
			router.GET("/api/v4/cas/loggedin", s.casCallbackHandler(p))
		case libsecurity.AuthMethodSaml:
			sp, err := saml.NewServiceProvider(s.newUserProvider(), client.Collection(mongo.RightsMongoCollection), s.SessionStore,
				s.enforcer, s.Config, s.GetTokenService(), s.GetTokenStore(), s.Logger)
			if err != nil {
				s.Logger.Err(err).Msg("RegisterCallbackRoutes: NewServiceProvider error")
				panic(err)
			}

			router.GET("/saml/metadata", sp.SamlMetadataHandler())
			router.GET("/saml/auth", sp.SamlAuthHandler())
			router.POST("/saml/acs", sp.SamlAcsHandler())
			router.GET("/saml/slo", sp.SamlSloHandler())
			router.GET("/api/v4/saml/metadata", sp.SamlMetadataHandler())
			router.GET("/api/v4/saml/auth", sp.SamlAuthHandler())
			router.POST("/api/v4/saml/acs", sp.SamlAcsHandler())
			router.GET("/api/v4/saml/slo", sp.SamlSloHandler())
		}
	}
}

func (s *security) GetAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Auth(s.GetHttpAuthProviders()),
		middleware.SessionAuth(s.DbClient, s.SessionStore),
	}
}

func (s *security) GetFileAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Auth([]libsecurity.HttpProvider{
			httpprovider.NewCookieProvider(s.GetTokenService(), s.GetTokenStore(),
				s.newUserProvider(), s.cookieOptions.FileAccessName, s.Logger),
		}),
	}
}

func (s *security) GetSessionStore() libsession.Store {
	return s.SessionStore
}

func (s *security) GetConfig() libsecurity.Config {
	return *s.Config
}

func (s *security) GetPasswordEncoder() password.Encoder {
	return password.NewSha1Encoder()
}

func (s *security) GetTokenService() token.Service {
	secretKey := os.Getenv(JwtSecretEnv)

	return token.NewJwtService([]byte(secretKey), canopsis.AppName, s.apiConfigProvider)
}
func (s *security) GetTokenStore() token.Store {
	return token.NewMongoStore(s.DbClient, s.Logger)
}
func (s *security) GetTokenProvider() libsecurity.TokenProvider {
	return tokenprovider.NewTokenProvider(s.GetTokenService(), s.GetTokenStore(), s.newUserProvider(), s.Logger)
}

func (s *security) GetCookieOptions() CookieOptions {
	return s.cookieOptions
}

type casLoginRequest struct {
	// Redirect is front-end url to redirect back after authentication.
	Redirect string `form:"redirect"`
	// Service is proxy url to callback handler to set cookie for front-end app.
	Service string `form:"service"`
}

// casSessionLoginHandler redirects to CAS login url and saves referer url to service url.
func (s *security) casSessionLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		casConfig, err := s.newConfigProvider().LoadCasConfig(c.Request.Context())
		if err != nil {
			panic(err)
		}

		casUrl, err := url.Parse(casConfig.LoginUrl)
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

// casSessionCallbackHandler validates CAS ticket, inits session and redirects to referer url.
func (s *security) casSessionCallbackHandler(p libsecurity.HttpProvider) gin.HandlerFunc {
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

		session := s.getSession(c)
		session.Values["user"] = user.ID
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			panic(err)
		}

		err = s.enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		c.Redirect(http.StatusPermanentRedirect, request.Redirect)
	}
}

// casLoginHandler redirects to CAS login url and saves referer url to service url.
func (s *security) casLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := casLoginRequest{}

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		casConfig, err := s.newConfigProvider().LoadCasConfig(c.Request.Context())
		if err != nil {
			panic(err)
		}

		casUrl, err := url.Parse(casConfig.LoginUrl)
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

// casCallbackHandler validates CAS ticket, creates access token and redirects to referer url.
func (s *security) casCallbackHandler(p libsecurity.HttpProvider) gin.HandlerFunc {
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

		err = s.enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		accessToken, expiresAt, err := s.GetTokenService().GenerateToken(user.ID)
		if err != nil {
			panic(err)
		}

		now := time.Now()
		err = s.GetTokenStore().Save(c.Request.Context(), token.Token{
			ID:       accessToken,
			User:     user.ID,
			Provider: libsecurity.AuthMethodCas,
			Created:  types.CpsTime{Time: now},
			Expired:  types.CpsTime{Time: expiresAt},
		})
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

func (s *security) newUserProvider() libsecurity.UserProvider {
	return userprovider.NewMongoProvider(s.DbClient)
}

func (s *security) newConfigProvider() libsecurity.ConfigProvider {
	return configprovider.NewMongoProvider(s.DbClient)
}

func (s *security) newBaseAuthProvider() libsecurity.Provider {
	return provider.NewBaseProvider(s.newUserProvider(), s.GetPasswordEncoder())
}

func (s *security) newLdapAuthProvider() libsecurity.Provider {
	return provider.NewLdapProvider(
		s.newConfigProvider(),
		s.newUserProvider(),
		provider.NewLdapDialer(),
		s.enforcer,
	)
}

func (s *security) getSession(c *gin.Context) *sessions.Session {
	session, err := s.SessionStore.Get(c.Request, libsecurity.SessionKey)
	if err != nil {
		panic(err)
	}

	return session
}
