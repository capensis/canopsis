package entityservice

import (
	"strconv"

	"github.com/go-redis/redis/v8"
)

// stateCounterCmds is a type containing redis command results for each field
// of a StateCounters struct. See alarmCounterCmds for more details.
type stateCountersCmds struct {
	Critical *redis.IntCmd
	Major    *redis.IntCmd
	Minor    *redis.IntCmd
	Ok       *redis.IntCmd
}

// alarmCounterCmds is a type containing redis command results for each field
// of an AlarmCounters struct.
// This type is used in countersCache to increment the fields of an
// AlarmCounters stored in redis in a transaction, and to read the result of
// the commands after the transaction has been executed.
type alarmCountersCmds struct {
	All                  *redis.IntCmd
	Active               *redis.IntCmd
	State                stateCountersCmds
	Acknowledged         *redis.IntCmd
	NotAcknowledged      *redis.IntCmd
	AcknowledgedUnderPbh *redis.IntCmd
	PbehaviorCounters    pbehaviorCountersCmd
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
	if r.Active.Err() != nil {
		return AlarmCounters{}, r.Active.Err()
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
	if r.State.Ok.Err() != nil {
		return AlarmCounters{}, r.State.Ok.Err()
	}
	if r.Acknowledged.Err() != nil {
		return AlarmCounters{}, r.Acknowledged.Err()
	}
	if r.NotAcknowledged.Err() != nil {
		return AlarmCounters{}, r.NotAcknowledged.Err()
	}
	if r.AcknowledgedUnderPbh.Err() != nil {
		return AlarmCounters{}, r.AcknowledgedUnderPbh.Err()
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
		Active: r.Active.Val(),
		State: StateCounters{
			Critical: r.State.Critical.Val(),
			Major:    r.State.Major.Val(),
			Minor:    r.State.Minor.Val(),
			Ok:       r.State.Ok.Val(),
		},
		Acknowledged:         r.Acknowledged.Val(),
		NotAcknowledged:      r.NotAcknowledged.Val(),
		AcknowledgedUnderPbh: r.AcknowledgedUnderPbh.Val(),
		PbehaviorCounters:    pbehaviorCounters,
	}, nil
}

type getStateCountersCmds struct {
	Critical *redis.StringCmd
	Major    *redis.StringCmd
	Minor    *redis.StringCmd
	Ok       *redis.StringCmd
}

type getAlarmCountersCmds struct {
	All                  *redis.StringCmd
	Active               *redis.StringCmd
	State                getStateCountersCmds
	Acknowledged         *redis.StringCmd
	NotAcknowledged      *redis.StringCmd
	AcknowledgedUnderPbh *redis.StringCmd
	PbehaviorCounters    getPbehaviorCountersCmd
}

type getPbehaviorCountersCmd struct {
	All *redis.StringStringMapCmd
}

func (r getAlarmCountersCmds) Result() (AlarmCounters, error) {
	if r.Active.Err() != nil {
		return AlarmCounters{}, r.Active.Err()
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
	if r.State.Ok.Err() != nil {
		return AlarmCounters{}, r.State.Ok.Err()
	}
	if r.Acknowledged.Err() != nil {
		return AlarmCounters{}, r.Acknowledged.Err()
	}
	if r.NotAcknowledged.Err() != nil {
		return AlarmCounters{}, r.NotAcknowledged.Err()
	}
	if r.AcknowledgedUnderPbh.Err() != nil {
		return AlarmCounters{}, r.AcknowledgedUnderPbh.Err()
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
	active, err := strconv.Atoi(r.Active.Val())
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
	acknowledgedUnderPbh, err := strconv.Atoi(r.AcknowledgedUnderPbh.Val())
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
	ok, err := strconv.Atoi(r.State.Ok.Val())
	if err != nil {
		return AlarmCounters{}, err
	}

	return AlarmCounters{
		All:    int64(all),
		Active: int64(active),
		State: StateCounters{
			Critical: int64(critical),
			Major:    int64(major),
			Minor:    int64(minor),
			Ok:       int64(ok),
		},
		Acknowledged:         int64(acknowledged),
		NotAcknowledged:      int64(notAcknowledged),
		AcknowledgedUnderPbh: int64(acknowledgedUnderPbh),
		PbehaviorCounters:    pbehaviorCounters,
	}, nil
}
