package pattern

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func Match(
	entity types.Entity,
	alarm types.Alarm,
	entityPattern Entity,
	alarmPattern Alarm,
) (bool, error) {
	if len(entityPattern) == 0 && len(alarmPattern) == 0 {
		return false, nil
	}

	if len(entityPattern) > 0 {
		ok, err := entityPattern.Match(entity)
		if err != nil {
			return false, fmt.Errorf("entity pattern is invalid: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	if len(alarmPattern) > 0 {
		ok, err := alarmPattern.Match(alarm)
		if err != nil {
			return false, fmt.Errorf("alarm pattern is invalid: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}
