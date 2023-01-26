package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	parseErrorMatches = 3
	execErrorMatches  = 5
)

const (
	ErrTypeUndefined = iota
	ErrTypeNoSuchMainVariable
	ErrTypeNoSuchSecondaryVariable
	ErrTypeUnexpectedBlock
	ErrTypeUnexpectedSymbol
	ErrTypeUnexpectedFunction
	ErrTypeUnexpectedEOF
)

const (
	WrnTypeOutsideBlockVar = iota
)

type RegexpInfo struct {
	errRegexp     *regexp.Regexp
	errType       int
	matchesNumber int
	getErrMessage func([]string) string
}

type Validator interface {
	ValidateDeclareTicketRuleTemplate(s string) (bool, *ErrReport, []WrnReport)
	ValidateScenarioTemplate(s string) (bool, *ErrReport, []WrnReport)
}

type validator struct {
	executor *template.Executor

	execErrorRegex       *regexp.Regexp
	secondaryVarErrRegex *regexp.Regexp

	parseErrorRegex         *regexp.Regexp
	parseErrorsMsgRegexInfo []RegexpInfo

	declareTicketTplData map[string]any
	scenarioTplData      map[string]any
}

func NewValidator(executor *template.Executor) Validator {
	alarm := types.Alarm{
		Value: types.AlarmValue{
			ACK:         &types.AlarmStep{},
			Canceled:    &types.AlarmStep{},
			Done:        &types.AlarmStep{},
			Snooze:      &types.AlarmStep{},
			State:       &types.AlarmStep{},
			Status:      &types.AlarmStep{},
			LastComment: &types.AlarmStep{},
			Ticket:      &types.AlarmStep{},
			Tickets:     []types.AlarmStep{{}},
			Steps:       []types.AlarmStep{{}},
		},
	}

	return &validator{
		executor: executor,

		execErrorRegex:       regexp.MustCompile("^template: (.+): executing (.+) at <(.+)>: (.+)$"),
		secondaryVarErrRegex: regexp.MustCompile("^can't evaluate field (.+) in type (.+)$"),

		parseErrorRegex: regexp.MustCompile("^template: (.+): (.+)$"),
		parseErrorsMsgRegexInfo: []RegexpInfo{
			{
				errType:       ErrTypeUnexpectedFunction,
				errRegexp:     regexp.MustCompile("^function \"(.+)\" not defined$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Invalid function \"%s\"", matches[1])
				},
			},
			{
				errType:       ErrTypeUnexpectedSymbol,
				errRegexp:     regexp.MustCompile("^unexpected \"(.+)\" in (.+)$"),
				matchesNumber: 3,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Unexpected \"%s\"", matches[1])
				},
			},
			{
				errType:       ErrTypeUnexpectedEOF,
				errRegexp:     regexp.MustCompile("^unexpected EOF$"),
				matchesNumber: 1,
				getErrMessage: func(_ []string) string {
					return "Parsing error: invalid template"
				},
			},
			{
				errType:       ErrTypeUnexpectedBlock,
				errRegexp:     regexp.MustCompile("^unexpected (.+)$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					if matches[1] == "{{end}}" {
						return "Function or block is missing"
					}

					return fmt.Sprintf("Unexpected \"%s\"", matches[1])
				},
			},
		},

		declareTicketTplData: map[string]any{
			"Alarms": []struct {
				types.Alarm
				Entity types.Entity
			}{{
				Alarm:  alarm,
				Entity: types.Entity{},
			}},
			"Response":    map[string]any{},
			"ResponseMap": map[string]any{},
			"Header":      map[string]string{},
		},
		scenarioTplData: map[string]any{
			"Alarm":          alarm,
			"Entity":         types.Entity{},
			"Children":       []types.Alarm{alarm},
			"Response":       map[string]any{},
			"ResponseMap":    map[string]any{},
			"Header":         map[string]string{},
			"AdditionalData": action.AdditionalData{},
		},
	}
}

