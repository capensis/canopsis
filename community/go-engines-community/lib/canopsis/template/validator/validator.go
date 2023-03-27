package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

const (
	locationStringMatch = 1
	messageStringMatch  = 2
	parseErrorMatches   = 3
)

const (
	ErrTypeUndefined = iota
	ErrTypeUnexpectedBlock
	ErrTypeUnexpectedSymbol
	ErrTypeUnexpectedFunction
	ErrTypeUnexpectedEOF
)

const (
	WrnTypeOutsideBlockVar = iota
)

type RegexpInfo struct {
	errRegexp *regexp.Regexp
	errType   int
	// matchesNumber equals number of matched groups in errRegexp + 1
	matchesNumber int
	getErrMessage func([]string) string
}

type Validator interface {
	ValidateDeclareTicketRuleTemplate(s string) (bool, *ErrReport, []WrnReport, error)
	ValidateScenarioTemplate(s string) (bool, *ErrReport, []WrnReport, error)
}

type validator struct {
	timezoneConfigProvider config.TimezoneConfigProvider

	parseErrorRegex         *regexp.Regexp
	parseErrorsMsgRegexInfo []RegexpInfo

	declareTicketTplDataKeys []string
	scenarioTplDataKeys      []string

	warningOutOfBlockRegex *regexp.Regexp
}

func NewValidator(timezoneConfigProvider config.TimezoneConfigProvider) Validator {
	return &validator{
		timezoneConfigProvider: timezoneConfigProvider,

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
		declareTicketTplDataKeys: []string{
			"." + libtemplate.EnvVar,
			".Alarms",
			".Response",
			".ResponseMap",
			".Header",
			// for range alarms case
			".Value",
			".Entity",
		},
		scenarioTplDataKeys: []string{
			"." + libtemplate.EnvVar,
			".Alarm",
			".Entity",
			".Children",
			".Response",
			".ResponseMap",
			".Header",
			".AdditionalData",
			// for range children case
			".Value",
		},
		warningOutOfBlockRegex: regexp.MustCompile(`\{{2}[^\}]*\}{2}`),
	}
}

type ErrReport struct {
	Line int `json:"line"`

	// Possible error type values.
	//   * `0` - Undefined error
	//   * `1` - Unexpected block
	//   * `2` - Unexpected symbol
	//   * `3` - Unexpected function
	//   * `4` - Unexpected EOF
	Type    int    `json:"type"`
	Message string `json:"message"`

	// Var is defined only for template exec errors
	Var string `json:"var,omitempty"`
}

type WrnReport struct {
	// Possible warning type values.
	//   * `0` - Might be unfinished variable block
	Type int `json:"type"`

	Message string `json:"message"`
	Var     string `json:"var,omitempty"`
}

func (v *validator) ValidateDeclareTicketRuleTemplate(s string) (bool, *ErrReport, []WrnReport, error) {
	return v.validate(s, v.declareTicketTplDataKeys)
}

func (v *validator) ValidateScenarioTemplate(s string) (bool, *ErrReport, []WrnReport, error) {
	return v.validate(s, v.scenarioTplDataKeys)
}

func (v *validator) validate(s string, tplKeys []string) (bool, *ErrReport, []WrnReport, error) {
	location := v.timezoneConfigProvider.Get().Location

	_, err := template.New("tpl").Funcs(libtemplate.GetFunctions(location)).Parse(s)
	if err != nil {
		fullErrString := err.Error()
		report := &ErrReport{
			Type:    ErrTypeUndefined,
			Message: fullErrString,
		}

		// parse template parse error
		tplErrorMatches := v.parseErrorRegex.FindStringSubmatch(fullErrString)
		if len(tplErrorMatches) == parseErrorMatches {
			report.Line, err = getLine(tplErrorMatches[locationStringMatch])
			if err != nil {
				return false, nil, nil, err
			}

			report.Message = tplErrorMatches[messageStringMatch]

			for _, regexInfo := range v.parseErrorsMsgRegexInfo {
				errMsgMatches := regexInfo.errRegexp.FindStringSubmatch(report.Message)
				if len(errMsgMatches) == regexInfo.matchesNumber {
					report.Type = regexInfo.errType
					report.Message = regexInfo.getErrMessage(errMsgMatches)

					break
				}
			}
		}

		return false, report, nil, nil
	}

	var warnings []WrnReport

	// try to find out of block variables in the filtered text
	for _, key := range tplKeys {
		if strings.Contains(v.warningOutOfBlockRegex.ReplaceAllString(s, ""), key) {
			warnings = append(warnings, WrnReport{
				Type:    WrnTypeOutsideBlockVar,
				Message: "Variable is out of a template block",
				Var:     key,
			})
		}
	}

	return true, nil, warnings, nil
}

func getLine(s string) (int, error) {
	locationSplit := strings.Split(s, ":")
	if len(locationSplit) < 2 {
		return 0, fmt.Errorf("template exec error contains invalid location value = %s", s)
	}

	line, err := strconv.Atoi(locationSplit[1])
	if err != nil {
		return 0, fmt.Errorf("convert line variable to int error = %w", err)
	}

	return line, nil
}
