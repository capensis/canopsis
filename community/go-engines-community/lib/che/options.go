package che

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

type Options struct {
	Version                 bool
	FeatureEventProcessing  bool
	FeatureContextCreation  bool
	Purge                   bool
	PrintEventOnError       bool
	ModeDebug               bool
	ConsumeFromQueue        string
	PublishToQueue          string
	DataSourceDirectory     string
	PeriodicalWaitTime      time.Duration
	InfosDictionaryWaitTime time.Duration
	ExternalDataApiTimeout  time.Duration
	FifoAckExchange         string
}

func ParseOptions() Options {
	opts := Options{}
	flag.BoolVar(&opts.FeatureEventProcessing, "processEvent", true, "enable event processing. enabled by default.")
	flag.BoolVar(&opts.FeatureContextCreation, "createContext", true, "enable context graph creation. enabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.AxeQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.CheQueueName, "Consume events from this queue.")
	flag.StringVar(&opts.DataSourceDirectory, "dataSourceDirectory", ".", "The path of the directory containing the event filter's data source plugins.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.Purge, "purge", false, "purge consumer queue(s) before work")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two runs of periodical process")
	flag.DurationVar(&opts.InfosDictionaryWaitTime, "infosDictionaryWaitTime", time.Hour, "Duration to wait between two runs of update entity infos dictionary process")
	flag.DurationVar(&opts.ExternalDataApiTimeout, "externalDataApiTimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")

	flag.Bool("enrichContext", false, "enable context graph enrichment from event. disabled by default. WARNING: disable the old context-graph engine when using this. - deprecated")
	flag.String("enrichExclude", "", "Coma separated list of fields that shall not be part of context enrichment. - deprecated")
	flag.String("enrichInclude", "", "Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored. - deprecated")

	flag.BoolVar(&opts.Version, "version", false, "Show the version information")

	flag.Parse()

	return opts
}
