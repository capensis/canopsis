package scenario

import (
	"encoding/json"
	libaction "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
)

type ModelTransformer interface {
	TransformEditRequestToModel(request EditRequest) libaction.Scenario
}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

type modelTransformer struct{}

func (modelTransformer) TransformEditRequestToModel(r EditRequest) libaction.Scenario {
	actions := make([]libaction.Action, len(r.Actions))
	for i := range actions {
		var params map[string]interface{}

		if r.Actions[i].Parameters != nil {
			b, _ := json.Marshal(r.Actions[i].Parameters)
			_ = json.Unmarshal(b, &params)
		}

		actions[i] = libaction.Action{
			Type:                     r.Actions[i].Type,
			Comment:                  r.Actions[i].Comment,
			Parameters:               params,
			AlarmPatterns:            r.Actions[i].AlarmPatterns,
			EntityPatterns:           r.Actions[i].EntityPatterns,
			DropScenarioIfNotMatched: *r.Actions[i].DropScenarioIfNotMatched,
			EmitTrigger:              *r.Actions[i].EmitTrigger,
		}
	}

	return libaction.Scenario{
		Name:                 r.Name,
		Author:               r.Author,
		Enabled:              *r.Enabled,
		DisableDuringPeriods: r.DisableDuringPeriods,
		Triggers:             r.Triggers,
		Actions:              actions,
		Priority:             r.Priority,
		Delay:                r.Delay,
	}
}
