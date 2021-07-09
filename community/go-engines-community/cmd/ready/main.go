package main

import (
	"flag"
	"log"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils/ready"
)

func main() {
	var flagRetries int
	var flagTimeout time.Duration

	flag.IntVar(&flagRetries, "retries", 10, "number of retries per check. if 0, infinite number of retries")
	flag.DurationVar(&flagTimeout, "timeout", time.Second*60, "timeout after given duration. never timeout if 0s")

	flag.Parse()

	if flagTimeout != time.Second*0 {
		go ready.Timeout(flagTimeout)
	}

	ready.Abort(ready.Check(ready.CheckRedis, "redis", time.Second, flagRetries))
	ready.Abort(ready.Check(ready.CheckMongo, "mongo", time.Second, flagRetries))
	ready.Abort(ready.Check(ready.CheckAMQP, "amqp", time.Second, flagRetries))

	log.Println("ready!")
}
