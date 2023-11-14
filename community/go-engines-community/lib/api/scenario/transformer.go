package scenario

import (
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	libaction "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
			Type:       r.Actions[i].Type,
			Comment:    r.Actions[i].Comment,
			Parameters: r.Actions[i].Parameters,
			EntityPatternFields: r.Actions[i].EntityPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInEntityPattern(mongo.ScenarioMongoCollection),
			),
			AlarmPatternFields: r.Actions[i].AlarmPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInAlarmPattern(mongo.ScenarioMongoCollection),
				common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ScenarioMongoCollection),
			),
			DropScenarioIfNotMatched: *r.Actions[i].DropScenarioIfNotMatched,
			EmitTrigger:              *r.Actions[i].EmitTrigger,
		}
	}

	triggers := make([]string, 0, len(r.Triggers))
	for _, triggerRequest := range r.Triggers {
		if triggerRequest.Type == string(types.AlarmChangeEventsCount) {
			triggers = append(triggers, triggerRequest.Type+strconv.Itoa(triggerRequest.Threshold))
		} else {
			triggers = append(triggers, triggerRequest.Type)
		}
	}

	return libaction.Scenario{
		Name:                 r.Name,
		Author:               r.Author,
		Enabled:              *r.Enabled,
		DisableDuringPeriods: r.DisableDuringPeriods,
		Triggers:             triggers,
		Actions:              actions,
		Priority:             r.Priority,
		Delay:                r.Delay,
	}
}
