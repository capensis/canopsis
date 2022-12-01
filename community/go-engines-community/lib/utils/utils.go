package utils

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/google/uuid"
)

const NamingCharacterSet = "abcdefghijklmnopqrstuvwxyz1234567890"

var NumberOfCharacter = int64(len(NamingCharacterSet))

// ApplyTranslation translate a templated string using informations in an Event
func ApplyTranslation(element interface{}, key, formatString string) (string, error) {
	t, err := template.New(key).Parse(formatString)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, element)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReadFile handles opening and reading file at given path
func ReadFile(path string) ([]byte, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("read file: %v", err)
	}
	return source, nil
}

// StringInSlice checks if a string is present in a list of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}

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

// FindStringSubmatchMap returns a map containing the values of the named
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

func TimeTrack(start time.Time, name string) {
	fmt.Printf("%s took %s\n", name, time.Since(start))
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
