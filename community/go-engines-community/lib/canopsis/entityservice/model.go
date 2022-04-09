package entityservice

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"html/template"
	"strings"
)

type EntityService struct {
	types.Entity   `bson:",inline"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	OutputTemplate string                    `bson:"output_template" json:"output_template"`
}

// GetServiceOutput returns the output of the service.
func GetServiceOutput(outputTemplate string, counters AlarmCounters) (string, error) {
	tpl, err := template.New("template").Parse(outputTemplate)
	if err != nil {
		return "", fmt.Errorf(
			"unable to parse output template for service %s: %w", outputTemplate, err)
	}

	b := strings.Builder{}
	err = tpl.Execute(&b, counters)
	if err != nil {
		return "", fmt.Errorf(
			"unable to execute output template for service %s: %w",
			outputTemplate, err)
	}

	return b.String(), nil
}

// GetServiceState returns the state of the service.
func GetServiceState(counters AlarmCounters) int {
	if counters.State.Critical > 0 {
		return types.AlarmStateCritical
	}
	if counters.State.Major > 0 {
		return types.AlarmStateMajor
	}
	if counters.State.Minor > 0 {
		return types.AlarmStateMinor
	}

	return types.AlarmStateOK
}