type ErrReport struct {
	Line     int `json:"line,omitempty"`
	Position int `json:"position,omitempty"`

	// Possible error type values.
	//   * `0` - Undefined error
	//   * `1` - No main variable
	//   * `2` - No secondary variable
	//   * `3` - Unexpected block
	//   * `4` - Unexpected symbol
	//   * `5` - Unexpected function
	//   * `6` - Unexpected EOF
	Type    int    `json:"type"`
	Message string `json:"message"`

	// Var is defined only for template exec errors
	Var string `json:"var,omitempty"`
}

type WrnReport struct {
	// Possible warning type variables
	//   * `0` - Might be unfinished variable block
	Type int `json:"type"`

	Message string `json:"message"`
	Var     string `json:"var,omitempty"`
}

func (v *validator) ValidateDeclareTicketRuleTemplate(s string) (bool, *ErrReport, []WrnReport) {
	return v.validate(s, v.declareTicketTplData)
}

func (v *validator) ValidateScenarioTemplate(s string) (bool, *ErrReport, []WrnReport) {
	return v.validate(s, v.scenarioTplData)
}

func (v *validator) validate(s string, tplData map[string]any) (bool, *ErrReport, []WrnReport) {
	res, err := v.executor.Execute(s, tplData)
	if err != nil {
		fullErrString := err.Error()
		report := &ErrReport{
			Type:    ErrTypeUndefined,
			Message: fullErrString,
		}

		// parse template exec error
		tplErrorMatches := v.execErrorRegex.FindStringSubmatch(fullErrString)
		if len(tplErrorMatches) == execErrorMatches {
			report.Line, report.Position = getLocation(tplErrorMatches[1])

			report.Var = tplErrorMatches[3]
			report.Message = tplErrorMatches[4]

			if strings.HasPrefix(report.Message, "map has no entry for key") {
				report.Type = ErrTypeNoSuchMainVariable
				report.Message = fmt.Sprintf("No such variable \"%s\"", report.Var)
			} else if strings.HasPrefix(report.Message, "can't evaluate field") {
				report.Type = ErrTypeNoSuchSecondaryVariable

				errMsgMatches := v.secondaryVarErrRegex.FindStringSubmatch(report.Message)
				if len(errMsgMatches) == 3 {
					report.Message = fmt.Sprintf("Invalid variable \"%s\"", errMsgMatches[1])
					report.Var = errMsgMatches[1]
				}
			}

			return false, report, nil
		}

		// parse template parse error
		tplErrorMatches = v.parseErrorRegex.FindStringSubmatch(fullErrString)
		if len(tplErrorMatches) == parseErrorMatches {
			report.Line, report.Position = getLocation(tplErrorMatches[1])
			report.Message = tplErrorMatches[2]

			for _, regexInfo := range v.parseErrorsMsgRegexInfo {
				errMsgMatches := regexInfo.errRegexp.FindStringSubmatch(report.Message)
				if len(errMsgMatches) == regexInfo.matchesNumber {
					report.Type = regexInfo.errType
					report.Message = regexInfo.getErrMessage(errMsgMatches)

					break
				}
			}
		}

		return false, report, nil
	}

	var warnings []WrnReport

	// try to find template variables in the result text
	for k := range tplData {
		val := "." + k
		if strings.Contains(res, val) {
			warnings = append(warnings, WrnReport{
				Type:    WrnTypeOutsideBlockVar,
				Message: "Variable are out of a template block",
				Var:     val,
			})
		}
	}

	return true, nil, warnings
}

func getLocation(s string) (int, int) {
	var err error
	line, position := 0, 0

	locationSplit := strings.Split(s, ":")
	if len(locationSplit) < 2 {
		panic(fmt.Errorf("template exec error contains invalid location value = %s", s))
	}

	line, err = strconv.Atoi(locationSplit[1])
	if err != nil {
		panic(fmt.Errorf("convert line variable to int error = %w", err))
	}

	if len(locationSplit) > 2 {
		position, err = strconv.Atoi(locationSplit[2])
		if err != nil {
			panic(fmt.Errorf("convert position variable to int error = %w", err))
		}
	}

	return line, position
}
