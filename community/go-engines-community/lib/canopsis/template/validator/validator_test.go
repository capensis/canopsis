package validator_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	"github.com/golang/mock/gomock"
)

type dataSet struct {
	testName           string
	testString         string
	expectedErrReport  validator.ErrReport
	expectedWrnReports []validator.WrnReport
	expectedValid      bool
}

func TestValidator_ValidateDeclareTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataSets := []dataSet{
		{
			testName:          "validate with range Alarms",
			testString:        "test {{ range .Alarms }} test {{ .Value.Output }} {{ .Entity.Name }} {{ end }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Response",
			testString:        "test {{ index .Response \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with ResponseMap",
			testString:        "test {{ index .ResponseMap \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Header",
			testString:        "test {{ index .Header \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:   "validate unexpected symbol in operand",
			testString: "test {{ range .Alarms }} {{ .Value.Re3source } test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol in operand in the second line",
			testString: "test {{ range .Alarms }}\n{{ .Value.Re3source } test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    2,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol at the end",
			testString: "test {{ range .Alarms }} {{ .Value.Resource }} test {{ end }",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected function 1",
			testString: "test {{ rangee .Alarms }} {{ .Value.Resource }} test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedFunction,
				Message: "Invalid function \"rangee\"",
			},
		},
		{
			testName:   "validate unexpected function 2",
			testString: "test {{ range .Alarms }} {{ .Value.Resource | some }} test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedFunction,
				Message: "Invalid function \"some\"",
			},
		},
		{
			testName:   "validate unexpected EOF 1",
			testString: "test {{ range .Alarms }} {{ .Value.Resource }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected EOF 2",
			testString: "test {{ if .Alarms }} {{ .Value.Resource }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected block",
			testString: "test {{end}}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedBlock,
				Message: "Function or block is missing",
			},
		},
		{
			testName:   "validate other errors",
			testString: "test {{ break }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUndefined,
				Message: "{{break}} outside {{range}}",
			},
		},
		{
			testName:   "validate with unfinished variable block",
			testString: "{ index .Response \"test\" }}",
			expectedWrnReports: []validator.WrnReport{
				{
					Type:    validator.WrnTypeOutsideBlockVar,
					Message: "Variable is out of a template block",
					Var:     ".Response",
				},
			},
			expectedValid: true,
		},
	}

	v := getValidator(ctrl)

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			isValid, errReport, wrnReports, err := v.ValidateDeclareTicketRuleTemplate(dataset.testString)
			if err != nil {
				t.Errorf("error is not expected, but got %s", err.Error())
			}

			if isValid != dataset.expectedValid {
				t.Errorf("expected valid = %t, got %t", dataset.expectedValid, isValid)
			}

			if !isValid {
				if errReport == nil {
					t.Error("report shouldn't be nil")
				} else if *errReport != dataset.expectedErrReport {
					t.Errorf("expected error report = %v, got %v", dataset.expectedErrReport, *errReport)
				}
			} else if !reflect.DeepEqual(dataset.expectedWrnReports, wrnReports) {
				t.Errorf("expected warning reports = %v, got %v", dataset.expectedWrnReports, wrnReports)
			}
		})
	}
}

func TestValidator_ValidateScenarioWebhookTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataSets := []dataSet{
		{
			testName:          "validate with Alarm",
			testString:        "test {{ .Alarm.Value.Output }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Entity",
			testString:        "test {{ .Entity.Name }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Children",
			testString:        "test {{ range .Children }} {{ .Value.Output }} {{ end }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Response",
			testString:        "test {{ index .Response \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with ResponseMap",
			testString:        "test {{ index .ResponseMap \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with Header",
			testString:        "test {{ index .Header \"test\" }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:          "validate with AdditionalData",
			testString:        "test {{ .AdditionalData.Author }}",
			expectedErrReport: validator.ErrReport{},
			expectedValid:     true,
		},
		{
			testName:   "validate unexpected symbol in operand",
			testString: "test {{ range .Children }} {{ .Value.Re3source } test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol in operand in the second line",
			testString: "test {{ range .Children }}\n{{ .Value.Re3source } test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    2,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected symbol at the end",
			testString: "test {{ range .Children }} {{ .Value.Resource }} test {{ end }",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedSymbol,
				Message: "Unexpected \"}\"",
			},
		},
		{
			testName:   "validate unexpected function 1",
			testString: "test {{ rangee .Children }} {{ .Value.Resource }} test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedFunction,
				Message: "Invalid function \"rangee\"",
			},
		},
		{
			testName:   "validate unexpected function 2",
			testString: "test {{ range .Children }} {{ .Value.Resource | some }} test {{ end }}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedFunction,
				Message: "Invalid function \"some\"",
			},
		},
		{
			testName:   "validate unexpected EOF 1",
			testString: "test {{ range .Children }} {{ .Value.Resource }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected EOF 2",
			testString: "test {{ if .Children }} {{ .Value.Resource }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedEOF,
				Message: "Parsing error: invalid template",
			},
		},
		{
			testName:   "validate unexpected block",
			testString: "test {{end}}",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUnexpectedBlock,
				Message: "Function or block is missing",
			},
		},
		{
			testName:   "validate other errors",
			testString: "test {{ break }} test",
			expectedErrReport: validator.ErrReport{
				Line:    1,
				Type:    validator.ErrTypeUndefined,
				Message: "{{break}} outside {{range}}",
			},
		},
		{
			testName:   "validate with unfinished variable block",
			testString: "test {{ index .ResponseMap \"key1\"}} { index .Response \"test\" }} test {{ range .Alarms }} {{ .Value.Resource }} test {{ end }}",
			expectedWrnReports: []validator.WrnReport{
				{
					Type:    validator.WrnTypeOutsideBlockVar,
					Message: "Variable is out of a template block",
					Var:     ".Response",
				},
			},
			expectedValid: true,
		},
	}

	v := getValidator(ctrl)

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			isValid, errReport, wrnReports, err := v.ValidateScenarioTemplate(dataset.testString)
			if err != nil {
				t.Errorf("error is not expected, but got %s", err.Error())
			}

			if isValid != dataset.expectedValid {
				t.Errorf("expected valid = %t, got %t", dataset.expectedValid, isValid)
			}

			if !isValid {
				if errReport == nil {
					t.Error("report shouldn't be nil")
				} else if *errReport != dataset.expectedErrReport {
					t.Errorf("expected error report = %v, got %v", dataset.expectedErrReport, *errReport)
				}
			} else if !reflect.DeepEqual(dataset.expectedWrnReports, wrnReports) {
				t.Errorf("expected warning reports = %v, got %v", dataset.expectedWrnReports, wrnReports)
			}
		})
	}
}

func getValidator(ctrl *gomock.Controller) validator.Validator {
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().Return(config.TimezoneConfig{}).AnyTimes()

	return validator.NewValidator(mockTimezoneConfigProvider)
}
