package che

import (
	"flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"time"
)

type Options struct {
	FeatureEventProcessing  bool
	FeatureContextCreation  bool
	FeatureContextEnrich    bool
	Purge                   bool
	PrintEventOnError       bool
	ModeDebug               bool
	ConsumeFromQueue        string
	PublishToQueue          string
	EnrichExclude           string
	EnrichInclude           string
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
	flag.BoolVar(&opts.FeatureContextEnrich, "enrichContext", false, "enable context graph enrichment from event. disabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.StringVar(&opts.PublishToQueue, "publishQueue", canopsis.AxeQueueName, "Publish event to this queue.")
	flag.StringVar(&opts.ConsumeFromQueue, "consumeQueue", canopsis.CheQueueName, "Consume events from this queue.")
	flag.StringVar(&opts.EnrichExclude, "enrichExclude", "", "Coma separated list of fields that shall not be part of context enrichment.")
	flag.StringVar(&opts.EnrichInclude, "enrichInclude", "", "Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored.")
	flag.StringVar(&opts.DataSourceDirectory, "dataSourceDirectory", ".", "The path of the directory containing the event filter's data source plugins.")
	flag.BoolVar(&opts.ModeDebug, "d", false, "debug")
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.Purge, "purge", false, "purge consumer queue(s) before work")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two runs of periodical process")
	flag.DurationVar(&opts.InfosDictionaryWaitTime, "infosDictionaryWaitTime", time.Hour, "Duration to wait between two runs of update entity infos dictionary process")
	flag.DurationVar(&opts.ExternalDataApiTimeout, "externalDataApiTimeout", 30*time.Second, "External API HTTP Request Timeout.")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.FIFOAckExchangeName, "Publish FIFO Ack event to this exchange.")

	flagVersion := flag.Bool("version", false, "version infos")

	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	return opts
}
