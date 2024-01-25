package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

const NamingCharacterSet = "abcdefghijklmnopqrstuvwxyz1234567890"

var NumberOfCharacter = int64(len(NamingCharacterSet))

// GetStringField read a string field in a struct, with introspection
func GetStringField(a interface{}, field string) string {
	r := reflect.ValueOf(a)
	f := reflect.Indirect(r).FieldByName(field)
	value := f.String()
	if value == "<invalid Value>" {
		// zero-ing invalid values
		value = ""
	}
	return value
}

// NewID generate an uuid
func NewID() string {
	return uuid.New().String()
}

func MatchWithRegexExpression(re RegexExpression, s string) bool {
	switch uthRegex := re.(type) {
	case WrapperBuiltInRegex:
		return uthRegex.MatchString(s)
	case WrapperRegex2:
		match, _ := uthRegex.MatchString(s)
		return match
	}

	return false
}

// FindStringSubmatchMapWithRegexExpression returns a map containing the values of the named
// subexpressions of the lefmost match of the regular expression re in s. A
// return value of nil indicates no match.
func FindStringSubmatchMapWithRegexExpression(re RegexExpression, s string) map[string]string {
	switch uthRegex := re.(type) {
	case WrapperBuiltInRegex:
		match := uthRegex.FindStringSubmatch(s)
		if match == nil {
			return nil
		}

		names := uthRegex.SubexpNames()
		submatches := make(map[string]string)
		unnamedCount := 1
		for i := 1; i < len(names); i++ {
			if names[i] == "" {
				submatches[strconv.Itoa(unnamedCount)] = match[i]
				unnamedCount++
			} else {
				submatches[names[i]] = match[i]
			}
		}
		return submatches
	case WrapperRegex2:
		match, err := uthRegex.FindStringMatch(s)
		if err != nil || match == nil {
			return nil
		}
		names := uthRegex.GetGroupNames()
		submatches := make(map[string]string)
		for i := 1; i < len(names); i++ {
			submatches[names[i]] = match.GroupByName(names[i]).String()
		}
		return submatches
	}

	return nil
}

// FindStringSubmatchMap returns a map containing the values of the named
// subexpressions of the lefmost match of the regular expression re in s. A
// return value of nil indicates no match.
func FindStringSubmatchMap(re *regexp.Regexp, s string) map[string]string {
	match := re.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	names := re.SubexpNames()
	submatches := make(map[string]string)
	for i := 1; i < len(names); i++ {
		submatches[names[i]] = match[i]
	}

	return submatches
}

// AsString tries to convert an interface{} into a string, and returns its
// value and an integer indicating whether it succeeded or not.
func AsString(value interface{}) (string, bool) {
	switch typedValue := value.(type) {
	case string:
		return typedValue, true
	case *string:
		if typedValue == nil {
			return "", false
		}
		return *typedValue, true
	default:
		return "", false
	}
}

func TruncateString(s string, chars int) string {
	if chars < 1 {
		return s
	}

	format := fmt.Sprintf("%%.%ds", chars)

	return fmt.Sprintf(format, s)
}

// RandString generate a random string
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		index, _ := crand.Int(crand.Reader, big.NewInt(NumberOfCharacter))
		b[i] = NamingCharacterSet[index.Int64()]
	}
	return string(b)
}
