package validator_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

func TestValidator_Validate(t *testing.T) {
	v := getValidator()
	f := func(str string, data any, expectedIsValid bool, expectedErrReport *validator.ErrReport) {
		t.Helper()

		isValid, errReport, err := v.Validate(str, data)
		if err != nil {
			t.Errorf("error is not expected, but got %s", err.Error())
		}

		if isValid != expectedIsValid {
			t.Errorf("expected valid = %t, got %t", expectedIsValid, isValid)
		}

		if !isValid {
			if errReport == nil {
				t.Error("report shouldn't be nil")
			} else if expectedErrReport == nil {
				t.Error("expected report shouldn't be nil")
			} else if *errReport != *expectedErrReport {
				t.Errorf("expected error report = %+v, got %+v", expectedErrReport, *errReport)
			}
		}
	}

	event := types.Event{
		Connector: "test-connector",
		ExtraInfos: map[string]any{
			"info1": "val1",
		},
	}
	data := struct {
		Event types.Event
	}{
		Event: event,
	}

	f("test {{ range .Alarms }} {{ .Value.Re3source } test {{ end }}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Unexpected \"}\"",
	})
	f("test {{ range .Alarms }}\n{{ .Value.Re3source } test {{ end }}", nil, false, &validator.ErrReport{
		Line:    2,
		Message: "Unexpected \"}\"",
	})
	f("test {{ range .Alarms }} {{ .Value.Resource }} test {{ end }", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Unexpected \"}\"",
	})
	f("test {{ rangee .Alarms }} {{ .Value.Resource }} test {{ end }}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Invalid function \"rangee\"",
	})
	f("test {{ range .Alarms }} {{ .Value.Resource | some }} test {{ end }}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Invalid function \"some\"",
	})
	f("test {{ range .Alarms }} {{ .Value.Resource }} test", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Parsing error: invalid template",
	})
	f("test {{ if .Alarms }} {{ .Value.Resource }} test", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Parsing error: invalid template",
	})
	f("test {{end}}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Function or block is missing",
	})
	f("test {{ break }} test", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "{{break}} outside {{range}}",
	})
	f("{ index .Response \"test\" }}", nil, true, nil)
	f("test {{ .Notexist }}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"Notexist\"",
	})
	f("test {{ .Notexist }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"Notexist\"",
	})
	f("test {{ .Event.Notexist }}", nil, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"Event\"",
	})
	f("test {{ .Event.Notexist }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"Notexist\"",
	})
	f("test {{ index .Event.ExtraInfos \"info1\" }}", data, true, nil)
	f("test {{ .Event.ExtraInfos.info1 }}", data, true, nil)
	f("test {{ index .Event.ExtraInfos \"notexist\" }}", data, true, nil)
	f("test {{ .Event.ExtraInfos.notexist }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"notexist\"",
	})
	f("test {{ index .Event.Entity.Infos \"info1\" }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Undefined key or field \".Event.Entity.Infos\"",
	})
	f("test {{ .Event.Entity.Infos.info1 }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Undefined key or field \".Event.Entity.Infos.info1\"",
	})
	f("test {{ .Event.Entity.Name }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Undefined key or field \".Event.Entity.Name\"",
	})
	f("test {{ .Event.Entity.Notexist }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Unknown key or field \"Notexist\"",
	})
	f("test {{ .Event.AlarmChange.Type }}", data, false, &validator.ErrReport{
		Line:    1,
		Message: "Undefined key or field \".Event.AlarmChange.Type\"",
	})
}

func getValidator() validator.Validator {
	cfg := config.CanopsisConf{}
	cfg.Timezone.Timezone = "Europe/Paris"
	templateConfigProvider := config.NewTemplateConfigProvider(cfg, zerolog.Nop())
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, zerolog.Nop())

	return validator.NewValidator(template.NewExecutor(templateConfigProvider, timezoneConfigProvider))
}
