package api

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/saml"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/configprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/httpprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/provider"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
)

// Security is used to init auth methods by config.
type Security interface {
	// GetHttpAuthProviders creates http providers which authenticates each API request.
	GetHttpAuthProviders() []libsecurity.HttpProvider
	// GetAuthProviders creates providers which are used in auth API request.
	GetAuthProviders() []libsecurity.Provider
	// RegisterCallbackRoutes registers callback routes for auth methods.
	RegisterCallbackRoutes(router gin.IRouter)
	// GetAuthMiddleware returns corresponding config auth middlewares.
	GetAuthMiddleware() []gin.HandlerFunc
	// GetWebsocketAuthMiddleware returns auth middlewares for websocket.
	GetWebsocketAuthMiddleware() []gin.HandlerFunc
	GetSessionStore() libsession.Store
	GetConfig() libsecurity.Config
	GetPasswordEncoder() password.Encoder
}

type security struct {
	Config       *libsecurity.Config
	DbClient     mongo.DbClient
	SessionStore libsession.Store
	Logger       zerolog.Logger
}

// NewSecurity creates new security.
func NewSecurity(
	config *libsecurity.Config,
	dbClient mongo.DbClient,
	sessionStore libsession.Store,
	logger zerolog.Logger,
) Security {
	return &security{
		Config:       config,
		DbClient:     dbClient,
		SessionStore: sessionStore,
		Logger:       logger,
	}
}

func (s *security) GetHttpAuthProviders() []libsecurity.HttpProvider {
	res := make([]libsecurity.HttpProvider, 0)

	for _, v := range s.Config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodBasic:
			baseProvider := s.newBaseAuthProvider()
			res = append(res, httpprovider.NewBasicProvider(baseProvider))
		case libsecurity.AuthMethodApiKey:
			res = append(res, httpprovider.NewApikeyProvider(s.newUserProvider()))
		case libsecurity.AuthMethodLdap:
			ldapProvider := s.newLdapAuthProvider()
			res = append(res, httpprovider.NewQueryProvider(ldapProvider))
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

func (s *security) RegisterCallbackRoutes(router gin.IRouter) {
	for _, v := range s.Config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodCas:
			p := httpprovider.NewCasProvider(
				http.DefaultClient,
				s.newConfigProvider(),
				s.newUserProvider(),
			)
			router.GET("/cas/login", s.casLoginHandler())
			router.GET("/cas/loggedin", s.casCallbackHandler(p))
		case libsecurity.AuthMethodSaml:
			sp, err := saml.NewServiceProvider(s.newUserProvider(), s.SessionStore, s.Config, s.Logger)
			if err != nil {
				s.Logger.Err(err).Msg("RegisterCallbackRoutes: NewServiceProvider error")
				panic(err)
			}

			router.GET("/saml/metadata", sp.SamlMetadataHandler())
			router.GET("/saml/auth", sp.SamlAuthHandler())
			router.POST("/saml/acs", sp.SamlAcsHandler())
			router.GET("/saml/slo", sp.SamlSloHandler())
		}
	}
}

func (s *security) GetAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Auth(s.GetHttpAuthProviders()),
		middleware.SessionAuth(s.DbClient, s.SessionStore),
	}
}

func (s *security) GetWebsocketAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Auth([]libsecurity.HttpProvider{
			httpprovider.NewApikeyProvider(s.newUserProvider()),
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

type casLoginRequest struct {
	// Redirect is front-end url to redirect back after authentication.
	Redirect string `form:"redirect"`
	// Service is proxy url to callback handler to set cookie for front-end app.
	Service string `form:"service"`
}

// casLoginHandler redirects to CAS login url and saves referer url to session.
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

// casCallbackHandler validates CAS ticket, inits session and redirects to referer url.
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

		session := s.getSession(c)
		session.Values["user"] = user.ID
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusPermanentRedirect, request.Redirect)
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
	)
}

func (s *security) getSession(c *gin.Context) *sessions.Session {
	session, err := s.SessionStore.Get(c.Request, libsecurity.SessionKey)
	if err != nil {
		panic(err)
	}

	return session
}
