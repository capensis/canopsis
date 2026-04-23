package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
)

const (
	DefaultRegex2MatchTimeout = time.Second
)

type RegexExpression interface {
	Match([]byte) bool
	String() string
}

type WrapperBuiltInRegex struct {
	*regexp.Regexp
}

type WrapperRegex2 struct {
	*regexp2.Regexp
}

func (r WrapperRegex2) Match(content []byte) bool {
	m, e := r.MatchString(string(content))
	if e != nil {
		return false
	}
	return m
}

// NewRegexExpression
// todo move to separate package
func NewRegexExpression(expr string) (RegexExpression, error) {
	if re, err := regexp.Compile(expr); err != nil {
		if re2, err := regexp2.Compile(expr, regexp2.RE2); err != nil {
			return nil, fmt.Errorf("unable to parse regex: %w", err)
		} else {
			re2.MatchTimeout = DefaultRegex2MatchTimeout
			return WrapperRegex2{re2}, nil
		}
	} else {
		return WrapperBuiltInRegex{re}, nil
	}
}

func EscapeRegex(v string) string {
	return regexp2.Escape(v)
}
