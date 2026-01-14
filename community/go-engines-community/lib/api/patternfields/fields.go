package patternfields

import (
	"cmp"
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	FieldTypeObject    = "object"
	FieldTypeDuration  = "duration"
	FieldTypeReference = "reference"
)

const aliasLimit = 500

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
		Type: FieldTypeDuration,
	},
	{
		Name: "v.infos",
		Type: FieldTypeObject,
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
		Type: FieldTypeReference,
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
		Type: FieldTypeReference,
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
		Name: "v.ticket.ticket_data",
		Type: FieldTypeObject,
	},
	{
		Name: "v.snooze",
		Type: FieldTypeReference,
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
		Type: FieldTypeReference,
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
		Type: FieldTypeReference,
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
		Type: FieldTypeReference,
	},
	{
		Name: "v.total_state_changes",
		Type: pattern.FieldTypeInt,
	},
	{
		Name: "v.meta",
		Type: FieldTypeReference,
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
		Type: FieldTypeObject,
	},
	{
		Name: "component_infos",
		Type: FieldTypeObject,
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
		Type: FieldTypeObject,
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

var disabledFieldsInAlarmPattern = map[string][]string{
	mongo.IdleRuleMongoCollection:          {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.MetaAlarmRulesMongoCollection:    {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.FlappingRuleMongoCollection:      {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.ResolveRuleMongoCollection:       {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.ScenarioCollection:               {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.InstructionMongoCollection:       {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.DeclareTicketRuleCollection:      {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.LinkRuleMongoCollection:          {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.DynamicInfosRulesMongoCollection: {"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "v.infos"},
	mongo.AlarmTagCollection:               {"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "tags"},
	mongo.WidgetFiltersMongoCollection:     {},
}

var onlyAbsoluteTimeCondFieldsInAlarmPattern = map[string][]string{
	mongo.IdleRuleMongoCollection:          {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.DynamicInfosRulesMongoCollection: {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.MetaAlarmRulesMongoCollection:    {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.FlappingRuleMongoCollection:      {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.ResolveRuleMongoCollection:       {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.ScenarioCollection:               {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.InstructionMongoCollection:       {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.DeclareTicketRuleCollection:      {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.LinkRuleMongoCollection:          {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.AlarmTagCollection:               {"v.creation_date", "v.ack.t", "v.activation_date"},
}

var disabledFieldsInEntityPattern = map[string][]string{
	mongo.StateSettingsMongoCollection:     {"last_event_date", "component", "component_infos"},
	mongo.EntityMongoCollection:            {"last_event_date", "connector", "component_infos"},
	mongo.PbehaviorMongoCollection:         {"last_event_date"},
	mongo.IdleRuleMongoCollection:          {"last_event_date"},
	mongo.DynamicInfosRulesMongoCollection: {"last_event_date"},
	mongo.MetaAlarmRulesMongoCollection:    {"last_event_date"},
	mongo.FlappingRuleMongoCollection:      {"last_event_date"},
	mongo.ResolveRuleMongoCollection:       {"last_event_date"},
	mongo.ScenarioCollection:               {"last_event_date"},
	mongo.InstructionMongoCollection:       {"last_event_date"},
	mongo.KpiFilterMongoCollection:         {"last_event_date"},
	mongo.DeclareTicketRuleCollection:      {"last_event_date"},
	mongo.LinkRuleMongoCollection:          {"last_event_date"},
	mongo.AlarmTagCollection:               {"last_event_date"},
	mongo.EventFilterRuleCollection:        {},
	mongo.WidgetFiltersMongoCollection:     {},
}

var disabledFieldsInEventPattern = map[string][]string{
	mongo.EventFilterRuleCollection: {},
}

var disabledFieldsInPbehaviorPattern = map[string][]string{
	mongo.WidgetFiltersMongoCollection: {},
}

var disabledFieldsInWeatherServicePattern = map[string][]string{
	mongo.WidgetFiltersMongoCollection: {},
}

type FieldGetter interface {
	Get(ctx context.Context, collection string) (FieldsResponse, error)
}

type FieldsResponse struct {
	AlarmPattern          []AlarmFieldResponse  `json:"alarm_pattern,omitempty"`
	EntityPattern         []EntityFieldResponse `json:"entity_pattern,omitempty"`
	EventPattern          []FieldResponse       `json:"event_pattern,omitempty"`
	PbehaviorPattern      []FieldResponse       `json:"pbehavior_pattern,omitempty"`
	WeatherServicePattern []FieldResponse       `json:"weather_service_pattern,omitempty"`
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
	g := &fieldGetter{
		entityPropCollection: dbClient.Collection(mongo.EntityInfosPropertyCollection),
	}

	g.diabledFieldsInAlarmPattern = make(map[string]map[string]bool, len(disabledFieldsInAlarmPattern))
	for k, v := range disabledFieldsInAlarmPattern {
		g.diabledFieldsInAlarmPattern[k] = g.sliceToMap(v)
	}

	g.onlyAbsoluteTimeCondFieldsInAlarmPattern = make(map[string]map[string]bool, len(onlyAbsoluteTimeCondFieldsInAlarmPattern))
	for k, v := range onlyAbsoluteTimeCondFieldsInAlarmPattern {
		g.onlyAbsoluteTimeCondFieldsInAlarmPattern[k] = g.sliceToMap(v)
	}

	g.disabledFieldsInEntityPattern = make(map[string]map[string]bool, len(disabledFieldsInEntityPattern))
	for k, v := range disabledFieldsInEntityPattern {
		g.disabledFieldsInEntityPattern[k] = g.sliceToMap(v)
	}

	g.disabledFieldsInEventPattern = make(map[string]map[string]bool, len(disabledFieldsInEventPattern))
	for k, v := range disabledFieldsInEventPattern {
		g.disabledFieldsInEventPattern[k] = g.sliceToMap(v)
	}

	g.disabledFieldsInPbehaviorPattern = make(map[string]map[string]bool, len(disabledFieldsInPbehaviorPattern))
	for k, v := range disabledFieldsInPbehaviorPattern {
		g.disabledFieldsInPbehaviorPattern[k] = g.sliceToMap(v)
	}

	g.disabledFieldsInWeatherServicePattern = make(map[string]map[string]bool, len(disabledFieldsInWeatherServicePattern))
	for k, v := range disabledFieldsInWeatherServicePattern {
		g.disabledFieldsInWeatherServicePattern[k] = g.sliceToMap(v)
	}

	return g
}

func GetForbiddenFieldsInAlarmPattern(collection string) []string {
	return disabledFieldsInAlarmPattern[collection]
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection string) []string {
	return onlyAbsoluteTimeCondFieldsInAlarmPattern[collection]
}

func GetForbiddenFieldsInEntityPattern(collection string) []string {
	return disabledFieldsInEntityPattern[collection]
}

type fieldGetter struct {
	entityPropCollection                     mongo.DbCollection
	diabledFieldsInAlarmPattern              map[string]map[string]bool
	onlyAbsoluteTimeCondFieldsInAlarmPattern map[string]map[string]bool
	disabledFieldsInEntityPattern            map[string]map[string]bool
	disabledFieldsInEventPattern             map[string]map[string]bool
	disabledFieldsInPbehaviorPattern         map[string]map[string]bool
	disabledFieldsInWeatherServicePattern    map[string]map[string]bool
}

func (g *fieldGetter) Get(ctx context.Context, collection string) (FieldsResponse, error) {
	ep, err := g.getEntityFields(ctx, collection)
	if err != nil {
		return FieldsResponse{}, err
	}

	return FieldsResponse{
		AlarmPattern:          g.getAlarmFields(collection),
		EntityPattern:         ep,
		EventPattern:          g.getFields(collection, eventFields, g.disabledFieldsInEventPattern),
		PbehaviorPattern:      g.getFields(collection, pbehaviorFields, g.disabledFieldsInPbehaviorPattern),
		WeatherServicePattern: g.getFields(collection, weatherServiceFields, g.disabledFieldsInWeatherServicePattern),
	}, nil
}

func (g *fieldGetter) getAlarmFields(collection string) []AlarmFieldResponse {
	disabledFields, ok := g.diabledFieldsInAlarmPattern[collection]
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

func (g *fieldGetter) getEntityFields(ctx context.Context, collection string) ([]EntityFieldResponse, error) {
	disabledFields, ok := g.disabledFieldsInEntityPattern[collection]
	if !ok {
		return nil, nil
	}

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
			FieldResponse: FieldResponse{
				Name:    alias.Alias,
				Type:    cmp.Or(typeMapping[alias.Type], pattern.FieldTypeString),
				Enabled: true,
			},
			Alias: true,
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

func (g *fieldGetter) sliceToMap(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}

	return m
}
