package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func (f *Feeder) modeFeeder(ctx context.Context) error {
	if err := f.setupAmqp(); err != nil {
		return err
	}
	nEntities := (f.flags.NComp) * (f.flags.NRes) * (f.flags.NConn)
	eventsPerSec := float64(nEntities) / f.flags.Freq.Seconds()
	nanosecSleep := time.Duration(float64(time.Second.Nanoseconds()) / eventsPerSec * float64(f.flags.NRes))
	pubcount := int64(0)

	f.logger.Info().Msgf("components: %d\n", f.flags.NComp)
	f.logger.Info().Msgf("resources: %d\n", f.flags.NRes)
	f.logger.Info().Msgf("connectors: %d\n", f.flags.NConn)
	f.logger.Info().Msgf("full loop event count: %d\n", nEntities)
	f.logger.Info().Msgf("entities per connector: %d\n", (f.flags.NComp)*(f.flags.NRes))
	f.logger.Info().Msgf("context graph entities: %d (%d resource/component + %d connectors + %d comp)\n",
		(f.flags.NConn + f.flags.NComp*(f.flags.NRes) + f.flags.NComp),
		f.flags.NComp*(f.flags.NRes),
		f.flags.NConn,
		f.flags.NComp)
	f.logger.Info().Msgf("events per second: %f\n", eventsPerSec)

	tstart := time.Now().UnixNano()
	checkEvery := int64(100)

	changeStateEvery := int64(100 / f.flags.Alarms)

	rand.Seed(time.Now().UnixNano())

	stateMap := make(map[string]int)

	f.logger.Info().Msg("pushing events")

	for {
		for ci := f.flags.CompStart; ci < f.flags.NComp; ci++ {
			for Ci := f.flags.ConnStart; Ci < f.flags.NConn; Ci++ {
				for ri := f.flags.ResStart; ri < f.flags.NRes; ri++ {
					eid := fmt.Sprintf("%d%d", ci, ri)
					state := stateMap[eid]

					if (ci*Ci*ri)%changeStateEvery == 0 {
						if state == types.AlarmStateOK {
							state = types.AlarmStateCritical
						} else {
							state = types.AlarmStateOK
						}
						//log.Printf("change state %d %d: %d -> %d", ci, ri, stateMap[eid], state)
					}

					stateMap[eid] = state

					err := f.send(ctx, int64(state), Ci, ci, ri)
					if err != nil {
						return err
					}
					pubcount++

					if pubcount%checkEvery == 0 {
						tsent := time.Now().UnixNano() - tstart
						adj := f.adjust(eventsPerSec, checkEvery, tsent)
						nanosecSleep = time.Duration(nanosecSleep.Nanoseconds() + adj)
						tstart = time.Now().UnixNano()
					}

				}

				if nanosecSleep.Nanoseconds() > 0 {
					time.Sleep(nanosecSleep)
				}
				f.flags.ResStart = int64(0)
			}
			f.flags.ConnStart = int64(0)
		}
		f.flags.CompStart = int64(0)
		f.logger.Info().Msg("full loop.")
	}
}
