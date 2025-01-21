package che

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
)

type Options struct {
	Version                 bool
	FeatureEventProcessing  bool
	FeatureContextCreation  bool
	Purge                   bool
	PrintEventOnError       bool
	ModeDebug               bool
	PeriodicalWaitTime      time.Duration
	InfosDictionaryWaitTime time.Duration
	ExternalDataApiTimeout  time.Duration
	SoftDeleteWaitTime      time.Duration
	CleanPerfDataWaitTime   time.Duration
	FifoAckExchange         string
	ExternalWorkers         int
	SystemWorkers           int
	UserWorkers             int
}

func ParseOptions() (Options, []string) {
	opts := Options{}
	flag.BoolVar(&opts.FeatureEventProcessing, "processEvent", true, "enable event processing. enabled by default.")
	flag.BoolVar(&opts.FeatureContextCreation, "createContext", true, "enable context graph creation. enabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.Purge, "purge", false, "purge consumer queue(s) before work")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two runs of periodical process")
	flag.DurationVar(&opts.InfosDictionaryWaitTime, "infosDictionaryWaitTime", time.Hour, "Duration to wait between two runs of update entity infos dictionary process")
	flag.DurationVar(&opts.SoftDeleteWaitTime, "softDeleteWaitTime", time.Hour, "Duration to keep soft deleted entities in the db until they will be removed")
	flag.DurationVar(&opts.CleanPerfDataWaitTime, "cleanPerfDataWaitTime", 24*time.Hour, "Duration to keep deleted perf data in entities")
	flag.DurationVar(&opts.ExternalDataApiTimeout, "externalDataApiTimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.DefaultExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.IntVar(&opts.ExternalWorkers, "externalWorkers", canopsis.DefaultExternalEventWorkers, "Amount of workers to process external event flow.")
	flag.IntVar(&opts.SystemWorkers, "systemWorkers", canopsis.DefaultSystemEventWorkers, "Amount of workers to process system event flow.")
	flag.IntVar(&opts.UserWorkers, "userWorkers", canopsis.DefaultUserEventWorkers, "Amount of workers to process user event flow.")

	flag.Int("workers", 0, "Deprecated: Amount of workers to process each event flow.")
	flag.String("publishQueue", "", "Deprecated: publish event to this queue.")
	flag.String("consumeQueue", "", "Deprecated: consume events from this queue.")

	flag.Parse()

	return opts, libflag.FindDeprecatedFlags("workers", "publishQueue", "consumeQueue")
}
