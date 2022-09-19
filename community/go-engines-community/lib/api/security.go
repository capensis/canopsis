package api

import (
	"net/http"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/cas"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/saml"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
	"github.com/rs/zerolog"
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
	GetTokenService() apisecurity.TokenService
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
	config       *libsecurity.Config
	dbClient     mongo.DbClient
	sessionStore libsession.Store
	enforcer     libsecurity.Enforcer
	logger       zerolog.Logger

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
		config:       config,
		dbClient:     dbClient,
		sessionStore: sessionStore,
		enforcer:     enforcer,
		logger:       logger,

		cookieOptions: cookieOptions,

		apiConfigProvider: apiConfigProvider,
	}
}

func (s *security) GetHttpAuthProviders() []libsecurity.HttpProvider {
	res := make([]libsecurity.HttpProvider, 0)

	for _, v := range s.config.Security.AuthProviders {
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

	for _, v := range s.config.Security.AuthProviders {
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
	for _, v := range s.config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodCas:
			p := httpprovider.NewCasProvider(
				http.DefaultClient,
				s.newConfigProvider(),
				s.newUserProvider(),
			)
			router.GET("/cas/login", cas.SessionLoginHandler(s.newConfigProvider()))
			router.GET("/cas/loggedin", cas.SessionCallbackHandler(p, s.enforcer, s.sessionStore))
			router.GET("/api/v4/cas/login", cas.LoginHandler(s.newConfigProvider()))
			router.GET("/api/v4/cas/loggedin", cas.CallbackHandler(p, s.enforcer, s.GetTokenService()))
		case libsecurity.AuthMethodSaml:
			sp, err := saml.NewServiceProvider(s.newUserProvider(), client.Collection(mongo.RightsMongoCollection), s.sessionStore,
				s.enforcer, s.config, s.GetTokenService(), s.logger)
			if err != nil {
				s.logger.Err(err).Msg("RegisterCallbackRoutes: NewServiceProvider error")
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
		middleware.SessionAuth(s.dbClient, s.sessionStore),
	}
}

func (s *security) GetFileAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Auth([]libsecurity.HttpProvider{
			httpprovider.NewCookieProvider(s.getJwtTokenService(), s.GetTokenStore(),
				s.newUserProvider(), s.cookieOptions.FileAccessName, s.logger),
		}),
	}
}

func (s *security) GetSessionStore() libsession.Store {
	return s.sessionStore
}

func (s *security) GetConfig() libsecurity.Config {
	return *s.config
}

func (s *security) GetPasswordEncoder() password.Encoder {
	return password.NewSha1Encoder()
}

func (s *security) GetTokenService() apisecurity.TokenService {
	return apisecurity.NewTokenService(s.dbClient, s.getJwtTokenService(), s.GetTokenStore())
}

func (s *security) GetTokenStore() token.Store {
	return token.NewMongoStore(s.dbClient, s.logger)
}

func (s *security) GetTokenProvider() libsecurity.TokenProvider {
	return tokenprovider.NewTokenProvider(s.getJwtTokenService(), s.GetTokenStore(), s.newUserProvider(), s.logger)
}

func (s *security) GetCookieOptions() CookieOptions {
	return s.cookieOptions
}

func (s *security) getJwtTokenService() token.Service {
	secretKey := os.Getenv(JwtSecretEnv)
	return token.NewJwtService([]byte(secretKey), canopsis.AppName, s.apiConfigProvider)
}

func (s *security) newUserProvider() libsecurity.UserProvider {
	return userprovider.NewMongoProvider(s.dbClient)
}

func (s *security) newConfigProvider() libsecurity.ConfigProvider {
	return configprovider.NewMongoProvider(s.dbClient)
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
