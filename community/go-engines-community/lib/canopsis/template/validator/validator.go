package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

const (
	locationStringMatch = 1
	messageStringMatch  = 2
	parseErrorMatches   = 3
)

type RegexpInfo struct {
	errRegexp     *regexp.Regexp
	fullErrRegexp *regexp.Regexp
	// matchesNumber equals number of matched groups in errRegexp + 1
	matchesNumber int
	getErrMessage func([]string) string
}

type Validator interface {
	Validate(s string, data any) (bool, *ErrReport, string, error)
}

type validator struct {
	templateExecutor libtemplate.Executor

	parseErrorRegex         *regexp.Regexp
	parseErrorsMsgRegexInfo []RegexpInfo
}

func NewValidator(templateExecutor libtemplate.Executor) Validator {
	return &validator{
		templateExecutor: templateExecutor,

		parseErrorRegex: regexp.MustCompile("^template: (.+): (.+)$"),
		parseErrorsMsgRegexInfo: []RegexpInfo{
			{
				errRegexp:     regexp.MustCompile("^function \"(.+)\" not defined$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Invalid function \"%s\"", matches[1])
				},
			},
			{
				errRegexp:     regexp.MustCompile("^unexpected \"(.+)\" in (.+)$"),
				matchesNumber: 3,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Unexpected \"%s\"", matches[1])
				},
			},
			{
				errRegexp:     regexp.MustCompile("^unexpected EOF$"),
				matchesNumber: 1,
				getErrMessage: func(_ []string) string {
					return "Parsing error: invalid template"
				},
			},
			{
				errRegexp:     regexp.MustCompile("^unexpected (.+)$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					if matches[1] == "{{end}}" {
						return "Function or block is missing"
					}

					return fmt.Sprintf("Unexpected \"%s\"", matches[1])
				},
			},
			{
				errRegexp:     regexp.MustCompile("^can't evaluate field (.+) in type"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Unknown key or field \"%s\"", matches[1])
				},
			},
			{
				errRegexp:     regexp.MustCompile("no entry for key \"(.+)\"$"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Unknown key or field \"%s\"", matches[1])
				},
			},
			{
				fullErrRegexp: regexp.MustCompile("executing \".+\" at <(.+)>: nil pointer evaluating"),
				matchesNumber: 2,
				getErrMessage: func(matches []string) string {
					return fmt.Sprintf("Undefined key or field \"%s\"", matches[1])
				},
			},
		},
	}
}

type ErrReport struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func (v *validator) Validate(s string, data any) (bool, *ErrReport, string, error) {
	p := v.templateExecutor.Parse(s)
	if tplErr := p.Err; tplErr != nil {
		report, err := v.getReport(tplErr)

		return false, report, "", err
	}

	res, tplErr := v.templateExecutor.ExecuteByTpl(p.Tpl, data)
	if tplErr != nil {
		report, err := v.getReport(tplErr)

		return false, report, "", err
	}

	return true, nil, res, nil
}

func (v *validator) getReport(tplErr error) (*ErrReport, error) {
	fullErrString := tplErr.Error()
	report := &ErrReport{
		Message: fullErrString,
	}

	// parse template error
	tplErrorMatches := v.parseErrorRegex.FindStringSubmatch(fullErrString)
	if len(tplErrorMatches) == parseErrorMatches {
		var err error
		report.Line, err = v.getLine(tplErrorMatches[locationStringMatch])
		if err != nil {
			return nil, err
		}

		report.Message = tplErrorMatches[messageStringMatch]
		for _, regexInfo := range v.parseErrorsMsgRegexInfo {
			if regexInfo.errRegexp != nil {
				errMsgMatches := regexInfo.errRegexp.FindStringSubmatch(report.Message)
				if len(errMsgMatches) == regexInfo.matchesNumber {
					report.Message = regexInfo.getErrMessage(errMsgMatches)

					break
				}
			} else if regexInfo.fullErrRegexp != nil {
				errMsgMatches := regexInfo.fullErrRegexp.FindStringSubmatch(fullErrString)
				if len(errMsgMatches) == regexInfo.matchesNumber {
					report.Message = regexInfo.getErrMessage(errMsgMatches)

					break
				}
			}
		}
	}

	return report, nil
}

func (v *validator) getLine(s string) (int, error) {
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
