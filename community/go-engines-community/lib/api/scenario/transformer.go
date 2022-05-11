package scenario

import (
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
		actions[i] = libaction.Action{
			Type:                     r.Actions[i].Type,
			Comment:                  r.Actions[i].Comment,
			Parameters:               r.Actions[i].Parameters,
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
		Priority:             *r.Priority,
		Delay:                r.Delay,
	}
}
