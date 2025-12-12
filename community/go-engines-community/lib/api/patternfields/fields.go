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

type FieldGetter interface {
	Get(ctx context.Context, collection string) (FieldsResponse, error)
}

type FieldsResponse struct {
	AlarmPattern  []AlarmFieldResponse  `json:"alarm_pattern,omitempty"`
	EntityPattern []EntityFieldResponse `json:"entity_pattern,omitempty"`
}

type AlarmFieldResponse struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Enabled              bool   `json:"enabled"`
	OnlyAbsoluteTimeCond *bool  `json:"only_absolute_time_cond,omitempty"`
}

type EntityFieldResponse struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
	Alias   bool   `json:"alias"`
}

func NewFieldGetter(dbClient mongo.DbClient) FieldGetter {
	return &fieldGetter{
		entityPropCollection: dbClient.Collection(mongo.EntityInfosPropertyCollection),
		alarmFields: []AlarmFieldResponse{
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
				Name: "v.tags",
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
		},
		entityFields: []EntityFieldResponse{
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
				Name: "last_event_data",
				Type: pattern.FieldTypeTimestamp,
			},
		},
	}
}

func GetForbiddenFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved"}
	case mongo.DynamicInfosRulesMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "v.infos"}
	case mongo.AlarmTagCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "tags"}
	default:
		return nil
	}
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"v.creation_date", "v.ack.t", "v.activation_date"}
	default:
		return nil
	}
}

func GetForbiddenFieldsInEntityPattern(collection string) []string {
	switch collection {
	case mongo.StateSettingsMongoCollection:
		return []string{"last_event_date", "component", "component_infos"}
	case mongo.EntityMongoCollection:
		return []string{"last_event_date", "connector", "component_infos"}
	case mongo.PbehaviorMongoCollection,
		mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.KpiFilterMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"last_event_date"}
	default:
		return nil
	}
}

type fieldGetter struct {
	alarmFields          []AlarmFieldResponse
	entityFields         []EntityFieldResponse
	entityPropCollection mongo.DbCollection
}

func (g *fieldGetter) Get(ctx context.Context, collection string) (FieldsResponse, error) {
	ep, err := g.getEntityFields(ctx, collection)
	if err != nil {
		return FieldsResponse{}, err
	}

	return FieldsResponse{
		AlarmPattern:  g.getAlarmFields(collection),
		EntityPattern: ep,
	}, nil
}

func (g *fieldGetter) getAlarmFields(collection string) []AlarmFieldResponse {
	fields := make([]AlarmFieldResponse, len(g.alarmFields))
	copy(fields, g.alarmFields)
	forbidden := g.sliceToMap(GetForbiddenFieldsInAlarmPattern(collection))
	absoluteTime := g.sliceToMap(GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection))
	for i := range fields {
		fields[i].Enabled = !forbidden[fields[i].Name]
		if fields[i].Type == pattern.FieldTypeTimestamp {
			b := absoluteTime[fields[i].Name]
			fields[i].OnlyAbsoluteTimeCond = &b
		}
	}

	return fields
}

func (g *fieldGetter) getEntityFields(ctx context.Context, collection string) ([]EntityFieldResponse, error) {
	fields := make([]EntityFieldResponse, len(g.entityFields))
	copy(fields, g.entityFields)
	forbidden := g.sliceToMap(GetForbiddenFieldsInEntityPattern(collection))
	for i := range fields {
		fields[i].Enabled = !forbidden[fields[i].Name]
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

func (g *fieldGetter) sliceToMap(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}

	return m
}
