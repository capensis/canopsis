package patternfields

import (
	"cmp"
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	aliasLimit = 500
)

var alarmFields = []FieldResponse{
	{
		Name: "v.display_name",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.state.val",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "v.status.val",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "v.component",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.connector",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.connector_name",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.resource",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.creation_date",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.duration",
		Type: pattern.FieldTypeDuration,
	},
	{
		Name: "v.infos",
		Type: pattern.FieldTypeObject,
	},
	{
		Name: "v.output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.last_event_date",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.last_update_date",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.ack",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.ack.t",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.ack.a",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.ack.m",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.ack.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.resolved",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.ticket",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.ticket.ticket",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.ticket.m",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.ticket.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.ticket.ticket_data",
		Type: pattern.FieldTypeObject,
	},
	{
		Name: "v.failed_ticket",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.snooze",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.snooze.a",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.snooze.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.canceled",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.canceled.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.last_comment.m",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.last_comment.a",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.last_comment.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "tags",
		Type: pattern.FieldTypeStringArray,
	},
	{
		Name: "v.activation_date",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.activation_date",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "v.long_output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.initial_output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.initial_long_output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.state.initiator",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "v.change_state",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.total_state_changes",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "v.meta",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "v.meta",
		Type: pattern.FieldTypeString,
	},
}

var entityFields = []FieldResponse{
	{
		Name: "_id",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "name",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "type",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "component",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "connector",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "infos",
		Type: pattern.FieldTypeObject,
	},
	{
		Name: "component_infos",
		Type: pattern.FieldTypeObject,
	},
	{
		Name: "category",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "impact_level",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "last_event_date",
		Type: pattern.FieldTypeTimestamp,
	},
	{
		Name: "pause",
		Type: pattern.FieldTypeReference,
	},
	{
		Name: "pause.author",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "pause.initiator",
		Type: pattern.FieldTypeString,
	},
}

var eventFields = []FieldResponse{
	{
		Name: "event_type",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "state",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "source_type",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "component",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "connector",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "connector_name",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "resource",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "extra",
		Type: pattern.FieldTypeObject,
	},
	{
		Name: "long_output",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "author",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "initiator",
		Type: pattern.FieldTypeString,
	},
}

var pbehaviorFields = []FieldResponse{
	{
		Name: "pbehavior_info.id",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "pbehavior_info.reason",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "pbehavior_info.type",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "pbehavior_info.canonical_type",
		Type: pattern.FieldTypeString,
	},
}

var weatherServiceFields = []FieldResponse{
	{
		Name: "is_grey",
		Type: pattern.FieldTypeBool,
	},
	{
		Name: "icon",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "secondary_icon",
		Type: pattern.FieldTypeString,
	},
	{
		Name: "state.val",
		Type: pattern.FieldTypeInt,
	},
}

type FieldsRequest struct {
	// Type should be only used to get state settings inherited pattern fields
	// if empty inherited pattern fields aren't included in the response.
	Type string `form:"type" binding:"oneoforempty=service component"`
}

type FieldGetter interface {
	Get(ctx context.Context, collection string, r FieldsRequest) (FieldsResponse, error)
}

type FieldsResponse struct {
	AlarmPattern           []AlarmFieldResponse  `json:"alarm_pattern,omitempty"`
	EntityPattern          []EntityFieldResponse `json:"entity_pattern,omitempty"`
	InheritedEntityPattern []EntityFieldResponse `json:"inherited_entity_pattern,omitempty"`
	EventPattern           []FieldResponse       `json:"event_pattern,omitempty"`
	PbehaviorPattern       []FieldResponse       `json:"pbehavior_pattern,omitempty"`
	WeatherServicePattern  []FieldResponse       `json:"weather_service_pattern,omitempty"`
}

type FieldResponse struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

type AlarmFieldResponse struct {
	FieldResponse
	OnlyAbsoluteTimeCond *bool `json:"only_absolute_time_cond,omitempty"`
}

type EntityFieldResponse struct {
	FieldResponse
	Alias bool `json:"alias"`
}

func NewFieldGetter(dbClient mongo.DbClient) FieldGetter {
	return &fieldGetter{
		entityPropCollection:                     dbClient.Collection(mongo.EntityInfosPropertyCollection),
		disabledFieldsInAlarmPattern:             patternfields.GetDisabledFieldsInAlarmPattern(),
		onlyAbsoluteTimeCondFieldsInAlarmPattern: patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(),
		disabledFieldsInEntityPattern:            patternfields.GetDisabledFieldsInEntityPattern(),
		disabledFieldsInInheritedEntityPattern:   patternfields.GetDisabledFieldsInInheritedEntityPattern(),
		disabledFieldsInEventPattern:             patternfields.GetDisabledFieldsInEventPattern(),
		disabledFieldsInPbehaviorPattern:         patternfields.GetDisabledFieldsInPbehaviorPattern(),
		disabledFieldsInWeatherServicePattern:    patternfields.GetDisabledFieldsInWeatherServicePattern(),
	}
}

