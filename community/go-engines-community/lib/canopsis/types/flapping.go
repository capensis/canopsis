package types

import (
	"sync"
	"time"
)

var flapping *flappingCheck

func init() {
	flapping = &flappingCheck{
		check: func(a *Alarm) bool {
			return false
		},
	}
}

type flappingCheck struct {
	mu    sync.RWMutex
	check func(a *Alarm) bool
}

func SetFlapping(check func(a *Alarm) bool) {
	flapping.mu.Lock()
	defer flapping.mu.Unlock()
	flapping.check = check
}

func IsFlapping(a *Alarm) bool {
	flapping.mu.RLock()
	defer flapping.mu.RUnlock()
	return flapping.check(a)
}

func IsFlappingWithDurationAndStep(a *Alarm, interval time.Duration, step int) bool {
	lastStepType := ""
	freq := 0

	for i := len(a.Value.Steps) - 1; i >= 0; i-- {
		s := a.Value.Steps[i]
		duration := time.Since(s.Timestamp.Time)
		if duration >= interval {
			break
		}

		if s.Type != lastStepType {
			switch s.Type {
			case AlarmStepStateIncrease, AlarmStepStateDecrease:
				lastStepType = s.Type
				freq++
			}
		}

		if freq > step {
			return true
		}
	}

	return false
}
