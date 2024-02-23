package event

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func getOpenAlarmMatch(event rpc.AxeEvent) bson.M {
	if event.Alarm != nil {
		return bson.M{
			"_id":        event.Alarm.ID,
			"v.resolved": nil,
		}
	}

	if event.AlarmID != "" {
		return bson.M{
			"_id":        event.AlarmID,
			"v.resolved": nil,
		}
	}

	return bson.M{
		"d":          event.Entity.ID,
		"v.resolved": nil,
	}
}

func getOpenAlarmMatchWithStepsLimit(event rpc.AxeEvent) bson.M {
	match := getOpenAlarmMatch(event)
	match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}

	return match
}

func stepUpdateQueryWithInPbhInterval(stepType, msg string, params rpc.AxeParameters) bson.M {
	newStep := types.NewAlarmStep(stepType, params.Timestamp, params.Author, msg, params.User, params.Role,
		params.Initiator, false)

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func valStepUpdateQueryWithInPbhInterval(stepType string, value types.CpsNumber, msg string, params rpc.AxeParameters) bson.M {
	newStep := types.NewAlarmStep(stepType, params.Timestamp, params.Author, msg, params.User, params.Role,
		params.Initiator, false)
	newStep.Value = value

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func execStepUpdateQueryWithInPbhInterval(stepType, displayGroup, msg string, params rpc.AxeParameters) bson.M {
	newStep := types.NewAlarmStep(stepType, params.Timestamp, params.Author, msg, params.User, params.Role,
		params.Initiator, false)
	newStep.Execution = params.Execution
	newStep.DisplayGroup = displayGroup

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func ticketStepUpdateQueryWithInPbhInterval(stepType string, displayGroup, msg string, params rpc.AxeParameters) bson.M {
	newStep := types.NewTicketStep(stepType, params.Timestamp, params.Author, msg, params.User, params.Role,
		params.Initiator, params.TicketInfo, false)
	newStep.Execution = params.Execution
	newStep.DisplayGroup = displayGroup

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func stepUpdateQueryWithInPbhIntervalByStep(newStep types.AlarmStep) bson.M {
	return bson.M{"$cond": bson.M{
		"if": bson.M{"$and": []bson.M{
			{"$eq": bson.A{bson.M{"$type": "$v.pbehavior_info.id"}, "string"}},
			{"$ne": bson.A{"$v.pbehavior_info.id", ""}},
		}},
		"then": bson.M{"$mergeObjects": bson.A{
			newStep,
			bson.M{"in_pbh": true},
		}},
		"else": newStep,
	}}
}

func addStepUpdateQuery(newStepQueries ...bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{"$v.steps", newStepQueries}}
}

func addTicketUpdateQuery(newStepQuery bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{
		bson.M{"$cond": bson.M{
			"if":   "$v.tickets",
			"then": "$v.tickets",
			"else": bson.A{},
		}},
		bson.A{newStepQuery},
	}}
}
