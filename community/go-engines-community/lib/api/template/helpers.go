package template

import (
	"context"
	"errors"
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
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
		Name:       "displayName",
		PluralName: "displayNames",
		Value:      "%var%.Value.DisplayName",
	},
	{
		Name:       "connector",
		PluralName: "connectors",
		Value:      "%var%.Value.Connector",
	},
	{
		Name:       "connectorName",
		PluralName: "connectorNames",
		Value:      "%var%.Value.ConnectorName",
	},
	{
		Name:       "component",
		PluralName: "components",
		Value:      "%var%.Value.Component",
	},
	{
		Name:       "resource",
		PluralName: "resources",
		Value:      "%var%.Value.Resource",
	},
	{
		Name:       "output",
		PluralName: "outputs",
		Value:      "%var%.Value.Output",
	},
	{
		Name:       "initialOutput",
		PluralName: "initialOutputs",
		Value:      "%var%.Value.InitialOutput",
	},
	{
		Name:       "stateMessage",
		PluralName: "stateMessages",
		Value:      "%var%.Value.State.Message",
	},
	{
		Name:       "stateValue",
		PluralName: "stateValues",
		Value:      "%var%.Value.State.Value",
	},
	{
		Name:       "statusValue",
		PluralName: "statusValues",
		Value:      "%var%.Value.Status.Value",
	},
	{
		Name:       "ticketAuthor",
		PluralName: "ticketAuthors",
		Value:      "%var%.Value.Ticket.Author",
	},
	{
		Name:       "ticketID",
		PluralName: "ticketIDs",
		Value:      "%var%.Value.Ticket.Ticket",
	},
	{
		Name:       "ticketMessage",
		PluralName: "ticketMessages",
		Value:      "%var%.Value.Ticket.Message",
	},
	{
		Name:       "ackAuthor",
		PluralName: "ackAuthors",
		Value:      "%var%.Value.ACK.Author",
	},
	{
		Name:       "ackMessage",
		PluralName: "ackMessages",
		Value:      "%var%.Value.ACK.Message",
	},
	{
		Name:       "lastCommentAuthor",
		PluralName: "lastCommentAuthors",
		Value:      "%var%.Value.LastComment.Author",
	},
	{
		Name:       "lastCommentMessage",
		PluralName: "lastCommentMessages",
		Value:      "%var%.Value.LastComment.Message",
	},
	{
		Name:       "infos",
		PluralName: "infos",
		Value:      "(index (index %var%.Value.Infos \"%rule_id%\") \"%infos_name%\")",
	},
}

var entityVars = []tplVar{
	{
		Name:       "name",
		PluralName: "names",
		Value:      "%var%.Name",
	},
	{
		Name:       "type",
		PluralName: "types",
		Value:      "%var%.Type",
	},
	{
		Name:       "infos",
		PluralName: "infos",
		Value:      "(index %var%.Infos \"%infos_name%\").Value",
	},
	{
		Name:       "componentInfos",
		PluralName: "componentInfos",
		Value:      "(index %var%.ComponentInfos \"%infos_name%\").Value",
	},
}

var eventVars = []tplVar{
	{
		Name:       "connector",
		PluralName: "connectors",
		Value:      "%var%.Connector",
	},
	{
		Name:       "connectorName",
		PluralName: "connectorNames",
		Value:      "%var%.ConnectorName",
	},
	{
		Name:       "component",
		PluralName: "components",
		Value:      "%var%.Component",
	},
	{
		Name:       "resource",
		PluralName: "resources",
		Value:      "%var%.Resource",
	},
	{
		Name:       "output",
		PluralName: "outputs",
		Value:      "%var%.Output",
	},
	{
		Name:       "extraInfos",
		PluralName: "extraInfos",
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
		Name:  "environmentVariables",
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
		lv, _ := l.Value.(string)
		rv, _ := r.Value.(string)

		return strings.Compare(lv, rv)
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
	err := collection.FindOne(ctx, bson.M{"_id": id, "type": TypeTestDataEvent}, options.FindOne().SetProjection(bson.M{"body": 1})).Decode(&m)
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

	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "type": TypeTestDataResponse},
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

func GetAlarmDataFromTest(ctx context.Context, collection mongo.DbCollection, testID string, ruleType int, ruleID string) (types.AlarmWithEntity, error) {
	var alarm types.AlarmWithEntity
	if testID == "" {
		return alarm, nil
	}

	test := TestModel{}
	err := collection.FindOne(ctx, bson.M{"_id": testID, "type": ruleType, "rule._id": ruleID}).Decode(&test)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return alarm, common.NewValidationError("testdata.test", "Test doesn't exist.")
		}

		return alarm, err
	}

	if test.Data.Alarm != nil {
		alarm.Alarm = test.Data.Alarm.Alarm
		alarm.Entity = test.Data.Alarm.Entity
	}

	return alarm, nil
}
