package validator_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	"github.com/golang/mock/gomock"
)

func TestValidator_ValidateDeclareTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataSets := []struct {
		testName       string
		testString     string
		expectedReport validator.Report
		expectedValid  bool
	}{
		{
			testName:       "validate with success",
			testString:     "test {{ range .Alarms }} test {{ end }}",
			expectedReport: validator.Report{},
			expectedValid:  true,
		},
		{
			testName:   "validate main variable",
			testString: "test {{ .Alarmmm }} test",
			expectedReport: validator.Report{
				Line:     1,
				Position: 8,
				Type:     validator.ErrorTypeNoSuchMainVariable,
				Message:  "No such variable \".Alarmmm\"",
				Var:      ".Alarmmm",
			},
		},
		{
			testName:   "validate main variable on the second line",
			testString: "test\n{{ .Alarmmm }} test",
			expectedReport: validator.Report{
				Line:     2,
				Position: 3,
				Type:     validator.ErrorTypeNoSuchMainVariable,
				Message:  "No such variable \".Alarmmm\"",
				Var:      ".Alarmmm",
			},
		},
		{
			testName:   "validate secondary variable",
			testString: "test {{ .AdditionalData.Some }}",
			expectedReport: validator.Report{
				Line:     1,
				Position: 23,
				Type:     validator.ErrorTypeNoSuchSecondaryVariable,
				Message:  "Invalid variable \"Some\"",
				Var:      "Some",
			},
		},
		{
			testName:   "validate unexpected symbol in operand",
			testString: "test {{ range .Alarms }} {{ .Value.Re3source } test {{ end }}",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol in operand in the second line",
			testString: "test {{ range .Alarms }}\n{{ .Value.Re3source } test {{ end }}",
			expectedReport: validator.Report{
				Line:    2,
				Type:    validator.ErrorTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol at the end",
			testString: "test {{ range .Alarms }} {{ .Value.Resource }} test {{ end }",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected function 1",
			testString: "test {{ rangee .Alarms }} {{ .Value.Resource }} test {{ end }}",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedFunction,
				Message: "Invalid function \"rangee\"",
			},
		},
		{
			testName:   "validate unexpected function 2",
			testString: "test {{ range .Alarms }} {{ .Value.Resource | some }} test {{ end }}",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedFunction,
				Message: "Invalid function \"some\"",
			},
		},
		{
			testName:   "validate unexpected EOF 1",
			testString: "test {{ range .Alarms }} {{ .Value.Resource }} test",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected EOF 2",
			testString: "test {{ if .Alarms }} {{ .Value.Resource }} test",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected block",
			testString: "test {{end}}",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUnexpectedBlock,
				Message: "Function or block is missing",
			},
		},
		{
			testName:   "validate other errors",
			testString: "test {{ break }} test",
			expectedReport: validator.Report{
				Line:    1,
				Type:    validator.ErrorTypeUndefined,
				Message: "{{break}} outside {{range}}",
			},
		},
	}

	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().Return(config.TimezoneConfig{}).AnyTimes()

	v := validator.NewValidator(template.NewExecutor(mockTimezoneConfigProvider))

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			isValid, report := v.ValidateDeclareTicketTemplate(dataset.testString)
			if isValid != dataset.expectedValid {
				t.Errorf("expected valid = %t, got %t", dataset.expectedValid, isValid)
			}

			if !isValid {
				if report == nil {
					t.Error("report shouldn't be nil")
				} else if *report != dataset.expectedReport {
					t.Errorf("expected report = %v, got %v", dataset.expectedReport, report)
				}
			}
		})
	}
}