type fieldGetter struct {
	entityPropCollection                     mongo.DbCollection
	disabledFieldsInAlarmPattern             map[string]map[string]bool
	onlyAbsoluteTimeCondFieldsInAlarmPattern map[string]map[string]bool
	disabledFieldsInEntityPattern            map[string]map[string]bool
	disabledFieldsInInheritedEntityPattern   map[string]map[string]map[string]bool
	disabledFieldsInEventPattern             map[string]map[string]bool
	disabledFieldsInPbehaviorPattern         map[string]map[string]bool
	disabledFieldsInWeatherServicePattern    map[string]map[string]bool
}

func (g *fieldGetter) Get(ctx context.Context, collection string, r FieldsRequest) (FieldsResponse, error) {
	var ep, iep []EntityFieldResponse
	var err error

	if disabledFields, ok := g.disabledFieldsInEntityPattern[collection]; ok {
		ep, err = g.getEntityFields(ctx, disabledFields)
		if err != nil {
			return FieldsResponse{}, err
		}
	}

	if disabledFieldsByType, ok := g.disabledFieldsInInheritedEntityPattern[collection]; ok {
		if disabledFields, ok := disabledFieldsByType[r.Type]; ok {
			iep, err = g.getEntityFields(ctx, disabledFields)
			if err != nil {
				return FieldsResponse{}, err
			}
		}
	}

	return FieldsResponse{
		AlarmPattern:           g.getAlarmFields(collection),
		EntityPattern:          ep,
		InheritedEntityPattern: iep,
		EventPattern:           g.getFields(collection, eventFields, g.disabledFieldsInEventPattern),
		PbehaviorPattern:       g.getFields(collection, pbehaviorFields, g.disabledFieldsInPbehaviorPattern),
		WeatherServicePattern:  g.getFields(collection, weatherServiceFields, g.disabledFieldsInWeatherServicePattern),
	}, nil
}

func (g *fieldGetter) getAlarmFields(collection string) []AlarmFieldResponse {
	disabledFields, ok := g.disabledFieldsInAlarmPattern[collection]
	if !ok {
		return nil
	}

	fields := make([]AlarmFieldResponse, len(alarmFields))
	for i := range fields {
		fields[i].FieldResponse = alarmFields[i]
		fields[i].Enabled = !disabledFields[fields[i].Name]
		if fields[i].Type == pattern.FieldTypeTimestamp {
			b := g.onlyAbsoluteTimeCondFieldsInAlarmPattern[collection][fields[i].Name]
			fields[i].OnlyAbsoluteTimeCond = &b
		}
	}

	return fields
}

func (g *fieldGetter) getEntityFields(ctx context.Context, disabledFields map[string]bool) ([]EntityFieldResponse, error) {
	fields := make([]EntityFieldResponse, len(entityFields))
	for i := range fields {
		fields[i].FieldResponse = entityFields[i]
		fields[i].Enabled = !disabledFields[fields[i].Name]
	}

	cursor, err := g.entityPropCollection.Find(ctx,
		bson.M{"alias": bson.M{"$nin": bson.A{nil, ""}}},
		options.Find().
			SetProjection(bson.M{"alias": 1, "type": 1}).
			SetSort(bson.M{"created": 1}).
			SetLimit(aliasLimit),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot find aliases: %w", err)
	}

	defer cursor.Close(ctx)
	typeMapping := map[int]string{
		types.EntityInfoTypeBoolean:     pattern.FieldTypeBool,
		types.EntityInfoTypeNumber:      pattern.FieldTypeInt,
		types.EntityInfoTypeTimestamp:   pattern.FieldTypeTimestamp,
		types.EntityInfoTypeString:      pattern.FieldTypeString,
		types.EntityInfoTypeStringArray: pattern.FieldTypeStringArray,
	}
	for cursor.Next(ctx) {
		alias := struct {
			Alias string `bson:"alias"`
			Type  int    `bson:"type"`
		}{}
		err = cursor.Decode(&alias)
		if err != nil {
			return nil, fmt.Errorf("cannot decode alias: %w", err)
		}

		fields = append(fields, EntityFieldResponse{
			Name:    alias.Alias,
			Type:    cmp.Or(typeMapping[alias.Type], pattern.FieldTypeString),
			Enabled: true,
			Alias:   true,
		})
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cannot fetch aliases: %w", err)
	}

	return fields, nil
}

func (g *fieldGetter) getFields(collection string, fields []FieldResponse, disabledFieldsByCollection map[string]map[string]bool) []FieldResponse {
	disabledFields, ok := disabledFieldsByCollection[collection]
	if !ok {
		return nil
	}

	res := make([]FieldResponse, len(fields))
	copy(res, fields)
	for i := range res {
		res[i].Enabled = !disabledFields[fields[i].Name]
	}

	return res
}
