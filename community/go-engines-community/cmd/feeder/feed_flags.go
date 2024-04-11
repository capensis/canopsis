package main

import (
	"errors"
	"flag"
	"os"
	"time"

	cps "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

func (f *Flags) ParseArgs() error {
	flagNComp := flag.Int64("ncomp", int64(1200), "number of components")
	flagNConn := flag.Int64("nconn", int64(10), "number of connectors")
	flagNRes := flag.Int64("nres", int64(10), "number of resources")
	flagFreq := flag.Duration("freq", time.Minute*5, "each comp/res will send an event each <duration>")
	flagAlarms := flag.Int("alarms", 20, "percent of alarms")
	flagCompStart := flag.Int64("compstart", int64(0), "start with this component")
	flagConnStart := flag.Int64("connstart", int64(0), "start with this connector")
	flagResStart := flag.Int64("resstart", int64(0), "start with this resource")
	flagDirtyEvent := flag.Bool("dirty", true, "dirty event to be cleaned by che. if not, compat event for old engines")
	flagVersion := flag.Bool("version", false, "Show the version information")
	flagFile := flag.String("file", "event.json", "send event from json file (-amqp|http) or directory (-http)")
	flagMode := flag.String("mode", "feeder", "mode: file (-amqp or -http) or feeder")
	flagExchange := flag.String("exchange", cps.CheExchangeName, "exchange name to publish events to")
	flagLoop := flag.Bool("loop", false, "constantly send events")

	flagHTTP := flag.Bool("http", false, "publish events overt http if -file is used")
	flagAMQP := flag.Bool("amqp", false, "publish events overt amqp if -file is used")
	flagHTTPURL := flag.String("url", os.Getenv("CPS_URL"), "http publish url. defaults to CPS_URL env var")
	flagAuthKey := flag.String("authkey", os.Getenv("CPS_AUTH_KEY"), "canopsis auth key. cannot be empty when using -http. defaults to CPS_AUTH_KEY")

	flagCheckJSON := flag.Bool("checkJSON", true, "check json validity before sending. defaults to true")

	flag.Parse()

	f.NComp = *flagNComp
	f.NConn = *flagNConn
	f.NRes = *flagNRes
	f.Freq = *flagFreq
	f.Alarms = *flagAlarms
	f.CompStart = *flagCompStart
	f.ConnStart = *flagConnStart
	f.ResStart = *flagResStart
	f.DirtyEvent = *flagDirtyEvent
	f.Version = *flagVersion
	f.File = *flagFile
	f.Mode = *flagMode
	f.ExchangeName = *flagExchange
	f.PubAMQP = *flagAMQP
	f.PubHTTP = *flagHTTP
	f.PubHTTPURL = *flagHTTPURL
	f.AuthKey = *flagAuthKey
	f.CheckJSON = *flagCheckJSON
	f.Loop = *flagLoop

	if f.PubAMQP && f.PubHTTP {
		return errors.New("-amqp and -http given")
	}

	return nil
}

type Flags struct {
	NComp        int64
	NConn        int64
	NRes         int64
	Freq         time.Duration
	Alarms       int
	CompStart    int64
	ConnStart    int64
	ResStart     int64
	DirtyEvent   bool
	Version      bool
	File         string
	Mode         string
	ExchangeName string
	PubAMQP      bool
	PubHTTP      bool
	PubHTTPURL   string
	AuthKey      string
	CheckJSON    bool
	Loop         bool
}
