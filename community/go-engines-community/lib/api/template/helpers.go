package template

import (
	"context"
	"errors"
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var alarmVars = []tplVar{
	{
		Name:       "Display name",
		PluralName: "Display names",
		Value:      "%var%.Value.DisplayName",
	},
	{
		Name:       "Connector",
		PluralName: "Connectors",
		Value:      "%var%.Value.Connector",
	},
	{
		Name:       "Connector name",
		PluralName: "Connector names",
		Value:      "%var%.Value.ConnectorName",
	},
	{
		Name:       "Component",
		PluralName: "Components",
		Value:      "%var%.Value.Component",
	},
	{
		Name:       "Resource",
		PluralName: "Resources",
		Value:      "%var%.Value.Resource",
	},
	{
		Name:       "Output",
		PluralName: "Outputs",
		Value:      "%var%.Value.Output",
	},
	{
		Name:       "Initial output",
		PluralName: "Initial outputs",
		Value:      "%var%.Value.InitialOutput",
	},
	{
		Name:       "State message",
		PluralName: "State messages",
		Value:      "%var%.Value.State.Message",
	},
	{
		Name:       "State value",
		PluralName: "State values",
		Value:      "%var%.Value.State.Value",
	},
	{
		Name:       "Status value",
		PluralName: "Status values",
		Value:      "%var%.Value.Status.Value",
	},
	{
		Name:       "Ticket author",
		PluralName: "Ticket authors",
		Value:      "%var%.Value.Ticket.Author",
	},
	{
		Name:       "Ticket id",
		PluralName: "Ticket ids",
		Value:      "%var%.Value.Ticket.Ticket",
	},
	{
		Name:       "Ticket message",
		PluralName: "Ticket messages",
		Value:      "%var%.Value.Ticket.Message",
	},
	{
		Name:       "Ack author",
		PluralName: "Ack authors",
		Value:      "%var%.Value.ACK.Author",
	},
	{
		Name:       "Ack message",
		PluralName: "Ack messages",
		Value:      "%var%.Value.ACK.Message",
	},
	{
		Name:       "Last comment author",
		PluralName: "Last comment authors",
		Value:      "%var%.Value.LastComment.Author",
	},
	{
		Name:       "Last comment message",
		PluralName: "Last comment messages",
		Value:      "%var%.Value.LastComment.Message",
	},
	{
		Name:       "Infos",
		PluralName: "Infos",
		Value:      "(index (index %var%.Value.Infos \"%rule_id%\") \"%infos_name%\")",
	},
}

var entityVars = []tplVar{
	{
		Name:       "Name",
		PluralName: "Names",
		Value:      "%var%.Name",
	},
	{
		Name:       "Type",
		PluralName: "Types",
		Value:      "%var%.Type",
	},
	{
		Name:       "Infos",
		PluralName: "Infos",
		Value:      "(index %var%.Infos \"%infos_name%\").Value",
	},
	{
		Name:       "Component infos",
		PluralName: "Component infos",
		Value:      "(index %var%.ComponentInfos \"%infos_name%\").Value",
	},
}

var eventVars = []tplVar{
	{
		Name:       "Connector",
		PluralName: "Connectors",
		Value:      "%var%.Connector",
	},
	{
		Name:       "Connector name",
		PluralName: "Connector names",
		Value:      "%var%.ConnectorName",
	},
	{
		Name:       "Component",
		PluralName: "Components",
		Value:      "%var%.Component",
	},
	{
		Name:       "Resource",
		PluralName: "Resources",
		Value:      "%var%.Resource",
	},
	{
		Name:       "Output",
		PluralName: "Outputs",
		Value:      "%var%.Output",
	},
	{
		Name:       "Extra infos",
		PluralName: "Extra infos",
		Value:      "index %var%.ExtraInfos \"%infos_name%\"",
	},
}

type ResponseTestData struct {
	Body    map[string]any    `bson:"body"`
	Headers map[string]string `bson:"headers"`
}

type tplVar struct {
	Name       string
	PluralName string
	Value      string
}

func Validate(tplValidator validator.Validator, str string, data any) (ValidateResponse, error) {
	isValid, errReport, res, err := tplValidator.Validate(str, data)
	if err != nil {
		return ValidateResponse{}, err
	}

	return ValidateResponse{
		IsValid: isValid,
		Err:     errReport,
		Result:  res,
	}, nil
}

func AddEnvVars(vars []VarResponse, tplConfigProvider config.TemplateConfigProvider) []VarResponse {
	envVars := GetEnvVars(tplConfigProvider)
	res := make([]VarResponse, 0, len(vars)+1)
	res = append(res, vars...)
	res = append(res, VarResponse{
		Name:  "Environment variables",
		Value: envVars,
	})

	return res
}

func GetEnvVars(tplConfigProvider config.TemplateConfigProvider) []VarResponse {
	vars := tplConfigProvider.Get().Vars
	keys := make([]VarResponse, 0, len(vars))
	for k, v := range vars {
		if k == config.SystemEnvVariablesKey {
			if m, ok := v.(map[string]string); ok {
				for sk := range m {
					keys = append(keys, VarResponse{Value: "{{ ." + libtemplate.EnvVar + "." + k + "." + sk + " }}"})
				}

				continue
			}
		}

		keys = append(keys, VarResponse{Value: "{{ ." + libtemplate.EnvVar + "." + k + " }}"})
	}

	slices.SortFunc(keys, func(l, r VarResponse) int {
		return strings.Compare(l.Value.(string), r.Value.(string))
	})

	return keys
}

func GetAlarmVars(prefixVar, suffixVar, alarmVar string, pluralName bool) []VarResponse {
	res := make([]VarResponse, len(alarmVars))
	for i, v := range alarmVars {
		name := ""
		if pluralName {
			name = v.PluralName
		} else {
			name = v.Name
		}

		res[i] = VarResponse{
			Name:  name,
			Value: prefixVar + strings.Replace(v.Value, "%var%", alarmVar, 1) + suffixVar,
		}
	}

	return res
}

func GetEntityVars(prefixVar, suffixVar, entityVar string, pluralName bool) []VarResponse {
	res := make([]VarResponse, len(entityVars))
	for i, v := range entityVars {
		name := ""
		if pluralName {
			name = v.PluralName
		} else {
			name = v.Name
		}

		res[i] = VarResponse{
			Name:  name,
			Value: prefixVar + strings.Replace(v.Value, "%var%", entityVar, 1) + suffixVar,
		}
	}

	return res
}

func GetEventVars(prefixVar, suffixVar, eventVar string, pluralName bool) []VarResponse {
	res := make([]VarResponse, len(eventVars))
	for i, v := range eventVars {
		name := ""
		if pluralName {
			name = v.PluralName
		} else {
			name = v.Name
		}

		res[i] = VarResponse{
			Name:  name,
			Value: prefixVar + strings.Replace(v.Value, "%var%", eventVar, 1) + suffixVar,
		}
	}

	return res
}

func GetEventData(ctx context.Context, collection mongo.DbCollection, id string, encoder encoding.Encoder, decoder encoding.Decoder) (*types.Event, error) {
	if id == "" {
		return nil, nil
	}

	m := DataResponse{}
	err := collection.FindOne(ctx, bson.M{"_id": id, "type": TypeEvent}, options.FindOne().SetProjection(bson.M{"body": 1})).Decode(&m)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	b, err := encoder.Encode(m.Body)
	if err != nil {
		return nil, err
	}

	event := types.Event{}
	err = decoder.Decode(b, &event)
	if err != nil {
		return nil, err
	}

	err = event.InjectExtraInfos(b)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func GetResponseData(ctx context.Context, collection mongo.DbCollection, idsByIndex map[int]string) (map[int]ResponseTestData, error) {
	if len(idsByIndex) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(idsByIndex))
	idxes := make(map[string][]int, len(ids))
	for i, id := range idsByIndex {
		ids = append(ids, id)
		idxes[id] = append(idxes[id], i)
	}

	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "type": TypeResponse},
		options.Find().SetProjection(bson.M{"body": 1, "headers": 1}))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := make(map[int]ResponseTestData, len(ids))
	for cursor.Next(ctx) {
		m := struct {
			ID               string `bson:"_id"`
			ResponseTestData `bson:",inline"`
		}{}
		err = cursor.Decode(&m)
		if err != nil {
			return nil, err
		}

		for _, i := range idxes[m.ID] {
			res[i] = m.ResponseTestData
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	if len(res) != len(ids) {
		return nil, nil
	}

	return res, nil
}
