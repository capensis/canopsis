package entityservice

import (
	"github.com/go-redis/redis/v8"
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

type getStateCountersCmds struct {
	Critical *redis.StringCmd
	Major    *redis.StringCmd
	Minor    *redis.StringCmd
	Info     *redis.StringCmd
}

type getAlarmCountersCmds struct {
	All               *redis.StringCmd
	Alarms            *redis.StringCmd
	State             getStateCountersCmds
	Acknowledged      *redis.StringCmd
	NotAcknowledged   *redis.StringCmd
	PbehaviorCounters getPbehaviorCountersCmd
}

type getPbehaviorCountersCmd struct {
	All *redis.StringStringMapCmd
}

func (r getAlarmCountersCmds) Result() (AlarmCounters, error) {
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

	pbehaviorCounters := make(map[string]int64)
	for pbhType, v := range r.PbehaviorCounters.All.Val() {
		count, err := strconv.Atoi(v)
		if err != nil {
			return AlarmCounters{}, err
		}
		pbehaviorCounters[pbhType] = int64(count)
	}

	all, err := strconv.Atoi(r.All.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	alarms, err := strconv.Atoi(r.Alarms.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	acknowledged, err := strconv.Atoi(r.Acknowledged.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	notAcknowledged, err := strconv.Atoi(r.NotAcknowledged.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	critical, err := strconv.Atoi(r.State.Critical.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	major, err := strconv.Atoi(r.State.Major.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	minor, err := strconv.Atoi(r.State.Minor.Val())
	if err != nil {
		return AlarmCounters{}, err
	}
	info, err := strconv.Atoi(r.State.Info.Val())
	if err != nil {
		return AlarmCounters{}, err
	}

	return AlarmCounters{
		All:    int64(all),
		Alarms: int64(alarms),
		State: StateCounters{
			Critical: int64(critical),
			Major:    int64(major),
			Minor:    int64(minor),
			Info:     int64(info),
		},
		Acknowledged:      int64(acknowledged),
		NotAcknowledged:   int64(notAcknowledged),
		PbehaviorCounters: pbehaviorCounters,
	}, nil
}
