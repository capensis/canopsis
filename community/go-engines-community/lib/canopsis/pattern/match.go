package pattern

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func Match(
	alarm types.Alarm,
	entity types.Entity,
	entityPattern Entity,
	alarmPattern Alarm,
	oldEntityPatterns oldpattern.EntityPatternList,
	oldAlarmPatterns oldpattern.AlarmPatternList,
) (bool, error) {
	if !oldEntityPatterns.IsSet() && len(entityPattern) == 0 &&
		!oldAlarmPatterns.IsSet() && len(alarmPattern) == 0 {
		return false, nil
	}

	if len(entityPattern) > 0 {
		ok, _, err := entityPattern.Match(entity)
		if err != nil || !ok {
			return false, err
		}
	} else if oldEntityPatterns.IsSet() {
		if !oldEntityPatterns.IsValid() {
			return false, errors.New("old entity pattern is not valid")
		}
		if !oldEntityPatterns.Matches(&entity) {
			return false, nil
		}
	}

	if len(alarmPattern) > 0 {
		ok, err := alarmPattern.Match(alarm)
		if err != nil || !ok {
			return false, err
		}
	} else if oldAlarmPatterns.IsSet() {
		if !oldAlarmPatterns.IsValid() {
			return false, errors.New("old alarm pattern is not valid")
		}
		if !oldAlarmPatterns.Matches(&alarm) {
			return false, nil
		}
	}

	return true, nil
}
