package functional

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	"github.com/cucumber/godog"
	redismod "github.com/go-redis/redis/v8"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test allowed only with "API_URL" environment variable
	if _, err := bdd.GetApiURL(); err != nil {
		log.Fatal(err)
	}

	var flags Flags
	flags.ParseArgs()

	var eventLogger zerolog.Logger
	if flags.eventLogs != "" {
		f, err := os.OpenFile(flags.eventLogs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		eventLogger = zerolog.New(&eventLogWriter{writer: f}).
			Level(zerolog.DebugLevel).
			With().Timestamp().
			Logger()
	}

	err := bdd.RunDummyHttpServer(ctx, fmt.Sprintf("localhost:%d", flags.dummyHttpPort))
	if err != nil {
		log.Fatal(err)
	}

	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(context.Background())

	amqpConnection, err := amqp.NewConnection(liblog.NewLogger(false), 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()

	redisClient, err := redis.NewSession(ctx, 0, liblog.NewLogger(false), 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	opts := godog.Options{
		StopOnFailure:  true,
		Format:         "pretty",
		Paths:          flags.paths,
		DefaultContext: ctx,
	}
	testSuiteInitializer := InitializeTestSuite(ctx, flags, dbClient, redisClient)
	scenarioInitializer, err := InitializeScenario(flags, dbClient, amqpConnection, eventLogger)
	if err != nil {
		log.Fatal(err)
	}

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

func InitializeTestSuite(ctx context.Context, flags Flags, dbClient mongo.DbClient, redisClient redismod.Cmdable) func(*godog.TestSuiteContext) {
	return func(godogCtx *godog.TestSuiteContext) {
		godogCtx.BeforeSuite(func() {
			err := clearStores(ctx, flags, dbClient, redisClient)
			if err != nil {
				panic(err)
			}
			time.Sleep(flags.periodicalWaitTime)
		})
		godogCtx.AfterSuite(func() {
			err := clearStores(ctx, flags, dbClient, redisClient)
			if err != nil {
				panic(err)
			}
		})
	}
}

func InitializeScenario(flags Flags, dbClient mongo.DbClient, amqpConnection amqp.Connection,
	eventLogger zerolog.Logger) (func(*godog.ScenarioContext), error) {
	apiClient, err := bdd.NewApiClient(dbClient)
	if err != nil {
		return nil, err
	}

	mongoClient, err := bdd.NewMongoClient(dbClient)
	if err != nil {
		return nil, err
	}

	amqpClient, err := bdd.NewAmqpClient(dbClient, amqpConnection,
		flags.eventWaitExchange, flags.eventWaitKey,
		libjson.NewEncoder(), libjson.NewDecoder(), eventLogger)
	if err != nil {
		return nil, err
	}

	return func(ctx *godog.ScenarioContext) {
		ctx.Before(apiClient.ResetResponse)
		ctx.Before(amqpClient.Reset)
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			eventLogger.Info().Str("file", sc.Uri).Msgf("%s", sc.Name)
			return ctx, nil
		})

		if flags.checkUncaughtEvents {
			ctx.After(func(ctx context.Context, sc *godog.Scenario, scErr error) (context.Context, error) {
				if scErr == nil {
					err := amqpClient.IWaitTheEndOfEventProcessing()
					if err == nil {
						return ctx, errors.New("caught event")
					}
				}

				return ctx, scErr
			})
		}

		ctx.Step(`^I am ([\w-]+)$`, apiClient.IAm)
		ctx.Step(`^I am authenticated with username "([^"]+)" and password "([^"]+)"$`, apiClient.IAmAuthenticatedByBasicAuth)
		ctx.Step(`^I send an event:$`, apiClient.ISendAnEvent)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+) and body is:$`, apiClient.IDoRequestUntilResponse)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+) and body contains:$`, apiClient.IDoRequestUntilResponseContains)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+) and response key \"([\w\.]+)\" is greater or equal than (\d+)$`, apiClient.IDoRequestUntilResponseKeyIsGreaterOrEqualThan)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+) and response array key \"([\w\.]+)\" contains:$`, apiClient.IDoRequestUntilResponseArrayKeyContains)
		ctx.Step(`^I do (\w+) (.+):$`, apiClient.IDoRequestWithBody)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+)$`, apiClient.IDoRequestUntilResponseCode)
		ctx.Step(`^I do (\w+) (.+)$`, apiClient.IDoRequest)
		ctx.Step(`^the response code should be (\d+)$`, apiClient.TheResponseCodeShouldBe)
		ctx.Step(`^the response body should be:$`, apiClient.TheResponseBodyShouldBe)
		ctx.Step(`^the response body should contain:$`, apiClient.TheResponseBodyShouldContain)
		ctx.Step(`^the response raw body should be:$`, apiClient.TheResponseRawBodyShouldBe)
		ctx.Step(`^the response key \"([\w\.]+)\" should not exist$`, apiClient.TheResponseKeyShouldNotExist)
		ctx.Step(`^the response key \"([\w\.]+)\" should not be \"([^\"]+)\"$`, apiClient.TheResponseKeyShouldNotBe)
		ctx.Step(`^the response key \"([\w\.]+)\" should be greater or equal than (\d+)$`, apiClient.TheResponseKeyShouldBeGreaterOrEqualThan)
		ctx.Step(`^the response array key \"([\w\.]+)\" should contain:$`, apiClient.TheResponseArrayKeyShouldContain)
		ctx.Step(`^the response array key \"([\w\.]+)\" should contain only one:$`, apiClient.TheResponseArrayKeyShouldContainOnlyOne)
		ctx.Step(`^the response array key \"([\w\.]+)\" should contain only:$`, apiClient.TheResponseArrayKeyShouldContainOnly)
		ctx.Step(`^I save response ([\w]+)=(.+)$`, apiClient.ISaveResponse)
		ctx.Step(`^\"([\w]+)\" (>|<|>=|<=) \"([\w]+)\"$`, apiClient.ValueShouldBeGteLteThan)
		ctx.Step(`^an alarm (.+) should be in the db$`, mongoClient.AlarmShouldBeInTheDb)
		ctx.Step(`^an entity (.+) should be in the db$`, mongoClient.EntityShouldBeInTheDb)
		ctx.Step(`^I set header ([\w\.\-]+)=(.+)$`, apiClient.ISetRequestHeader)
		ctx.Step(`^I wait (\w+)$`, func(str string) error {
			duration, err := time.ParseDuration(str)
			if err != nil {
				return err
			}
			time.Sleep(duration)
			return nil
		})
		ctx.Step(`^I wait the next periodical process$`, func() error {
			time.Sleep(flags.periodicalWaitTime)
			return nil
		})
		ctx.Step(`^I wait the end of event processing$`, amqpClient.IWaitTheEndOfEventProcessing)
		ctx.Step(`^I wait the end of (\d+) events processing$`, amqpClient.IWaitTheEndOfEventsProcessing)
		ctx.Step(`^I call RPC to engine-axe with alarm ([^:]+):$`, amqpClient.ICallRPCAxeRequest)
		ctx.Step(`^I call RPC to engine-webhook with alarm ([^:]+):$`, amqpClient.ICallRPCWebhookRequest)
	}, nil
}

func clearStores(ctx context.Context, flags Flags, dbClient mongo.DbClient, redisClient redismod.Cmdable) error {
	err := fixtures.Load(ctx, dbClient, flags.mongoFixtures)
	if err != nil {
		return fmt.Errorf("cannot load mongo fixtures: %w", err)
	}

	pgConnStr, err := postgres.GetConnStr()
	if err != nil {
		return err
	}

	p := &pgx.Postgres{}
	driver, err := p.Open(pgConnStr)
	if err != nil {
		return fmt.Errorf("cannot connect to timescale for migrations: %w", err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", flags.timescaleMigrations), "pgx", driver)
	if err != nil {
		return fmt.Errorf("cannot init timescale migrations: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cannot apply timescale migrations: %w", err)
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

	err = redisClient.FlushAll(ctx).Err()
	if err != nil {
		return err
	}

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
