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
	ParseErrorMatches = 3
	ExecErrorMatches  = 5
)

const (
	ErrorTypeUndefined = iota
	ErrorTypeNoSuchMainVariable
	ErrorTypeNoSuchSecondaryVariable
	ErrorTypeUnexpectedBlock
	ErrorTypeUnexpectedSymbol
	ErrorTypeUnexpectedFunction
	ErrorTypeUnexpectedEOF
)

type RegexpInfo struct {
	errRegexp     *regexp.Regexp
	errType       int
	matchesNumber int
	getErrMessage func([]string) string
}

type Validator interface {
	ValidateDeclareTicketTemplate(s string) (bool, *Report)
}

type validator struct {
	executor *template.Executor

	execErrorRegex       *regexp.Regexp
	secondaryVarErrRegex *regexp.Regexp

	parseErrorRegex         *regexp.Regexp
	parseErrorsMsgRegexInfo []RegexpInfo
}

func NewValidator(executor *template.Executor) Validator {
	return &validator{
		executor: executor,

		execErrorRegex:       regexp.MustCompile("^template: (.+): executing (.+) at <(.+)>: (.+)$"),
		secondaryVarErrRegex: regexp.MustCompile("^can't evaluate field (.+) in type (.+)$"),

		parseErrorRegex: regexp.MustCompile("^template: (.+): (.+)$"),
		parseErrorsMsgRegexInfo: []RegexpInfo{
			{
				errType:       ErrorTypeUnexpectedFunction,
				errRegexp:     regexp.MustCompile("^function \"(.+)\" not defined$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Invalid function \"%s\"", matches[1])
				},
			},
			{
				errType:       ErrorTypeUnexpectedSymbol,
				errRegexp:     regexp.MustCompile("^unexpected \"(.+)\" in (.+)$"),
				matchesNumber: 3,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Unexpected \"%s\"", matches[1])
				},
			},
			{
				errType:       ErrorTypeUnexpectedEOF,
				errRegexp:     regexp.MustCompile("^unexpected EOF$"),
				matchesNumber: 1,
				getErrMessage: func(_ []string) string {
					return "Parsing error: invalid template"
				},
			},
			{
				errType:       ErrorTypeUnexpectedBlock,
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
	}
}

type Report struct {
	Line     int `json:"line,omitempty"`
	Position int `json:"position,omitempty"`

	// Possible trigger values.
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

func (v *validator) ValidateDeclareTicketTemplate(s string) (bool, *Report) {
	tplData := map[string]interface{}{
		"Alarm":    types.Alarm{},
		"Entity":   types.Entity{},
		"Children": make([]types.Alarm, 0),
		"Alarms": []struct {
			types.Alarm
			Entity types.Entity
		}{{}},
		"Response":       map[string]any{},
		"ResponseMap":    map[string]any{},
		"Header":         map[string]string{},
		"AdditionalData": action.AdditionalData{},
	}

	_, err := v.executor.Execute(s, tplData)
	if err != nil {
		fullErrString := err.Error()
		report := &Report{
			Type:    ErrorTypeUndefined,
			Message: fullErrString,
		}

		// parse template exec error
		tplErrorMatches := v.execErrorRegex.FindStringSubmatch(fullErrString)
		if len(tplErrorMatches) == ExecErrorMatches {
			report.Line, report.Position = getLocation(tplErrorMatches[1])

			report.Var = tplErrorMatches[3]
			report.Message = tplErrorMatches[4]

			if strings.HasPrefix(report.Message, "map has no entry for key") {
				report.Type = ErrorTypeNoSuchMainVariable
				report.Message = fmt.Sprintf("No such variable \"%s\"", report.Var)
			} else if strings.HasPrefix(report.Message, "can't evaluate field") {
				report.Type = ErrorTypeNoSuchSecondaryVariable

				errMsgMatches := v.secondaryVarErrRegex.FindStringSubmatch(report.Message)
				if len(errMsgMatches) == 3 {
					report.Message = fmt.Sprintf("Invalid variable \"%s\"", errMsgMatches[1])
					report.Var = errMsgMatches[1]
				}
			}

			return false, report
		}

		// parse template parse error
		tplErrorMatches = v.parseErrorRegex.FindStringSubmatch(fullErrString)
		if len(tplErrorMatches) == ParseErrorMatches {
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

		return false, report
	}

	return true, nil
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
