package api

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	libtoken "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/httpprovider"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/rs/zerolog"
)

type Services struct {
	Logger                 zerolog.Logger
	Binder                 binding.Binder
	Responder              response.Responder
	PrimaryDbClient        mongo.DbClient
	SecondaryDbClient      mongo.DbClient
	AMQPPubChPool          libamqp.ChannelPool
	AMQPConsumeChPool      libamqp.ChannelPool
	Config                 config.CanopsisConf
	SecurityConfig         libsecurity.Config
	Enforcer               libsecurity.Enforcer
	AuthProviders          []libsecurity.HttpProvider
	APIConfigProvider      config.ApiConfigProvider
	TimezoneConfigProvider config.TimezoneConfigProvider
	TemplateConfigProvider config.TemplateConfigProvider
	JobPublisher           asyncjob.Publisher
	UserMetaUpdater        metrics.MetaUpdater
	LinkGenerator          link.Generator
	UserProvider           liblink.UserProvider

	configAdapter          config.Adapter
	apiConfigProvider      *config.BaseApiConfigProvider
	timezoneConfigProvider *config.BaseTimezoneConfigProvider
	templateConfigProvider *config.BaseTemplateConfigProvider
}

func NewServices(ctx context.Context, opts Options, logger zerolog.Logger) (s Services, cleanup func(), err error) {
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

	s.AMQPConsumeChPool = libamqp.NewChannelPool(amqpConsumeConn, canopsis.DefaultAMQPPublishPoolSize)
	cleanups = append(cleanups, s.AMQPConsumeChPool.Close)

	s.SecurityConfig, err = libsecurity.LoadConfig(opts.ConfigDir, logger)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot load security config: %w", err)
	}

	s.apiConfigProvider = config.NewApiConfigProvider(s.Config, logger)
	s.timezoneConfigProvider = config.NewTimezoneConfigProvider(s.Config, logger)
	s.templateConfigProvider = config.NewTemplateConfigProvider(s.Config, logger)
	s.APIConfigProvider = s.apiConfigProvider
	s.TimezoneConfigProvider = s.timezoneConfigProvider
	s.TemplateConfigProvider = s.templateConfigProvider

	validator, uniTrans, err := RegisterValidators(s.SecurityConfig)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot register validators: %w", err)
	}

	s.Binder = binding.NewBinder(json.NewDecoder(), validator)
	s.Responder = response.NewResponder(json.NewEncoder(), validation.NewErrorTranslator(uniTrans, logger), logger)
	s.Enforcer, err = libsecurity.NewEnforcer(opts.ConfigDir, s.PrimaryDbClient)
	if err != nil {
		return s, cleanup, fmt.Errorf("cannot create enforcer: %w", err)
	}

	s.AuthProviders = []libsecurity.HttpProvider{
		httpprovider.NewBearerProvider([]libsecurity.TokenProvider{
			libtoken.NewProvider(s.PrimaryDbClient, userprovider.NewMongoProvider(s.PrimaryDbClient, s.APIConfigProvider)),
		}),
	}

	s.JobPublisher = asyncjob.NewPublisher(libamqp.NewPooledPublisher(s.AMQPPubChPool), json.NewEncoder())

	s.UserMetaUpdater = metrics.NewNullMetaUpdater()

	tplExecutor := template.NewExecutor(s.TemplateConfigProvider, s.TimezoneConfigProvider)
	s.LinkGenerator = link.NewGenerator(s.SecondaryDbClient, externaldata.NewNullGetter(), tplExecutor, s.Logger)
	s.UserProvider = liblink.NewUserProvider(s.SecondaryDbClient)

	return s, cleanup, nil
}
