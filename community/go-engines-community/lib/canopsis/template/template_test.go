package template

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestFunctions(t *testing.T) {
	f := func(tplStr string, tplData any, expectedRes string, expectedErr error, tz string) {
		t.Helper()
		var loc *time.Location
		if tz != "" {
			loc, _ = time.LoadLocation("Australia/Queensland")
		}

		tpl, err := template.New("test").
			Funcs(GetFunctions(loc)).
			Parse(tplStr)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		res, err := executeTemplate(tpl, tplData)
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected err %v but got %v", expectedErr, err)
		}

		if res != expectedRes {
			t.Fatalf("expected res %q but got %q", expectedRes, res)
		}
	}

	f(
		`{{ split .Sep .Index .Input }}`,
		map[string]any{
			"Sep":   ",",
			"Index": 0,
			"Input": "NgocHa,MinhNghia,Minh",
		},
		"NgocHa",
		nil,
		"",
	)
	f(
		`{{ split .Sep .Index .Input }}`,
		map[string]any{
			"Sep":   ",",
			"Index": -1,
			"Input": "NgocHa,MinhNghia,Minh",
		},
		"",
		nil,
		"",
	)
	f(
		`{{ split .Sep .Index .Input }}`,
		map[string]any{
			"Sep":   " ",
			"Index": 1,
			"Input": "This is space",
		},
		"is",
		nil,
		"",
	)
	f(
		`{{ trim . }}`,
		"  ",
		"",
		nil,
		"",
	)
	f(
		`{{ trim . }}`,
		" kratos ",
		"kratos",
		nil,
		"",
	)
	f(
		`{{ trim . }}`,
		"\tkratos\n ",
		"kratos",
		nil,
		"",
	)
	f(
		`{{ uppercase . }}`,
		"  ",
		"  ",
		nil,
		"",
	)
	f(
		`{{ uppercase . }}`,
		"kratos",
		"KRATOS",
		nil,
		"",
	)
	f(
		`{{ uppercase . }}`,
		"KraTos",
		"KRATOS",
		nil,
		"",
	)
	f(
		`{{ lowercase . }}`,
		"  ",
		"  ",
		nil,
		"",
	)
	f(
		`{{ lowercase . }}`,
		"kratos",
		"kratos",
		nil,
		"",
	)
	f(
		`{{ lowercase . }}`,
		"KraTos",
		"kratos",
		nil,
		"",
	)
	f(
		`{{ lowercase . }}`,
		"KRATOS",
		"kratos",
		nil,
		"",
	)
	f(
		`{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" "Australia/Queensland" }}`,
		map[string]any{
			"TestDate": datetime.CpsTime{
				Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
			},
		},
		"Thu, 28 Oct 2021 17:05:00 AEST",
		nil,
		"",
	)
	f(
		`{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" }}`,
		map[string]any{
			"TestDate": datetime.CpsTime{
				Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
			},
		},
		"Thu, 28 Oct 2021 17:05:00 AEST",
		nil,
		"Australia/Queensland",
	)
	f(
		`{{tag_has_key .Alarm.Tags "Tag1" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"true",
		nil,
		"",
	)
	f(
		`{{tag_has_key .Alarm.Tags "Tag2" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"true",
		nil,
		"",
	)
	f(
		`{{tag_has_key .Alarm.Tags "Tag3" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"false",
		nil,
		"",
	)
	f(
		`{{get_tag .Alarm.Tags "Tag1" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"Value1",
		nil,
		"",
	)
	f(
		`{{get_tag .Alarm.Tags "Tag2" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"",
		nil,
		"",
	)
	f(
		`{{get_tag .Alarm.Tags "Tag3" }}`,
		map[string]any{
			"Alarm": types.Alarm{
				Tags: []string{"Tag1: Value1", "Tag2"},
			},
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrLeft .String 5 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello",
		nil,
		"",
	)
	f(
		`{{substrLeft .String 0 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrLeft .String 13 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello, world!",
		nil,
		"",
	)
	f(
		`{{substrLeft .String 100000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello, world!",
		nil,
		"",
	)
	f(
		`{{substrLeft .String -100000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrLeft .String 8 }}`,
		map[string]any{
			"String": "élémentaire",
		},
		"élémenta",
		nil,
		"",
	)
	f(
		`{{substrRight .String 6 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"world!",
		nil,
		"",
	)
	f(
		`{{substrRight .String 0 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrRight .String 13 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello, world!",
		nil,
		"",
	)
	f(
		`{{substrRight .String 100000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello, world!",
		nil,
		"",
	)
	f(
		`{{substrRight .String -100000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrRight .String 6 }}`,
		map[string]any{
			"String": "élémenta",
		},
		"émenta",
		nil,
		"",
	)
	f(
		`{{substr .String 0 5 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello",
		nil,
		"",
	)
	f(
		`{{substr .String 7 5 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"world",
		nil,
		"",
	)
	f(
		`{{substr .String 7 1000000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"world!",
		nil,
		"",
	)
	f(
		`{{substr .String 0 1000000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello, world!",
		nil,
		"",
	)
	f(
		`{{substr .String -1 1000000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substr .String 5 0 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substr .String 0 -1000000 }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substr .String 1 2 }}`,
		map[string]any{
			"String": "élémenta",
		},
		"lé",
		nil,
		"",
	)
	f(
		`{{strlen .String }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"13",
		nil,
		"",
	)
	f(
		`{{strlen .String }}`,
		map[string]any{
			"String": "élémenta",
		},
		"8",
		nil,
		"",
	)
	f(
		`{{strpos .String "w" }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"7",
		nil,
		"",
	)
	f(
		`{{strpos .String "ém" }}`,
		map[string]any{
			"String": "élémenta",
		},
		"2",
		nil,
		"",
	)
	f(
		`{{strpos .String "invalid" }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"-1",
		nil,
		"",
	)
	f(
		`{{ $var := add .A .B }}{{ substrLeft .String $var }}`,
		map[string]any{
			"A":      1,
			"B":      2,
			"String": "qwerty",
		},
		"qwe",
		nil,
		"",
	)
	f(
		`{{substrLeft .String (sub (strlen .String) 8) }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello",
		nil,
		"",
	)
	f(
		`{{substrLeft .String (strpos .String ",") }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"Hello",
		nil,
		"",
	)
	f(
		`{{substrLeft .String (sub (strlen .String) 6) }}`,
		map[string]any{
			"String": "élémenta",
		},
		"él",
		nil,
		"",
	)
	f(
		`{{substrLeft .String (strpos .String "ém") }}`,
		map[string]any{
			"String": "élémenta",
		},
		"él",
		nil,
		"",
	)
	f(
		`{{substrLeft .String (strpos .String "invalid") }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{substrRight .String (sub (strlen .String) 7) }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"world!",
		nil,
		"",
	)
	f(
		`{{substrRight .String (sub (strlen .String) (strpos .String "w")) }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"world!",
		nil,
		"",
	)
	f(
		`{{substrRight .String (sub (strlen .String) 2) }}`,
		map[string]any{
			"String": "élémenta",
		},
		"émenta",
		nil,
		"",
	)
	f(
		`{{substrRight .String (sub (strlen .String) (strpos .String "ém")) }}`,
		map[string]any{
			"String": "élémenta",
		},
		"émenta",
		nil,
		"",
	)
	f(
		`{{substrRight .String (strpos .String "invalid") }}`,
		map[string]any{
			"String": "Hello, world!",
		},
		"",
		nil,
		"",
	)
	f(
		`{{ add .A .B }}`,
		map[string]any{
			"A": 2,
			"B": 2,
		},
		"4",
		nil,
		"",
	)
	f(
		`{{ add .A .B }}`,
		map[string]any{
			"A": "2.5",
			"B": 2.5,
		},
		"",
		ErrFailedConvertToInt64,
		"",
	)
	f(
		`{{ sub .A .B }}`,
		map[string]any{
			"A": 4,
			"B": 2,
		},
		"2",
		nil,
		"",
	)
	f(
		`{{ sub .A .B }}`,
		map[string]any{
			"A": "5",
			"B": 2.5,
		},
		"",
		ErrFailedConvertToInt64,
		"",
	)
	f(
		`{{ mult .A .B }}`,
		map[string]any{
			"A": 2,
			"B": 2,
		},
		"4",
		nil,
		"",
	)
	f(
		`{{ mult .A .B }}`,
		map[string]any{
			"A": "2.5",
			"B": 2.5,
		},
		"",
		ErrFailedConvertToInt64,
		"",
	)
	f(
		`{{ div .A .B }}`,
		map[string]any{
			"A": 4,
			"B": 2,
		},
		"2",
		nil,
		"",
	)
	f(
		`{{ div .A .B }}`,
		map[string]any{
			"A": 13,
			"B": 2,
		},
		"6",
		nil,
		"",
	)
	f(
		`{{ div .A .B }}`,
		map[string]any{
			"A": 6,
			"B": 0,
		},
		"",
		ErrDivisionByZero,
		"",
	)
	f(
		`{{ div .A .B }}`,
		map[string]any{
			"A": "6.25",
			"B": 2.5,
		},
		"",
		ErrFailedConvertToInt64,
		"",
	)
	f(
		`{{ strjoin .A "," }}`,
		map[string]any{
			"A": []string{"a", "b", "c"},
		},
		"a,b,c",
		nil,
		"",
	)
	f(
		`{{ strjoin .A "," }}`,
		map[string]any{
			"A": []any{"a", "b", "c"},
		},
		"a,b,c",
		nil,
		"",
	)
	f(
		`{{ strjoin .A "," }}`,
		map[string]any{
			"A": []any{"a"},
		},
		"a",
		nil,
		"",
	)
	f(
		`{{ strjoin .A "," }}`,
		map[string]any{
			"A": []any{},
		},
		"",
		nil,
		"",
	)
	f(
		`{{ strjoin .A "," }}`,
		map[string]any{
			"A": []any{1, 2, 3},
		},
		"",
		ErrFailedConvertToStringSlice,
		"",
	)

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

	f := func(data, expected any) {
		t.Helper()
		res := addDefaultTplVarsToData(data, map[string]any{EnvVar: envVars})
		if diff := pretty.Compare(res, expected); diff != "" {
			t.Fatalf("unexpected result %s", diff)
		}
	}

	f(
		map[string]types.Alarm{
			"Alarm": alarm,
		},
		map[string]any{
			"Alarm": alarm,
			"Env":   envVars,
		},
	)
	f(
		alarm,
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		&alarm,
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		Activatable(&alarm),
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		alarmWithEnabled{
			Alarm:   alarm,
			Enabled: true,
		},
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		alarmPtrWithEnabled{
			Alarm:   &alarm,
			Enabled: true,
		},
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		activatableWithEnabled{
			Activatable: Activatable(&alarm),
			Enabled:     true,
		},
		map[string]any{
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
			"CopyTagsFromChildren":              alarm.CopyTagsFromChildren,
			"FilterChildrenTagsByLabel":         alarm.FilterChildrenTagsByLabel,
			"EntityInfosFromChildren":           alarm.EntityInfosFromChildren,
			"Time":                              alarm.Time,
			"Value":                             alarm.Value,
			"NotAckedMetricSendTime":            alarm.NotAckedMetricSendTime,
			"NotAckedMetricType":                alarm.NotAckedMetricType,
			"NotAckedSince":                     alarm.NotAckedSince,
			"InactiveAutoInstructionInProgress": alarm.InactiveAutoInstructionInProgress,
			"InactiveDelayMetaAlarmInProgress":  alarm.InactiveDelayMetaAlarmInProgress,
			"MetaAlarmInactiveDelay":            alarm.MetaAlarmInactiveDelay,
			"Healthcheck":                       alarm.Healthcheck,
			"Env":                               envVars,
		},
	)
	f(
		map[int]types.Alarm{
			1: alarm,
		},
		map[int]types.Alarm{
			1: alarm,
		},
	)
	f([]types.Alarm{alarm}, []types.Alarm{alarm})
	f(1, 1)
}

func executeTemplate(tmpl *template.Template, payload interface{},
) (string, error) {
	var b bytes.Buffer
	err := tmpl.Execute(io.Writer(&b), payload)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
