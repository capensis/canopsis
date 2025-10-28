package validation

import (
	"fmt"
	"regexp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
)

type duplicateErrorParser struct {
	dupErrorRegexp *regexp.Regexp
	fields         map[string]string
}

type DuplicateErrorParser interface {
	Parse(err error) error
}

func NewDuplicateErrorParser(fields map[string]string) DuplicateErrorParser {
	return &duplicateErrorParser{
		dupErrorRegexp: regexp.MustCompile(`{ ([^:]+)`),
		fields:         fields,
	}
}

func (p *duplicateErrorParser) Parse(err error) error {
	match := p.dupErrorRegexp.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		matchedStr := match[1]

		for k, v := range p.fields {
			if matchedStr == k {
				return common.NewValidationError(k, v)
			}
		}

		return common.NewValidationError(matchedStr, matchedStr+" already exists.")
	}

	return fmt.Errorf("can't parse duplication error: %w", err)
}
