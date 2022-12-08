package functional

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/bdd"
	libjson "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	liblog "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/cucumber/godog"
	redismod "github.com/go-redis/redis/v8"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/rs/zerolog"
)

const (
	websocketScheme = "ws"
	websocketRoute  = "/api/v4/ws"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := liblog.NewLogger(true)

	apiUrl, err := bdd.GetApiURL()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	var flags Flags
	flags.ParseArgs()

	var eventLogger zerolog.Logger
	if flags.eventLogs != "" {
		f, err := os.OpenFile(flags.eventLogs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		defer f.Close()
		eventLogger = zerolog.New(&eventLogWriter{writer: f}).
			Level(zerolog.DebugLevel).
			With().Timestamp().
			Logger()
	}

	err = bdd.RunDummyHttpServer(ctx, fmt.Sprintf("localhost:%d", flags.dummyHttpPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer func() {
		if err = dbClient.Disconnect(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}()

	amqpConnection, err := amqp.NewConnection(logger, 0, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer amqpConnection.Close()

	redisClient, err := redis.NewSession(ctx, 0, logger, 0, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer redisClient.Close()

	loader := fixtures.NewLoader(dbClient, flags.mongoFixtures,
		fixtures.NewParser(fixtures.NewFaker(password.NewSha1Encoder())), logger)
	opts := godog.Options{
		StopOnFailure:  true,
		Format:         "pretty",
		Paths:          flags.paths,
		DefaultContext: ctx,
		Randomize:      flags.randomize,
		Concurrency:    flags.concurrency,
		Tags:           flags.tags,
	}

	if flags.onlyFixtures {
		err := clearStores(ctx, flags, loader, redisClient, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}

		return
	}

	templater := bdd.NewTemplater(map[string]interface{}{"apiURL": apiUrl})
	apiClient := bdd.NewApiClient(dbClient, apiUrl, templater)
	amqpClient := bdd.NewAmqpClient(dbClient, amqpConnection, flags.eventWaitExchange, flags.eventWaitKey,
		libjson.NewEncoder(), libjson.NewDecoder(), eventLogger, templater)
	mongoClient := bdd.NewMongoClient(dbClient)
	wsUrl, err := url.Parse(apiUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("invalid api url")
	}
	wsUrl.Scheme = websocketScheme
	wsUrl.Path = websocketRoute
	websocketClient := bdd.NewWebsocketClient(wsUrl.String(), templater)

	testSuiteInitializer := InitializeTestSuite(ctx, flags, loader, redisClient, logger)
	scenarioInitializer := InitializeScenario(flags, apiClient, amqpClient, mongoClient, websocketClient, loader, redisClient, logger)
	status := godog.TestSuite{
		Name:                 "canopsis",
		TestSuiteInitializer: testSuiteInitializer,
		ScenarioInitializer:  scenarioInitializer,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func InitializeTestSuite(
	ctx context.Context,
	flags Flags,
	loader fixtures.Loader,
	redisClient redismod.Cmdable,
	logger zerolog.Logger,
) func(*godog.TestSuiteContext) {
	return func(godogCtx *godog.TestSuiteContext) {
		if !flags.clearOnScenario {
			godogCtx.BeforeSuite(func() {
				err := clearStores(ctx, flags, loader, redisClient, logger)
				if err != nil {
					logger.Fatal().Err(err).Msg("")
				}

				logger.Info().Msg("waiting the next periodical process")
				time.Sleep(flags.periodicalWaitTime)
			})
		}
		godogCtx.AfterSuite(func() {
			err := clearStores(ctx, flags, loader, redisClient, logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("")
			}
		})
	}
}

func InitializeScenario(
	flags Flags,
	apiClient *bdd.ApiClient,
	amqpClient *bdd.AmqpClient,
	mongoClient *bdd.MongoClient,
	websocketClient *bdd.WebsocketClient,
	loader fixtures.Loader,
	redisClient redismod.Cmdable,
	logger zerolog.Logger,
) func(*godog.ScenarioContext) {
	return func(scenarioCtx *godog.ScenarioContext) {
		if flags.checkUncaughtEvents {
			scenarioCtx.After(func(ctx context.Context, sc *godog.Scenario, scErr error) (context.Context, error) {
				if scErr == nil {
					err := amqpClient.IWaitTheEndOfEventProcessing(ctx)
					if err == nil {
						return ctx, errors.New("caught event")
					}
				}

				return ctx, scErr
			})
		}

		scenarioCtx.Before(amqpClient.BeforeScenario)
		scenarioCtx.After(amqpClient.AfterScenario)
		scenarioCtx.After(websocketClient.AfterScenario)

		if flags.clearOnScenario {
			scenarioCtx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
				err := clearStores(ctx, flags, loader, redisClient, logger)
				if err != nil {
					return ctx, err
				}
				logger.Info().Msg("waiting the next periodical process")
				time.Sleep(flags.periodicalWaitTime)
				return ctx, nil
			})
		}

		scenarioCtx.Step(`^I am ([\w-]+)$`, apiClient.IAm)
		scenarioCtx.Step(`^I am authenticated with username "([^"]+)" and password "([^"]+)"$`, apiClient.IAmAuthenticatedByBasicAuth)
		scenarioCtx.Step(`^I send an event:$`, apiClient.ISendAnEvent)
		scenarioCtx.Step(`^I do (\w+) (.+) until response code is (\d+) and body is:$`, apiClient.IDoRequestUntilResponse)
		scenarioCtx.Step(`^I do (\w+) (.+) until response code is (\d+) and body contains:$`, apiClient.IDoRequestUntilResponseContains)
		scenarioCtx.Step(`^I do (\w+) (.+) until response code is (\d+) and response key \"([\w\.]+)\" is greater or equal than (\d+)$`, apiClient.IDoRequestUntilResponseKeyIsGreaterOrEqualThan)
		scenarioCtx.Step(`^I do (\w+) (.+) until response code is (\d+) and response array key \"([\w\.]+)\" contains:$`, apiClient.IDoRequestUntilResponseArrayKeyContains)
		scenarioCtx.Step(`^I do (\w+) (.+):$`, apiClient.IDoRequestWithBody)
		scenarioCtx.Step(`^I do (\w+) (.+) until response code is (\d+)$`, apiClient.IDoRequestUntilResponseCode)
		scenarioCtx.Step(`^I do (\w+) (.+)$`, apiClient.IDoRequest)
		scenarioCtx.Step(`^the response code should be (\d+)$`, apiClient.TheResponseCodeShouldBe)
		scenarioCtx.Step(`^the response body should be:$`, apiClient.TheResponseBodyShouldBe)
		scenarioCtx.Step(`^the response body should contain:$`, apiClient.TheResponseBodyShouldContain)
		scenarioCtx.Step(`^the response raw body should be:$`, apiClient.TheResponseRawBodyShouldBe)
		scenarioCtx.Step(`^the response key \"([\w\.]+)\" should exist$`, apiClient.TheResponseKeyShouldExist)
		scenarioCtx.Step(`^the response key \"([\w\.]+)\" should not exist$`, apiClient.TheResponseKeyShouldNotExist)
		scenarioCtx.Step(`^the response key \"([\w\.]+)\" should not be \"([^\"]+)\"$`, apiClient.TheResponseKeyShouldNotBe)
		scenarioCtx.Step(`^the difference between ([\w\.]+) ([\w\.]+) is in range (-?\d+\.?\d*),(-?\d+\.?\d*)$`, apiClient.TheDifferenceBetweenValues)
		scenarioCtx.Step(`^the response key \"([\w\.]+)\" should be greater or equal than (\d+)$`, apiClient.TheResponseKeyShouldBeGreaterOrEqualThan)
		scenarioCtx.Step(`^the response key \"([\w\.]+)\" should be greater than (\d+)$`, apiClient.TheResponseKeyShouldBeGreaterThan)
		scenarioCtx.Step(`^the response array key \"([\w\.]+)\" should contain:$`, apiClient.TheResponseArrayKeyShouldContain)
		scenarioCtx.Step(`^the response array key \"([\w\.]+)\" should contain only:$`, apiClient.TheResponseArrayKeyShouldContainOnly)
		scenarioCtx.Step(`^I save response ([\w]+)=(.+)$`, apiClient.ISaveResponse)
		scenarioCtx.Step(`^\"([\w]+)\" (>|<|>=|<=) \"([\w]+)\"$`, apiClient.ValueShouldBeGteLteThan)
		scenarioCtx.Step(`^an alarm (.+) should be in the db$`, mongoClient.AlarmShouldBeInTheDb)
		scenarioCtx.Step(`^an entity (.+) should be in the db$`, mongoClient.EntityShouldBeInTheDb)
		scenarioCtx.Step(`^an entity (.+) should not be in the db$`, mongoClient.EntityShouldNotBeInTheDb)
		scenarioCtx.Step(`^I set header ([\w\.\-]+)=(.+)$`, apiClient.ISetRequestHeader)
		scenarioCtx.Step(`^I wait (\w+)$`, func(str string) error {
			duration, err := time.ParseDuration(str)
			if err != nil {
				return err
			}
			time.Sleep(duration)
			return nil
		})
		scenarioCtx.Step(`^I wait the next periodical process$`, func() error {
			time.Sleep(flags.periodicalWaitTime)
			return nil
		})
		scenarioCtx.Step(`^I wait the end of event processing$`, amqpClient.IWaitTheEndOfEventProcessing)
		scenarioCtx.Step(`^I wait the end of (\d+) events processing$`, amqpClient.IWaitTheEndOfEventsProcessing)
		scenarioCtx.Step(`^I wait the end of event processing which contains:$`, amqpClient.IWaitTheEndOfEventProcessingWhichContains)
		scenarioCtx.Step(`^I wait the end of events processing which contain:$`, amqpClient.IWaitTheEndOfEventsProcessingWhichContain)
		scenarioCtx.Step(`^I wait the end of one of events processing which contain:$`, amqpClient.IWaitTheEndOfOneOfEventsProcessingWhichContain)
		scenarioCtx.Step(`^I send an event and wait the end of event processing:$`, func(ctx context.Context, doc string) (context.Context, error) {
			ctx, err := apiClient.ISendAnEvent(ctx, doc)
			if err != nil {
				return ctx, err
			}

			return ctx, amqpClient.IWaitTheEndOfSentEventProcessing(ctx, doc)
		})
		scenarioCtx.Step(`^I call RPC to engine-axe with alarm ([^:]+):$`, amqpClient.ICallRPCAxeRequest)
		scenarioCtx.Step(`^I call RPC to engine-webhook with alarm ([^:]+):$`, amqpClient.ICallRPCWebhookRequest)
		scenarioCtx.Step(`^I subscribe to websocket room \"([^\"]+)\"$`, websocketClient.ISubscribeToRoom)
		scenarioCtx.Step(`^I wait message from websocket room \"([^\"]+)\":$`, websocketClient.IWaitMessageFromRoom)
		scenarioCtx.Step(`^I wait message from websocket room \"([^\"]+)\" which contains:$`, websocketClient.IWaitMessageFromRoomWhichContains)
	}
}

func clearStores(
	ctx context.Context,
	flags Flags,
	loader fixtures.Loader,
	redisClient redismod.Cmdable,
	logger zerolog.Logger,
) error {
	err := loader.Load(ctx)
	if err != nil {
		return fmt.Errorf("cannot load mongo fixtures: %w", err)
	}

	logger.Info().Msg("MongoDB fixtures are applied")
	pgConnStr, err := postgres.GetConnStr()
	if err != nil {
		return err
	}

	pgDb, err := sql.Open("pgx", pgConnStr)
	if err != nil {
		return fmt.Errorf("cannot connect to timescale for fixtures: %w", err)
	}
	defer pgDb.Close()

	tsFixtures, err := testfixtures.New(
		testfixtures.Database(pgDb),
		testfixtures.Dialect("timescaledb"),
		testfixtures.UseAlterConstraint(),
		testfixtures.Paths(flags.timescaleFixtures...),
	)
	if err != nil {
		return fmt.Errorf("cannot init timescale fixtures: %w", err)
	}

	err = tsFixtures.Load()
	if err != nil {
		return fmt.Errorf("cannot load timescale fixtures: %w", err)
	}

	logger.Info().Msg("PostgresSQL fixtures are applied")

	err = redisClient.FlushAll(ctx).Err()
	if err != nil {
		return err
	}

	logger.Info().Msg("Redis is flushed")

	return nil
}

type eventLogWriter struct {
	writer io.Writer
}

func (w *eventLogWriter) Write(p []byte) (int, error) {
	var msg map[string]interface{}
	err := json.Unmarshal(p, &msg)
	if err != nil {
		return 0, err
	}

	fieldsStr := ""
	for k, v := range msg {
		switch k {
		case zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.MessageFieldName:
		default:
			s, err := json.Marshal(v)
			if err != nil {
				return 0, err
			}
			fieldsStr += fmt.Sprintf("%s=%s ", k, s)
		}
	}

	formattedMsg := fmt.Sprintf("%s %s > %s %s\n", msg[zerolog.TimestampFieldName],
		msg[zerolog.LevelFieldName], msg[zerolog.MessageFieldName], fieldsStr)

	return w.writer.Write([]byte(formattedMsg))
}
