package webhook

import (
	"encoding/json"
	"errors"
	"fmt"

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

const MultipleURLsDelimiter = ","

type History struct {
	BaseHistory `bson:",inline"`

	Execution         string   `bson:"execution" json:"execution"`
	Alarms            []string `bson:"alarms,omitempty" json:"alarms,omitempty"`
	Scenario          string   `bson:"scenario,omitempty" json:"scenario,omitempty"`
	DeclareTicketRule string   `bson:"declare_ticket_rule,omitempty" json:"declare_ticket_rule,omitempty"`
	Name              string   `bson:"name" json:"name"`

	Index              int64  `bson:"index" json:"index"`
	NextExec           string `bson:"next_exec,omitempty" json:"next_exec,omitempty"`
	StopOnFail         bool   `bson:"stop_on_fail,omitempty" json:"stop_on_fail,omitempty"`
	StopOnSuccess      bool   `bson:"stop_on_success,omitempty" json:"stop_on_success,omitempty"`
	MultipleURLs       bool   `bson:"multiple_urls,omitempty" json:"multiple_urls,omitempty"`
	ResolvedRequestURL string `bson:"resolved_request_url,omitempty" json:"resolved_request_url,omitempty"`

	SystemName      string                    `bson:"system_name,omitempty" json:"system_name,omitempty"`
	EmitTrigger     bool                      `bson:"emit_trigger,omitempty" json:"emit_trigger,omitempty"`
	Comment         string                    `bson:"comment,omitempty" json:"comment,omitempty"`
	AuthToken       *request.WebhookAuthToken `bson:"auth_token,omitempty" json:"auth_token,omitempty"`
	DeclareTicket   *DeclareTicket            `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	TicketResources bool                      `bson:"ticket_resources,omitempty" json:"ticket_resources,omitempty"`
	UserID          string                    `bson:"user,omitempty" json:"user,omitempty"`
	Username        string                    `bson:"username,omitempty" json:"username,omitempty"`
	Initiator       string                    `bson:"initiator,omitempty" json:"initiator,omitempty"`
	EventInitiator  string                    `bson:"event_initiator,omitempty" json:"event_initiator,omitempty"`
	EventOutput     string                    `bson:"event_output,omitempty" json:"event_output,omitempty"`
	Trigger         string                    `bson:"trigger,omitempty" json:"trigger,omitempty"`

	ResponseCode   int64             `bson:"response_code,omitempty" json:"response_code,omitempty"`
	ResponseHeader map[string]string `bson:"response_header,omitempty" json:"response_header,omitempty"`
	ResponseBody   map[string]any    `bson:"response_body,omitempty" json:"response_body,omitempty"`

	TicketID   string            `bson:"ticket_id,omitempty" json:"ticket_id,omitempty"`
	TicketURL  string            `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketData map[string]string `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`

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
	responseTplVars ResponseTplVars,
) map[string]any {
	res := make(map[string]any)
	if responseTplVars.Response != nil {
		res["Response"] = responseTplVars.Response
	}

	if responseTplVars.ResponseMap != nil {
		res["ResponseMap"] = responseTplVars.ResponseMap
	}

	if responseTplVars.Header != nil {
		res["Header"] = responseTplVars.Header
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

type TokenHistory struct {
	BaseHistory        `bson:",inline"`
	Rule               string                    `bson:"rule" json:"rule"`
	ResponseField      string                    `bson:"response_field,omitempty" json:"response_field,omitempty"`
	Template           string                    `bson:"template,omitempty" json:"template,omitempty"`
	ExpirationDuration datetime.DurationWithUnit `bson:"expiration_duration" json:"expiration_duration"`
	Token              string                    `bson:"token,omitempty" json:"token,omitempty"`
	ExpiredAt          datetime.MicroTime        `bson:"expired_at,omitempty" json:"expired_at,omitempty"`
}

type BaseHistory struct {
	ID string `bson:"_id" json:"_id"`

	Status     int64  `bson:"status" json:"status"`
	FailReason string `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`

	Request     request.Parameters `bson:"request" json:"request"`
	RawRequest  string             `bson:"raw_request,omitempty" json:"raw_request,omitempty"`
	RawResponse string             `bson:"raw_response,omitempty" json:"raw_response,omitempty"`

	CreatedAt   datetime.MicroTime `bson:"created_at" json:"created_at"`
	LaunchedAt  datetime.MicroTime `bson:"launched_at,omitempty" json:"launched_at,omitempty"`
	CompletedAt datetime.MicroTime `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	LastPing    datetime.MicroTime `bson:"last_ping,omitempty" json:"last_ping,omitempty"`
}

type DeclareTicket struct {
	EmptyResponse     bool               `bson:"empty_response" json:"empty_response"`
	IsRegexp          bool               `bson:"is_regexp" json:"is_regexp"`
	TicketID          string             `bson:"ticket_id,omitempty" json:"ticket_id"`
	TicketIDTpl       string             `bson:"ticket_id_tpl,omitempty" json:"ticket_id_tpl" binding:"template"`
	TicketURL         string             `bson:"ticket_url,omitempty" json:"ticket_url"`
	TicketURLTpl      string             `bson:"ticket_url_tpl,omitempty" json:"ticket_url_tpl" binding:"template"`
	TicketURLTitle    string             `bson:"ticket_url_title,omitempty" json:"ticket_url_title"`
	CheckTicketStatus *CheckTicketStatus `bson:"check_ticket_status,omitempty" json:"check_ticket_status,omitempty"`
	CustomFields      map[string]string  `bson:",inline"`
}

func (t *DeclareTicket) UnmarshalJSON(b []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if emptyResponse, ok := m["empty_response"]; ok {
		if boolVal, ok := emptyResponse.(bool); ok {
			t.EmptyResponse = boolVal
			delete(m, "empty_response")
		} else {
			return errors.New("invalid type of empty_response")
		}
	}

	if isRegexp, ok := m["is_regexp"]; ok {
		if boolVal, ok := isRegexp.(bool); ok {
			t.IsRegexp = boolVal
			delete(m, "is_regexp")
		} else {
			return errors.New("invalid type of is_regexp")
		}
	}

	if checkTicketStatus, ok := m["check_ticket_status"]; ok {
		if checkTicketStatus != nil {
			raw, err := json.Marshal(checkTicketStatus)
			if err != nil {
				return fmt.Errorf("invalid type of check_ticket_status: %w", err)
			}

			t.CheckTicketStatus = &CheckTicketStatus{}
			if err := json.Unmarshal(raw, t.CheckTicketStatus); err != nil {
				return fmt.Errorf("invalid type of check_ticket_status: %w", err)
			}
		}

		delete(m, "check_ticket_status")
	}

	customFields := make(map[string]string)
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			switch k {
			case "ticket_id":
				t.TicketID = strVal
			case "ticket_id_tpl":
				t.TicketIDTpl = strVal
			case "ticket_url":
				t.TicketURL = strVal
			case "ticket_url_tpl":
				t.TicketURLTpl = strVal
			case "ticket_url_title":
				t.TicketURLTitle = strVal
			default:
				customFields[k] = strVal
			}
		} else {
			return fmt.Errorf("invalid type of %s", k)
		}
	}
	t.CustomFields = customFields

	return nil
}

func (t DeclareTicket) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"empty_response":   t.EmptyResponse,
		"is_regexp":        t.IsRegexp,
		"ticket_id":        t.TicketID,
		"ticket_id_tpl":    t.TicketIDTpl,
		"ticket_url":       t.TicketURL,
		"ticket_url_tpl":   t.TicketURLTpl,
		"ticket_url_title": t.TicketURLTitle,
	}

	if t.CheckTicketStatus != nil {
		m["check_ticket_status"] = t.CheckTicketStatus
	}

	for k, v := range t.CustomFields {
		m[k] = v
	}

	return json.Marshal(m)
}

type Webhook struct {
	Request       *request.Parameters       `bson:"request,omitempty" json:"request,omitempty"`
	AuthToken     *request.WebhookAuthToken `bson:"auth_token,omitempty" json:"auth_token,omitempty"`
	DeclareTicket *DeclareTicket            `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	StopOnFail    *bool                     `bson:"stop_on_fail,omitempty" json:"stop_on_fail,omitempty"`
	StopOnSuccess *bool                     `bson:"stop_on_success,omitempty" json:"stop_on_success,omitempty"`
	MultipleURLs  *bool                     `bson:"multiple_urls,omitempty" json:"multiple_urls,omitempty"`
}

type CheckTicketStatus struct {
	Request             request.Parameters        `bson:"request" json:"request"`
	AuthToken           *request.WebhookAuthToken `bson:"auth_token,omitempty" json:"auth_token,omitempty"`
	ReuseHeadersAndAuth bool                      `bson:"reuse_headers_and_auth" json:"reuse_headers_and_auth"`
	StatusMapping       map[string]int            `bson:"status_mapping" json:"status_mapping"`
	TicketStatus        string                    `bson:"ticket_status" json:"ticket_status"`
	TicketStatusTpl     string                    `bson:"ticket_status_tpl" json:"ticket_status_tpl"`

	ResolvedRequest request.Parameters `bson:"resolved_request" json:"-" binding:"-"`
}

type CheckTicketStatusJob struct {
	ID                     string `bson:"_id"`
	Status                 int    `bson:"status"`
	HistoryID              string `bson:"history_id"`
	TicketID               string `bson:"ticket_id"`
	TicketSystemName       string `bson:"ticket_system_name"`
	PrevTicketStatus       int    `bson:"prev_ticket_status"`
	TicketStatus           int    `bson:"ticket_status"`
	TicketSourceStatus     string `bson:"ticket_source_status"`
	PrevTicketSourceStatus string `bson:"prev_ticket_source_status"`

	AlarmIDs          []string          `bson:"alarm_ids"`
	CheckTicketStatus CheckTicketStatus `bson:"check_ticket_status"`

	CreatedAt datetime.CpsTime `bson:"created_at"`
	CheckedAt datetime.CpsTime `bson:"checked_at,omitempty"`
}

type ResponseTplVars struct {
	Header      map[string]string
	Response    map[string]any
	ResponseMap map[string]any
}
