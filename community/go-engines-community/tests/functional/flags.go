package functional

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

const (
	dirMongoFixtures     = "testdata/fixtures/mongo"
	dirTimescaleFixtures = "testdata/fixtures/timescale"
)

type Flags struct {
	paths               arrayFlag
	mongoFixtures       arrayFlag
	timescaleFixtures   arrayFlag
	periodicalWaitTime  time.Duration
	dummyHttpPort       int64
	eventWaitKey        string
	eventWaitExchange   string
	eventsLog           string
	requestsLog         string
	checkUncaughtEvents bool
	onlyFixtures        bool
	randomize           int64
	concurrency         int
	tags                string
	clearOnScenario     bool
}

type arrayFlag []string

func (f *arrayFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *arrayFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *Flags) ParseArgs() {
	flag.Var(&f.paths, "paths", "All feature file paths.")
	flag.Var(&f.mongoFixtures, "mongoFixtures", "Mongo fixtures dirs.")
	flag.Var(&f.timescaleFixtures, "timescaleFixtures", "TimescaleDB fixtures dirs.")
	flag.DurationVar(&f.periodicalWaitTime, "pwt", 2200*time.Millisecond, "Duration to wait the end of next periodical process of all engines.")
	flag.StringVar(&f.eventWaitExchange, "ewe", "amq.direct", "Consume from exchange to detect the end of event processing.")
	flag.StringVar(&f.eventWaitKey, "ewk", canopsis.FIFOAckQueueName, "Consume by routing key to detect the end of event processing.")
	flag.StringVar(&f.eventsLog, "eventslog", "", "Log all received events.")
	flag.StringVar(&f.requestsLog, "requestslog", "", "Log all called API requests.")
	flag.Int64Var(&f.dummyHttpPort, "dummyHttpPort", 3000, "Port for dummy http server.")
	flag.BoolVar(&f.checkUncaughtEvents, "checkUncaughtEvents", false, "Enable catching event after each scenario.")
	flag.BoolVar(&f.onlyFixtures, "onlyFixtures", false, "Only apply fixtures.")
	flag.Int64Var(&f.randomize, "godog.randomize", 0, "Enable random order.")
	flag.IntVar(&f.concurrency, "godog.concurrency", 0, "Concurrency rate.")
	flag.StringVar(&f.tags, "godog.tags", "", "Filter scenarios.")
	flag.BoolVar(&f.clearOnScenario, "clearOnScenario", false, "Clear stores on each scenario.")
	flag.Parse()

	if !f.onlyFixtures && len(f.paths) == 0 {
		log.Fatal(fmt.Errorf("paths cannot be empty"))
	}

	paths := make([]string, 0, len(f.paths))
	for _, p := range f.paths {
		matches, err := filepath.Glob(p)
		if err == nil && matches != nil {
			paths = append(paths, matches...)
		} else {
			paths = append(paths, p)
		}
	}
	f.paths = paths

	if len(f.mongoFixtures) == 0 {
		f.mongoFixtures = []string{dirMongoFixtures}
	}

	if len(f.timescaleFixtures) == 0 {
		f.timescaleFixtures = []string{dirTimescaleFixtures}
	}
}
