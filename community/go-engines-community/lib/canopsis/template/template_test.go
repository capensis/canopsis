package template

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestFunctions(t *testing.T) {
	dataSets := map[string][]struct {
		Tpl         string
		TplData     any
		Tz          string
		ExpectedRes string
		ExpectedErr error
	}{
		"split": {
			{
				Tpl: `{{ split .Sep .Index .Input }}`,
				TplData: map[string]any{
					"Sep":   ",",
					"Index": 0,
					"Input": "NgocHa,MinhNghia,Minh",
				},
				ExpectedRes: "NgocHa",
			},
			{
				Tpl: `{{ split .Sep .Index .Input }}`,
				TplData: map[string]any{
					"Sep":   ",",
					"Index": -1,
					"Input": "NgocHa,MinhNghia,Minh",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{ split .Sep .Index .Input }}`,
				TplData: map[string]any{
					"Sep":   " ",
					"Index": 1,
					"Input": "This is space",
				},
				ExpectedRes: "is",
			},
		},
		"trim": {
			{
				Tpl:         `{{ trim . }}`,
				TplData:     "  ",
				ExpectedRes: "",
			},
			{
				Tpl:         `{{ trim . }}`,
				TplData:     " kratos ",
				ExpectedRes: "kratos",
			},
			{
				Tpl:         `{{ trim . }}`,
				TplData:     "\tkratos\n ",
				ExpectedRes: "kratos",
			},
		},
		"uppercase": {
			{
				Tpl:         `{{ uppercase . }}`,
				TplData:     "  ",
				ExpectedRes: "  ",
			},
			{
				Tpl:         `{{ uppercase . }}`,
				TplData:     "kratos",
				ExpectedRes: "KRATOS",
			},
			{
				Tpl:         `{{ uppercase . }}`,
				TplData:     "KraTos",
				ExpectedRes: "KRATOS",
			},
		},
		"lowercase": {
			{
				Tpl:         `{{ lowercase . }}`,
				TplData:     "  ",
				ExpectedRes: "  ",
			},
			{
				Tpl:         `{{ lowercase . }}`,
				TplData:     "kratos",
				ExpectedRes: "kratos",
			},
			{
				Tpl:         `{{ lowercase . }}`,
				TplData:     "KraTos",
				ExpectedRes: "kratos",
			},
			{
				Tpl:         `{{ lowercase . }}`,
				TplData:     "KRATOS",
				ExpectedRes: "kratos",
			},
		},
		"localtime": {
			{
				Tpl: `{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" "Australia/Queensland" }}`,
				TplData: map[string]any{
					"TestDate": datetime.CpsTime{
						Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
					},
				},
				ExpectedRes: "Thu, 28 Oct 2021 17:05:00 AEST",
			},
			{
				Tpl: `{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" }}`,
				TplData: map[string]any{
					"TestDate": datetime.CpsTime{
						Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
					},
				},
				Tz:          "Australia/Queensland",
				ExpectedRes: "Thu, 28 Oct 2021 17:05:00 AEST",
			},
		},
		"tag_has_key": {
			{
				Tpl: `{{tag_has_key .Alarm.Tags "Tag1" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "true",
			},
			{
				Tpl: `{{tag_has_key .Alarm.Tags "Tag2" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "true",
			},
			{
				Tpl: `{{tag_has_key .Alarm.Tags "Tag3" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "false",
			},
		},
		"get_tag": {
			{
				Tpl: `{{get_tag .Alarm.Tags "Tag1" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "Value1",
			},
			{
				Tpl: `{{get_tag .Alarm.Tags "Tag2" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{get_tag .Alarm.Tags "Tag3" }}`,
				TplData: map[string]any{
					"Alarm": types.Alarm{
						Tags: []string{"Tag1: Value1", "Tag2"},
					},
				},
				ExpectedRes: "",
			},
		},
		"substrLeft": {
			{
				Tpl: `{{substrLeft .String 5 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello",
			},
			{
				Tpl: `{{substrLeft .String 0 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substrLeft .String 13 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello, world!",
			},
			{
				Tpl: `{{substrLeft .String 100000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello, world!",
			},
			{
				Tpl: `{{substrLeft .String -100000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substrLeft .String 8 }}`,
				TplData: map[string]any{
					"String": "élémentaire",
				},
				ExpectedRes: "élémenta",
			},
		},
		"substrRight": {
			{
				Tpl: `{{substrRight .String 6 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "world!",
			},
			{
				Tpl: `{{substrRight .String 0 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substrRight .String 13 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello, world!",
			},
			{
				Tpl: `{{substrRight .String 100000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello, world!",
			},
			{
				Tpl: `{{substrRight .String -100000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substrRight .String 6 }}`,
				TplData: map[string]any{
					"String": "élémenta",
				},
				ExpectedRes: "émenta",
			},
		},
		"substr": {
			{
				Tpl: `{{substr .String 0 5 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello",
			},
			{
				Tpl: `{{substr .String 7 5 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "world",
			},
			{
				Tpl: `{{substr .String 7 1000000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "world!",
			},
			{
				Tpl: `{{substr .String 0 1000000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "Hello, world!",
			},
			{
				Tpl: `{{substr .String -1 1000000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substr .String 5 0 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substr .String 0 -1000000 }}`,
				TplData: map[string]any{
					"String": "Hello, world!",
				},
				ExpectedRes: "",
			},
			{
				Tpl: `{{substr .String 1 2 }}`,
				TplData: map[string]any{
					"String": "élémenta",
				},
				ExpectedRes: "lé",
			},
		},
		"add": {
			{
				Tpl: `{{ add .A .B }}`,
				TplData: map[string]any{
					"A": 2,
					"B": 2,
				},
				ExpectedRes: "4",
			},
			{
				Tpl: `{{ add .A .B }}`,
				TplData: map[string]any{
					"A": "2.5",
					"B": 2.5,
				},
				ExpectedErr: ErrFailedConvertToInt64,
			},
		},
		"sub": {
			{
				Tpl: `{{ sub .A .B }}`,
				TplData: map[string]any{
					"A": 4,
					"B": 2,
				},
				ExpectedRes: "2",
			},
			{
				Tpl: `{{ sub .A .B }}`,
				TplData: map[string]any{
					"A": "5",
					"B": 2.5,
				},
				ExpectedErr: ErrFailedConvertToInt64,
			},
		},
		"mult": {
			{
				Tpl: `{{ mult .A .B }}`,
				TplData: map[string]any{
					"A": 2,
					"B": 2,
				},
				ExpectedRes: "4",
			},
			{
				Tpl: `{{ mult .A .B }}`,
				TplData: map[string]any{
					"A": "2.5",
					"B": 2.5,
				},
				ExpectedErr: ErrFailedConvertToInt64,
			},
		},
		"div": {
			{
				Tpl: `{{ div .A .B }}`,
				TplData: map[string]any{
					"A": 4,
					"B": 2,
				},
				ExpectedRes: "2",
			},
			{
				Tpl: `{{ div .A .B }}`,
				TplData: map[string]any{
					"A": 13,
					"B": 2,
				},
				ExpectedRes: "6",
			},
			{
				Tpl: `{{ div .A .B }}`,
				TplData: map[string]any{
					"A": 6,
					"B": 0,
				},
				ExpectedErr: ErrDivisionByZero,
			},
			{
				Tpl: `{{ div .A .B }}`,
				TplData: map[string]any{
					"A": "6.25",
					"B": 2.5,
				},
				ExpectedErr: ErrFailedConvertToInt64,
			},
		},
		"substr_add_example": {
			{
				Tpl: `{{ $var := add .A .B }}{{ substrLeft .String $var }}`,
				TplData: map[string]any{
					"A":      1,
					"B":      2,
					"String": "qwerty",
				},
				ExpectedRes: "qwe",
			},
		},
	}

	for name, v := range dataSets {
		for i, data := range v {
			t.Run(name+"/"+strconv.Itoa(i), func(t *testing.T) {
				var loc *time.Location
				if data.Tz != "" {
					loc, _ = time.LoadLocation("Australia/Queensland")
				}

				tpl, err := template.New("test").
					Funcs(GetFunctions(loc)).
					Parse(data.Tpl)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				res, err := executeTemplate(tpl, data.TplData)
				if !errors.Is(err, data.ExpectedErr) {
					t.Errorf("expected err %v but got %v", data.ExpectedErr, err)
				}
				if res != data.ExpectedRes {
					t.Errorf("expected res %q but got %q", data.ExpectedRes, res)
				}
			})
		}
	}
}

func TestAddEnvVarsToData(t *testing.T) {
	alarm := types.Alarm{ID: "test-alarm"}
	envVars := map[string]any{
		"Location": "FR",
	}
	// the interface needed to test reflection in templates.
	type Activatable interface {
		IsActivated() bool
	}
	type activatableWithEnabled struct {
		Activatable
		Enabled bool
	}
	type alarmWithEnabled struct {
		types.Alarm
		Enabled bool
	}
	type alarmPtrWithEnabled struct {
		*types.Alarm
		Enabled bool
	}

	dataSet := []struct {
		Data        any
		ExpectedRes any
	}{
		{
			Data: map[string]types.Alarm{
				"Alarm": alarm,
			},
			ExpectedRes: map[string]any{
				"Alarm": alarm,
				"Env":   envVars,
			},
		},
		{
			Data: alarm,
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: &alarm,
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: Activatable(&alarm),
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: alarmWithEnabled{
				Alarm:   alarm,
				Enabled: true,
			},
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"Enabled":                           true,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: alarmPtrWithEnabled{
				Alarm:   &alarm,
				Enabled: true,
			},
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"Enabled":                           true,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: activatableWithEnabled{
				Activatable: Activatable(&alarm),
				Enabled:     true,
			},
			ExpectedRes: map[string]any{
				"EntityID":                          alarm.EntityID,
				"Enabled":                           true,
				"ID":                                alarm.ID,
				"KpiAssignedInstructions":           alarm.KpiAssignedInstructions,
				"KpiExecutedInstructions":           alarm.KpiExecutedInstructions,
				"KpiAssignedAutoInstructions":       alarm.KpiAssignedAutoInstructions,
				"KpiExecutedAutoInstructions":       alarm.KpiExecutedAutoInstructions,
				"Tags":                              alarm.Tags,
				"InternalTags":                      alarm.InternalTags,
				"InternalTagsUpdated":               datetime.MicroTime{},
				"ExternalTags":                      alarm.ExternalTags,
				"Time":                              alarm.Time,
				"Value":                             alarm.Value,
				"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
				"NotAckedMetricType":                alarm.NotAckedMetricType,
				"NotAckedSince":                     alarm.NotAckedSince,
				"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
				"Healthcheck":                       alarm.Healthcheck,
				"Env":                               envVars,
			},
		},
		{
			Data: map[int]types.Alarm{
				1: alarm,
			},
			ExpectedRes: map[int]types.Alarm{
				1: alarm,
			},
		},
		{
			Data:        []types.Alarm{alarm},
			ExpectedRes: []types.Alarm{alarm},
		},
		{
			Data:        1,
			ExpectedRes: 1,
		},
	}

	for i, data := range dataSet {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := addEnvVarsToData(data.Data, envVars)
			if diff := pretty.Compare(res, data.ExpectedRes); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func executeTemplate(tmpl *template.Template, payload interface{}) (string, error) {
	var b bytes.Buffer
	err := tmpl.Execute(io.Writer(&b), payload)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
