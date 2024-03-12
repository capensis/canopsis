package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
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

// FindAllStringSubmatchMapWithRegexExpression returns a slice of maps containing the values of the named
// subexpressions of all matches of the regular expression re in s. A
// return value of nil indicates no match.
func FindAllStringSubmatchMapWithRegexExpression(re RegexExpression, s string) []map[string]string {
	switch regex := re.(type) {
	case WrapperBuiltInRegex:
		matches := regex.FindAllStringSubmatch(s, -1)
		if matches == nil {
			return nil
		}

		names := regex.SubexpNames()
		submatches := make([]map[string]string, 0, len(names))
		for i, match := range matches {
			unnamedCount := 1
			for j := 1; j < len(names); j++ {
				if len(submatches) == i {
					submatches = append(submatches, make(map[string]string))
				}
				if names[j] == "" {
					submatches[i][strconv.Itoa(unnamedCount)] = match[j]
					unnamedCount++
				} else {
					submatches[i][names[j]] = match[j]
				}
			}
		}
		return submatches
	case WrapperRegex2:
		match, err := regex.FindStringMatch(s)
		if err != nil || match == nil {
			return nil
		}
		submatches := make([]map[string]string, 0, match.GroupCount())

		for i := 0; err == nil && match != nil; i++ {
			names := regex.GetGroupNames()
			if len(names) > 1 {
				submatches = append(submatches, make(map[string]string, len(names)-1))
				for _, name := range names[1:] {
					submatches[i][name] = match.GroupByName(name).String()
				}
			}
			match, err = regex.FindNextMatch(match)
		}
		return submatches
	}

	return nil
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
