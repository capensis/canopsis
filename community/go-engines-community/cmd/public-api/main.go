package main

//go:generate go tool github.com/swaggo/swag/cmd/swag init -parseDependency -d ../../lib/api/public -g ../../../cmd/public-api/$GOFILE -o ../../lib/api/public/docs --outputTypes yaml --instanceName schemas

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	api "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
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
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 30 * time.Second
	jobWorkers        = 5
)

const JwtSecretEnv = "CPS_JWT_SECRET" //nolint:gosec

func main() {
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var opts api.Options
	api.BindCmdFlags(&opts)
	flag.Parse()
	if opts.Version {
		canopsis.PrintVersionInfo()

		return
	}

	logger := log.NewLogger(ctx, opts.Options)
	s, cleanup, err := newServices(ctx, opts, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create services")
	}

	defer cleanup()

	r, err := api.RegisterRoutes(s.Services, opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register routes")
	}

	wg := runBackground(ctx, s, cancel)
	hs := serve(r, opts.Port, s)

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()
	err = hs.Shutdown(shutdownCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to shutdown http server")
	}

	logger.Info().Msg("http server stopped")

	// Let the worker finish the running jobs before connections are closed.
	err = wg.Wait()
	if err != nil {
		logger.Fatal().Err(err).Msg("workers stopped with error")
	}
}

type services struct {
	api.Services

	configAdapter          config.Adapter
	apiConfigProvider      *config.BaseApiConfigProvider
	timezoneConfigProvider *config.BaseTimezoneConfigProvider
	consumeChPool          libamqp.ChannelPool
}

func newServices(ctx context.Context, opts api.Options, logger zerolog.Logger) (s services, cleanup func(), err error) {
	s.Logger = logger
	var cleanups []func() error
	cleanup = func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			err := cleanups[i]()
			if err != nil {
				logger.Err(err).Msg("cannot close connection")
			}
		}
	}

	s.PrimaryDbClient, err = mongo.NewClient(ctx)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot connect to mongodb: %w", err)
	}

	cleanups = append(cleanups, func() error {
		return s.PrimaryDbClient.Disconnect(context.WithoutCancel(ctx))
	})

	s.configAdapter = config.NewAdapter(s.PrimaryDbClient)
	s.Config, err = s.configAdapter.GetConfig(ctx)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot load config: %w", err)
	}

	config.SetDbClientRetry(s.PrimaryDbClient, s.Config)

	s.SecondaryDbClient, err = mongo.NewClient(ctx, mongo.ClientOptions{
		RetryCount:      s.Config.Global.ReconnectRetries,
		MinRetryTimeout: s.Config.Global.GetReconnectTimeout(),
		ReadPreference:  mongo.SecondaryPreferred(),
	})
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot connect to mongodb: %w", err)
	}

	cleanups = append(cleanups, func() error {
		return s.SecondaryDbClient.Disconnect(context.WithoutCancel(ctx))
	})

	amqpPubConn, err := libamqp.New(-1, s.Config.Global.GetReconnectTimeout(), logger)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot connect to rabbitmq: %w", err)
	}

	cleanups = append(cleanups, amqpPubConn.Close)

	amqpConsumeConn, err := libamqp.New(-1, s.Config.Global.GetReconnectTimeout(), logger)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot connect to rabbitmq: %w", err)
	}

	cleanups = append(cleanups, amqpConsumeConn.Close)

	s.AMQPPubChPool = libamqp.NewChannelPool(amqpPubConn, canopsis.DefaultAMQPPublishPoolSize)
	cleanups = append(cleanups, s.AMQPPubChPool.Close)

	s.consumeChPool = libamqp.NewChannelPool(amqpConsumeConn, canopsis.DefaultAMQPPublishPoolSize)
	cleanups = append(cleanups, s.consumeChPool.Close)

	s.SecurityConfig, err = libsecurity.LoadConfig(opts.ConfigDir, logger)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot load security config: %w", err)
	}

	s.apiConfigProvider = config.NewApiConfigProvider(s.Config, logger)
	s.timezoneConfigProvider = config.NewTimezoneConfigProvider(s.Config, logger)
	s.APIConfigProvider = s.apiConfigProvider
	s.TimezoneConfigProvider = s.timezoneConfigProvider

	validator, uniTrans, err := api.RegisterValidators()
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot register validators: %w", err)
	}

	s.Binder = binding.NewBinder(json.NewDecoder(), validator)
	s.Responder = response.NewResponder(json.NewEncoder(), validation.NewErrorTranslator(uniTrans, logger), logger)
	s.Enforcer, err = libsecurity.NewEnforcer(opts.ConfigDir, s.PrimaryDbClient)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot create enforcer: %w", err)
	}

	tokenGenerator := token.NewJwtGenerator([]byte(os.Getenv(JwtSecretEnv)), canopsis.AppName, s.APIConfigProvider)
	s.AuthProviders = []libsecurity.HttpProvider{
		httpprovider.NewBearerProvider([]libsecurity.TokenProvider{
			tokenprovider.NewTokenProvider(tokenGenerator, token.NewMongoStore(s.PrimaryDbClient, logger), userprovider.NewMongoProvider(s.PrimaryDbClient, s.APIConfigProvider), logger),
		}),
	}

	s.JobPublisher = asyncjob.NewPublisher(libamqp.NewPooledPublisher(s.AMQPPubChPool), json.NewEncoder())

	return s, cleanup, nil
}

// serve starts the http server and returns it, so it can be shut down gracefully.
func serve(r http.Handler, port int, s services) *http.Server {
	hs := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
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

	return hs
}

// runBackground starts background works and returns the group to wait on once the http server stopped.
// cancel stops the whole api when the job worker fails.
func runBackground(ctx context.Context, s services, cancel context.CancelFunc) *errgroup.Group {
	jobHandlers, jobMaintainers := api.AsyncJobs(s.Services)
	jobWorker := asyncjob.NewWorker(
		s.consumeChPool,
		json.NewDecoder(),
		jobHandlers,
		jobWorkers,
		s.Logger,
	)

	wg, wctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		err := jobWorker.Run(wctx)
		if err != nil {
			cancel()

			return fmt.Errorf("async job worker stopped with error: %w", err)
		}

		return nil
	})

	wg.Go(func() error {
		asyncjob.RunMaintenance(wctx, jobMaintainers, s.Logger)

		return nil
	})

	wg.Go(func() error {
		config.RunReload(wctx, s.configAdapter, canopsis.PeriodicalWaitTime, s.Logger, s.apiConfigProvider, s.timezoneConfigProvider)

		return nil
	})

	return wg
}
