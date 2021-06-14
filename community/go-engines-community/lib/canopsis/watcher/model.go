package watcher

import (
	"fmt"
	"strings"
	"text/template"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// Watcher is a structure representing a watcher type entity document.
type Watcher struct {
	types.Entity   `bson:",inline"`          // inherits from entity
	Entities       pattern.EntityPatternList `bson:"entities"`
	State          map[string]interface{}    `bson:"state"`
	OutputTemplate string                    `bson:"output_template"`
}

// CheckEntityInWatcher checks if the entity is watched by the watcher. It
// returns true when the entity is matched by a pattern in the watcher, false
// otherwise.
func (w Watcher) CheckEntityInWatcher(entity types.Entity) bool {
	return w.Entities.Matches(&entity)
}

// GetOutput returns the output of the watcher.
func (w Watcher) GetOutput(counters AlarmCounters) (string, error) {
	tpl, err := template.New("template").Parse(w.OutputTemplate)
	if err != nil {
		return "", fmt.Errorf(
			"unable to parse output template for watcher %s: %+v",
			w.ID, err)
	}

	b := strings.Builder{}
	err = tpl.Execute(&b, counters)
	if err != nil {
		return "", fmt.Errorf(
			"unable to execute output template for watcher %s: %+v",
			w.ID, err)
	}

	return b.String(), nil
}

// GetState returns the state of the watcher.
func (w Watcher) GetState(counters AlarmCounters) (int, error) {
	switch w.State["method"] {
	case MethodWorst:
		return worst(counters), nil
	default:
		return 0, fmt.Errorf("unknown watcher method : %s", w.State["method"])
	}
}
