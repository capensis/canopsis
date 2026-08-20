package che

import (
	"flag"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libflag "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

type Options struct {
	log.Options
	Version                 bool
	Purge                   bool
	PrintEventOnError       bool
	PeriodicalWaitTime      time.Duration
	InfosDictionaryWaitTime time.Duration
	SoftDeleteWaitTime      time.Duration
	FifoAckExchange         string
	ExternalWorkers         int
	SystemWorkers           int
	UserWorkers             int
	RpcWorkers              int
}

func ParseOptions() (Options, []string) {
	opts := Options{}
	log.BindCmdFlags(&opts.Options)
	flag.BoolVar(&opts.PrintEventOnError, "printEventOnError", false, "Print event on processing error")
	flag.BoolVar(&opts.Purge, "purge", false, "purge consumer queue(s) before work")
	flag.DurationVar(&opts.PeriodicalWaitTime, "periodicalWaitTime", canopsis.PeriodicalWaitTime, "Duration to wait between two runs of periodical process")
	flag.DurationVar(&opts.InfosDictionaryWaitTime, "infosDictionaryWaitTime", time.Hour, "Duration to wait between two runs of update entity infos dictionary process")
	flag.DurationVar(&opts.SoftDeleteWaitTime, "softDeleteWaitTime", time.Hour, "Duration to keep soft deleted entities in the db until they will be removed")
	flag.StringVar(&opts.FifoAckExchange, "fifoAckExchange", canopsis.DefaultExchangeName, "Publish FIFO Ack event to this exchange.")
	flag.BoolVar(&opts.Version, "version", false, "Show the version information")
	flag.IntVar(&opts.ExternalWorkers, "externalWorkers", canopsis.DefaultExternalEventWorkers, "Amount of workers to process external event flow.")
	flag.IntVar(&opts.SystemWorkers, "systemWorkers", canopsis.DefaultSystemEventWorkers, "Amount of workers to process system event flow.")
	flag.IntVar(&opts.UserWorkers, "userWorkers", canopsis.DefaultUserEventWorkers, "Amount of workers to process user event flow.")
	flag.IntVar(&opts.RpcWorkers, "rpcWorkers", canopsis.DefaultRpcWorkers, "Amount of workers to process rpc event flow.")

	flag.Bool("processEvent", true, "Deprecated: enable event processing. enabled by default.")
	flag.Bool("createContext", true, "Deprecated: enable context graph creation. enabled by default. WARNING: disable the old context-graph engine when using this.")
	flag.Duration("externalDataApiTimeout", 30*time.Second, "Deprecated: External API HTTP Request Timeout.")

	flag.Parse()

	return opts, libflag.FindDeprecatedFlags("processEvent", "createContext", "externalDataApiTimeout")
}
