package webhook

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	StatusCreated = iota
	StatusRunning
	StatusSucceeded
	StatusFailed
	StatusAborted
)

type History struct {
	ID                string   `bson:"_id" json:"_id"`
	Alarms            []string `bson:"alarms" json:"alarms"`
	Index             int64    `bson:"index" json:"index"`
	Scenario          string   `bson:"scenario,omitempty" json:"scenario,omitempty"`
	DeclareTicketRule string   `bson:"declare_ticket_rule,omitempty" json:"declare_ticket_rule,omitempty"`
	NextExec          string   `bson:"next_exec,omitempty" json:"next_exec,omitempty"`
	StopOnFail        bool     `bson:"stop_on_fail,omitempty" json:"stop_on_fail,omitempty"`
	Execution         string   `bson:"execution" json:"execution"`
	Name              string   `bson:"name" json:"name"`

	SystemName      string                        `bson:"system_name,omitempty" json:"system_name,omitempty"`
	EmitTrigger     bool                          `bson:"emit_trigger,omitempty" json:"emit_trigger,omitempty"`
	Status          int64                         `bson:"status" json:"status"`
	Comment         string                        `bson:"comment,omitempty" json:"comment,omitempty"`
	Request         request.Parameters            `bson:"request" json:"request"`
	DeclareTicket   *request.WebhookDeclareTicket `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	TicketResources bool                          `bson:"ticket_resources,omitempty" json:"ticket_resources,omitempty"`
	FailReason      string                        `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`
	RawRequest      string                        `bson:"raw_request,omitempty" json:"raw_request,omitempty"`
	RawResponse     string                        `bson:"raw_response,omitempty" json:"raw_response,omitempty"`
	ResponseCode    int64                         `bson:"response_code,omitempty" json:"response_code,omitempty"`
	ResponseHeader  map[string]string             `bson:"response_header,omitempty" json:"response_header,omitempty"`
	ResponseBody    map[string]any                `bson:"response_body,omitempty" json:"response_body,omitempty"`
	TicketID        string                        `bson:"ticket_id,omitempty" json:"ticket_id,omitempty"`
	TicketURL       string                        `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketData      map[string]string             `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`
	UserID          string                        `bson:"user,omitempty" json:"user,omitempty"`
	Username        string                        `bson:"username,omitempty" json:"username,omitempty"`
	Initiator       string                        `bson:"initiator,omitempty" json:"initiator,omitempty"`
	CreatedAt       datetime.CpsTime              `bson:"created_at" json:"created_at"`
	LaunchedAt      datetime.CpsTime              `bson:"launched_at,omitempty" json:"launched_at,omitempty"`
	CompletedAt     datetime.CpsTime              `bson:"completed_at,omitempty" json:"completed_at,omitempty"`

	EventInitiator string `bson:"event_initiator,omitempty" json:"event_initiator,omitempty"`
	EventOutput    string `bson:"event_output,omitempty" json:"event_output,omitempty"`
	Trigger        string `bson:"trigger,omitempty" json:"trigger,omitempty"`

	IsTest bool `bson:"is_test,omitempty" json:"is_test,omitempty"`
}

type TplAlarm struct {
	types.Alarm `bson:",inline"`
	Entity      types.Entity `bson:"entity" json:"entity"`
	Children    []struct {
		types.Alarm `bson:",inline"`
		Entity      types.Entity `bson:"entity" json:"entity"`
	} `bson:"children" json:"children"`
}

func FetchAlarmsForTplPipeline(ids []string) []bson.M {
	return []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$project": bson.M{"v.steps": 0}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$addFields": bson.M{
			"data": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "v.children",
			"foreignField": "d",
			"pipeline": []bson.M{
				{"$match": bson.M{"v.resolved": nil}},
				{"$project": bson.M{"v.steps": 0}},
			},
			"as": "children",
		}},
		{"$unwind": bson.M{
			"path":                       "$children",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "child_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "children.d",
			"foreignField": "_id",
			"as":           "children.entity",
		}},
		{"$unwind": bson.M{"path": "$children.entity", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"child_index": 1}},
		{"$group": bson.M{
			"_id":      "$_id",
			"data":     bson.M{"$first": "$data"},
			"children": bson.M{"$push": "$children"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"children": "$children"},
			}},
		}},
		{"$sort": bson.M{"t": -1}},
	}
}

func NewTplData(
	forMultiple bool,
	alarms []TplAlarm,
	additionalData types.AdditionalData,
	response, responseMap map[string]any,
	header map[string]string,
) map[string]any {
	res := make(map[string]any)
	if response != nil {
		res["Response"] = response
	}

	if responseMap != nil {
		res["ResponseMap"] = responseMap
	}

	if header != nil {
		res["Header"] = header
	}

	if forMultiple {
		res["Alarms"] = alarms
	} else if len(alarms) > 0 {
		res["Alarm"] = alarms[0].Alarm
		res["Entity"] = alarms[0].Entity
		res["Children"] = alarms[0].Children
	}

	res["AdditionalData"] = additionalData

	return res
}
