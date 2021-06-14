package watcher

import (
	"github.com/go-redis/redis/v7"
	"strconv"
)

// stateCounterCmds is a type containing redis command results for each field
// of a StateCounters struct. See alarmCounterCmds for more details.
type stateCountersCmds struct {
	Critical *redis.IntCmd
	Major    *redis.IntCmd
	Minor    *redis.IntCmd
	Info     *redis.IntCmd
}

// alarmCounterCmds is a type containing redis command results for each field
// of an AlarmCounters struct.
// This type is used in countersCache to increment the fields of an
// AlarmCounters stored in redis in a transaction, and to read the result of
// the commands after the transaction has been executed.
type alarmCountersCmds struct {
	All               *redis.IntCmd
	Alarms            *redis.IntCmd
	State             stateCountersCmds
	Acknowledged      *redis.IntCmd
	NotAcknowledged   *redis.IntCmd
	PbehaviorCounters pbehaviorCountersCmd
}

type pbehaviorCountersCmd struct {
	All  *redis.StringStringMapCmd
	Incr map[string]*redis.IntCmd
}

// Result reads the result of the redis commands, and return them in an
// AlarmCounters struct.
// If any of the redis command failed, an AlarmCounters with all the fields set
// to zero is returned, as well as the error.
func (r alarmCountersCmds) Result() (AlarmCounters, error) {
	if r.Alarms.Err() != nil {
		return AlarmCounters{}, r.Alarms.Err()
	}
	if r.State.Critical.Err() != nil {
		return AlarmCounters{}, r.State.Critical.Err()
	}
	if r.State.Major.Err() != nil {
		return AlarmCounters{}, r.State.Major.Err()
	}
	if r.State.Minor.Err() != nil {
		return AlarmCounters{}, r.State.Minor.Err()
	}
	if r.State.Info.Err() != nil {
		return AlarmCounters{}, r.State.Info.Err()
	}
	if r.Acknowledged.Err() != nil {
		return AlarmCounters{}, r.Acknowledged.Err()
	}
	if r.NotAcknowledged.Err() != nil {
		return AlarmCounters{}, r.NotAcknowledged.Err()
	}
	if r.PbehaviorCounters.All.Err() != nil {
		return AlarmCounters{}, r.PbehaviorCounters.All.Err()
	}

	for _, cmd := range r.PbehaviorCounters.Incr {
		if err := cmd.Err(); err != nil {
			return AlarmCounters{}, err
		}
	}

	all := r.PbehaviorCounters.All.Val()
	pbehaviorCounters := make(map[string]int64)
	for pbhType, v := range all {
		count, err := strconv.Atoi(v)
		if err != nil {
			return AlarmCounters{}, err
		}
		pbehaviorCounters[pbhType] = int64(count)
	}

	return AlarmCounters{
		All:    r.All.Val(),
		Alarms: r.Alarms.Val(),
		State: StateCounters{
			Critical: r.State.Critical.Val(),
			Major:    r.State.Major.Val(),
			Minor:    r.State.Minor.Val(),
			Info:     r.State.Info.Val(),
		},
		Acknowledged:      r.Acknowledged.Val(),
		NotAcknowledged:   r.NotAcknowledged.Val(),
		PbehaviorCounters: pbehaviorCounters,
	}, nil
}
