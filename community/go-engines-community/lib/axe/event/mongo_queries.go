package event

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func getOpenAlarmMatch(event rpc.AxeEvent) bson.M {
	m := getExactOpenAlarmMatch(event)
	if m != nil {
		return m
	}

	return bson.M{
		"d":          event.Entity.ID,
		"v.resolved": nil,
	}
}

func getExactAlarmMatch(event rpc.AxeEvent) bson.M {
	if event.Alarm != nil {
		return bson.M{"_id": event.Alarm.ID}
	}

	if event.AlarmID != "" {
		return bson.M{"_id": event.AlarmID}
	}

	return nil
}

func getExactOpenAlarmMatch(event rpc.AxeEvent) bson.M {
	m := getExactAlarmMatch(event)
	if m == nil {
		return nil
	}

	m["v.resolved"] = nil

	return m
}

func getOpenAlarmMatchWithStepsLimit(event rpc.AxeEvent) bson.M {
	match := getOpenAlarmMatch(event)
	match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}

	return match
}

func getExactAlarmMatchWithStepsLimit(event rpc.AxeEvent) bson.M {
	match := getExactAlarmMatch(event)
	if match == nil {
		return nil
	}

	match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}

	return match
}

func stepUpdateQueryWithInPbhInterval(stepType, msg string, params rpc.AxeParameters) bson.M {
	newStep := NewAlarmStep(stepType, params, false)
	newStep.Message = msg

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func valStepUpdateQueryWithInPbhInterval(stepType string, value types.CpsNumber, msg string, params rpc.AxeParameters) bson.M {
	newStep := NewAlarmStep(stepType, params, false)
	newStep.Message = msg
	newStep.Value = value

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func execStepUpdateQueryWithInPbhInterval(stepType, displayGroup, msg string, params rpc.AxeParameters) bson.M {
	newStep := NewAlarmStep(stepType, params, false)
	newStep.Message = msg
	newStep.Execution = params.Execution
	newStep.DisplayGroup = displayGroup

	return stepUpdateQueryWithInPbhIntervalByStep(newStep)
}

func ticketStepUpdateQueryWithInPbhInterval(stepType string, displayGroup, msg string, params rpc.AxeParameters) bson.M {
	newStep := NewAlarmStep(stepType, params, false)
	newStep.Message = msg
	newStep.TicketInfo = params.TicketInfo
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
			bson.M{"$literal": newStep},
			bson.M{"in_pbh": true},
		}},
		"else": bson.M{"$literal": newStep},
	}}
}

func addStepUpdateQuery(newStepQueries ...bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{"$v.steps", newStepQueries}}
}

func addTicketUpdateQuery(newStepQuery bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{
		bson.M{"$ifNull": bson.A{"$v.tickets", bson.A{}}},
		bson.A{newStepQuery},
	}}
}

func addCommentsUpdateQuery(newStepQuery bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{
		bson.M{"$ifNull": bson.A{"$v.comments", bson.A{}}},
		bson.A{newStepQuery},
	}}
}
