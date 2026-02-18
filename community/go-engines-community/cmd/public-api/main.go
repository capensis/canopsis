package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	api "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/httpprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/tokenprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 30 * time.Second
)

const JwtSecretEnv = "CPS_JWT_SECRET" //nolint:gosec

func main() {
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var opts api.Options
	api.BindCmdFlags(&opts)
	flag.Parse()
	if opts.Version {
		canopsis.PrintVersionInfo()

		return
	}

	s := api.Services{}
	s.Logger = log.NewLogger(ctx, opts.Options)
	var err error
	s.PrimaryDbClient, err = mongo.NewClient(ctx)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot connect to mongodb")
	}

	defer func() {
		err = s.PrimaryDbClient.Disconnect(context.Background())
		if err != nil {
			s.Logger.Err(err).Msg("cannot disconnect from mongodb")
		}
	}()

	configAdapter := config.NewAdapter(s.PrimaryDbClient)
	s.Config, err = configAdapter.GetConfig(ctx)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot load config")
	}

	// Set mongodb setting.
	config.SetDbClientRetry(s.PrimaryDbClient, s.Config)
	// Connect to rmq.
	amqpConn, err := libamqp.NewConnection(s.Logger, -1, s.Config.Global.GetReconnectTimeout())
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot connect to rabbitmq")
	}

	defer func() {
		err = amqpConn.Close()
		if err != nil {
			s.Logger.Err(err).Msg("cannot disconnect from rabbitmq")
		}
	}()

	s.AMQPPublisher, err = amqpConn.Channel()
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot connect to rabbitmq")
	}

	defer func() {
		err = s.AMQPPublisher.Close()
		if err != nil {
			s.Logger.Err(err).Msg("cannot disconnect from rabbitmq")
		}
	}()

	s.SecurityConfig, err = libsecurity.LoadConfig(opts.ConfigDir)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot load security config")
	}

	s.APIConfigProvider = config.NewApiConfigProvider(s.Config, s.Logger)
	_, uniTrans, err := api.RegisterValidators(s.SecurityConfig)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot register validators")
	}

	s.Responder = response.NewResponder(json.NewEncoder(), validation.NewErrorTranslator(uniTrans, s.Logger), s.Logger)
	s.Enforcer, err = libsecurity.NewEnforcer(opts.ConfigDir, s.PrimaryDbClient)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("cannot create enforcer")
	}

	secretKey := os.Getenv(JwtSecretEnv)
	tokenGenerator := token.NewJwtGenerator([]byte(secretKey), canopsis.AppName, s.APIConfigProvider)
	s.AuthProviders = []libsecurity.HttpProvider{
		httpprovider.NewBearerProvider([]libsecurity.TokenProvider{
			tokenprovider.NewTokenProvider(tokenGenerator, token.NewMongoStore(s.PrimaryDbClient, s.Logger), userprovider.NewMongoProvider(s.PrimaryDbClient, s.APIConfigProvider), s.Logger),
		}),
	}
	r, err := api.RegisterRoutes(s, opts)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("failed to register routes")
	}

	hs := &http.Server{
		Addr:              ":" + strconv.Itoa(opts.Port),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	go func() {
		s.Logger.Info().Msg("http server started")
		err := hs.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Fatal().Err(err).Msg("failed to run http server")
		}
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()
	if err := hs.Shutdown(shutdownCtx); err != nil {
		s.Logger.Fatal().Err(err).Msg("failed to shutdown http server")

		return
	}

	s.Logger.Info().Msg("http server stopped")
}
