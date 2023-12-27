package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/cas"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/saml"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/httpprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/provider"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/sharetoken"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/tokenprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const JwtSecretEnv = "CPS_JWT_SECRET" //nolint:gosec
const sessionStoreSessionMaxAge = 24 * time.Hour

// Security is used to init auth methods by config.
type Security interface {
	// GetHttpAuthProviders creates http providers which authenticates each API request.
	GetHttpAuthProviders() []libsecurity.HttpProvider
	// GetAuthProviders creates providers which are used in auth API request.
	GetAuthProviders() []libsecurity.Provider
	// RegisterCallbackRoutes registers callback routes for auth methods.
	RegisterCallbackRoutes(ctx context.Context, router gin.IRouter, client mongo.DbClient)
	// GetAuthMiddleware returns corresponding config auth middlewares.
	GetAuthMiddleware() []gin.HandlerFunc
	// GetFileAuthMiddleware returns auth middleware for files.
	GetFileAuthMiddleware() gin.HandlerFunc
	GetSessionStore() libsession.Store
	GetConfig() libsecurity.Config
	GetPasswordEncoder() password.Encoder
	GetTokenService() apisecurity.TokenService
	GetTokenGenerator() token.Generator
	GetTokenProviders() []libsecurity.TokenProvider
	GetCookieOptions() CookieOptions
}

type CookieOptions struct {
	FileAccessName string
	MaxAge         int
}

func DefaultCookieOptions() CookieOptions {
	return CookieOptions{
		FileAccessName: "token",
		MaxAge:         int(sessionStoreSessionMaxAge.Seconds()),
	}
}

type security struct {
	config       libsecurity.Config
	dbClient     mongo.DbClient
	sessionStore libsession.Store
	enforcer     libsecurity.Enforcer
	logger       zerolog.Logger

	apiConfigProvider  config.ApiConfigProvider
	maintenanceAdapter config.MaintenanceAdapter

	cookieOptions CookieOptions
}

// NewSecurity creates new security.
func NewSecurity(
	config libsecurity.Config,
	dbClient mongo.DbClient,
	sessionStore libsession.Store,
	enforcer libsecurity.Enforcer,
	apiConfigProvider config.ApiConfigProvider,
	maintenanceAdapter config.MaintenanceAdapter,
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

		apiConfigProvider:  apiConfigProvider,
		maintenanceAdapter: maintenanceAdapter,
	}
}

func (s *security) GetHttpAuthProviders() []libsecurity.HttpProvider {
	res := make([]libsecurity.HttpProvider, 0)

	for _, v := range s.config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodBasic:
			baseProvider := s.newBaseAuthProvider()
			res = append(res, httpprovider.NewBasicProvider(baseProvider))
			res = append(res, httpprovider.NewBearerProvider(s.GetTokenProviders()))
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

func (s *security) RegisterCallbackRoutes(ctx context.Context, router gin.IRouter, client mongo.DbClient) {
	for _, v := range s.config.Security.AuthProviders {
		switch v {
		case libsecurity.AuthMethodCas:
			casConfig := s.config.Security.Cas
			p := httpprovider.NewCasProvider(
				http.DefaultClient,
				casConfig,
				s.newUserProvider(),
			)
			router.GET("/cas/login", cas.SessionLoginHandler(casConfig))
			router.GET("/cas/loggedin", cas.SessionCallbackHandler(p, s.enforcer, s.sessionStore))
			router.GET("/api/v4/cas/login", cas.LoginHandler(casConfig))
			router.GET("/api/v4/cas/loggedin", cas.CallbackHandler(p, s.enforcer, s.GetTokenService(), s.maintenanceAdapter)) //nolint: contextcheck
		case libsecurity.AuthMethodSaml:
			sp, err := saml.NewServiceProvider(ctx, s.newUserProvider(), client.Collection(mongo.RoleCollection), s.sessionStore,
				s.enforcer, s.config, s.GetTokenService(), s.maintenanceAdapter, s.logger)
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
		middleware.Auth(s.GetHttpAuthProviders(), s.maintenanceAdapter, s.enforcer),
		middleware.SessionAuth(s.dbClient, s.apiConfigProvider, s.sessionStore),
	}
}

func (s *security) GetFileAuthMiddleware() gin.HandlerFunc {
	return middleware.Auth([]libsecurity.HttpProvider{
		httpprovider.NewCookieProvider(s.GetTokenProviders(), s.cookieOptions.FileAccessName, s.logger),
	}, s.maintenanceAdapter, s.enforcer)
}

func (s *security) GetSessionStore() libsession.Store {
	return s.sessionStore
}

func (s *security) GetConfig() libsecurity.Config {
	return s.config
}

func (s *security) GetPasswordEncoder() password.Encoder {
	return password.NewBcryptEncoder()
}

func (s *security) GetTokenService() apisecurity.TokenService {
	return apisecurity.NewTokenService(s.config, s.dbClient, s.GetTokenGenerator(), token.NewMongoStore(s.dbClient, s.logger))
}

func (s *security) GetTokenProviders() []libsecurity.TokenProvider {
	return []libsecurity.TokenProvider{
		tokenprovider.NewTokenProvider(s.GetTokenGenerator(), token.NewMongoStore(s.dbClient, s.logger), s.newUserProvider(), s.logger),
		tokenprovider.NewTokenProvider(s.GetTokenGenerator(), sharetoken.NewMongoStore(s.dbClient, s.logger), s.newUserProvider(), s.logger),
	}
}

func (s *security) GetCookieOptions() CookieOptions {
	return s.cookieOptions
}

func (s *security) GetTokenGenerator() token.Generator {
	secretKey := os.Getenv(JwtSecretEnv)
	return token.NewJwtGenerator([]byte(secretKey), canopsis.AppName, s.apiConfigProvider)
}

func (s *security) newUserProvider() libsecurity.UserProvider {
	return userprovider.NewMongoProvider(s.dbClient, s.apiConfigProvider)
}

func (s *security) newBaseAuthProvider() libsecurity.Provider {
	return provider.NewBaseProvider(s.newUserProvider(), []password.Encoder{
		s.GetPasswordEncoder(),
		// todo deprecated encoder
		password.NewSha1Encoder(),
	})
}

func (s *security) newLdapAuthProvider() libsecurity.Provider {
	return provider.NewLdapProvider(
		s.config.Security.Ldap,
		s.newUserProvider(),
		provider.NewLdapDialer(),
		s.enforcer,
	)
}
