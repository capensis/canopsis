package utils

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

const NamingCharacterSet = "abcdefghijklmnopqrstuvwxyz1234567890"

var NumberOfCharacter = int64(len(NamingCharacterSet))

// NewID generate an uuid
func NewID() string {
	id, err := uuid.NewV7()
	if err != nil {
		// error is extremely rare so panic is ok
		panic(fmt.Errorf("cannot generate new ID: %w", err))
	}

	return id.String()
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
	if chars < 1 || chars >= len(s) || chars >= utf8.RuneCountInString(s) {
		return s
	}

	// Check if the string is ASCII
	if len(s) == utf8.RuneCountInString(s) {
		return s[:chars]
	}

	// Handle multi-byte characters
	truncated := make([]rune, 0, chars)
	for _, r := range s {
		if len(truncated) >= chars {
			break
		}
		truncated = append(truncated, r)
	}

	return string(truncated)
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

// RandBase64String generates a random base64 string by generated n bytes
func RandBase64String(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ToObjectIDHex(s string) string {
	const hexLen = 24
	re := regexp.MustCompile(fmt.Sprintf(`^[0-9a-fA-F]{%d}$`, hexLen))
	if len(s) != hexLen || !re.MatchString(s) {
		hash := sha3.New384()
		hash.Write([]byte(s))
		s = hex.EncodeToString(hash.Sum(nil)[:hexLen/2])
	}
	return s
}
