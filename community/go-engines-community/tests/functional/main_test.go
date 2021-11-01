package functional

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/bdd"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libjson "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	liblog "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/cucumber/godog"
	redismod "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Flags struct {
	paths               arrayFlag
	fixtures            arrayFlag
	periodicalWaitTime  time.Duration
	dummyHttpPort       int64
	eventWaitKey        string
	eventWaitExchange   string
	eventLogs           string
	checkUncaughtEvents bool
}

type arrayFlag []string

func (f *arrayFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *arrayFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test allowed only with "API_URL" environment variable
	if _, err := bdd.GetApiURL(); err != nil {
		os.Exit(0)
	}

	var flags Flags
	flag.Var(&flags.paths, "paths", "All feature file paths.")
	flag.Var(&flags.fixtures, "fixtures", "All fixtures dirs.")
	flag.DurationVar(&flags.periodicalWaitTime, "pwt", 2200*time.Millisecond, "Duration to wait the end of next periodical process of all engines.")
	flag.StringVar(&flags.eventWaitExchange, "ewe", "amq.direct", "Consume from exchange to detect the end of event processing.")
	flag.StringVar(&flags.eventWaitKey, "ewk", canopsis.FIFOAckQueueName, "Consume by routing key to detect the end of event processing.")
	flag.StringVar(&flags.eventLogs, "eventlogs", "", "Log all received events.")
	flag.Int64Var(&flags.dummyHttpPort, "dummyHttpPort", 3000, "Port for dummy http server.")
	flag.BoolVar(&flags.checkUncaughtEvents, "checkUncaughtEvents", false, "Enable catching event after each scenario.")
	flag.Parse()

	if len(flags.paths) == 0 {
		log.Fatal(fmt.Errorf("paths cannot be empty"))
	}

	if len(flags.fixtures) == 0 {
		flags.fixtures = []string{"../../database/fixtures/test"}
	}

	logger := liblog.NewLogger(true)
	var eventLogger zerolog.Logger
	if flags.eventLogs != "" {
		f, err := os.OpenFile(flags.eventLogs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		eventLogger = zerolog.New(&eventLogWriter{writer: f}).
			Level(zerolog.DebugLevel).
			With().Timestamp().
			Logger()
	}

	paths := make([]string, 0, len(flags.paths))
	for _, p := range flags.paths {
		matches, err := filepath.Glob(p)
		if err == nil && matches != nil {
			paths = append(paths, matches...)
		} else {
			paths = append(paths, p)
		}
	}
	err := bdd.RunDummyHttpServer(ctx, fmt.Sprintf("localhost:%d", flags.dummyHttpPort))
	if err != nil {
		log.Fatal(err)
	}

	dbClient, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = dbClient.Disconnect(context.Background())
		if err != nil {
			log.Print(err)
		}
	}()

	amqpConnection, err := amqp.NewConnection(logger, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = amqpConnection.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	redisClient, err := redis.NewSession(ctx, 0, logger, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = redisClient.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	loader := fixtures.NewLoader(dbClient, flags.fixtures, true,
		fixtures.NewParser(password.NewSha1Encoder()), logger)

	opts := godog.Options{
		StopOnFailure:  true,
		Format:         "pretty",
		Paths:          paths,
		DefaultContext: ctx,
	}
	testSuiteInitializer := InitializeTestSuite(ctx, flags, loader, redisClient)
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

func InitializeTestSuite(ctx context.Context, flags Flags, loader fixtures.Loader, redisClient redismod.Cmdable) func(*godog.TestSuiteContext) {
	return func(godogCtx *godog.TestSuiteContext) {
		godogCtx.BeforeSuite(func() {
			err := clearStores(ctx, loader, redisClient)
			if err != nil {
				panic(err)
			}
			time.Sleep(flags.periodicalWaitTime)
		})
		godogCtx.AfterSuite(func() {
			err := clearStores(ctx, loader, redisClient)
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
		ctx.Step(`^I do (\w+) ([^:]+) until response code is (\d+) and body is:$`, apiClient.IDoRequestUntilResponse)
		ctx.Step(`^I do (\w+) ([^:]+) until response code is (\d+) and body contains:$`, apiClient.IDoRequestUntilResponseContains)
		ctx.Step(`^I do (\w+) ([^:]+):$`, apiClient.IDoRequestWithBody)
		ctx.Step(`^I do (\w+) (.+) until response code is (\d+)$`, apiClient.IDoRequestUntilResponseCode)
		ctx.Step(`^I do (\w+) (.+)$`, apiClient.IDoRequest)
		ctx.Step(`^the response code should be (\d+)$`, apiClient.TheResponseCodeShouldBe)
		ctx.Step(`^the response body should be:$`, apiClient.TheResponseBodyShouldBe)
		ctx.Step(`^the response body should contain:$`, apiClient.TheResponseBodyShouldContain)
		ctx.Step(`^the response raw body should be:$`, apiClient.TheResponseRawBodyShouldBe)
		ctx.Step(`^the response key \"([\w\.]+)\" should not exist$`, apiClient.TheResponseKeyShouldNotExist)
		ctx.Step(`^the response key \"([\w\.]+)\" should not be \"([^\"]+)\"$`, apiClient.TheResponseKeyShouldNotBe)
		ctx.Step(`^I save response ([\w]+)=(.+)$`, apiClient.ISaveResponse)
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

func clearStores(ctx context.Context, loader fixtures.Loader, redisClient redismod.Cmdable) error {
	err := loader.Load(ctx)
	if err != nil {
		return err
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
