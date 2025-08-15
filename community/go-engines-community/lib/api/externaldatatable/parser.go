package externaldatatable

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"github.com/valyala/fastjson"
)

const (
	DelimiterComma = ","
	DelimiterDot   = "."
	DelimiterSpace = " "
	DelimiterNone  = ""
)

const (
	jsonArray = iota + 1
	delimiterArray
)

var errStringIsNotAValidNumber = errors.New("string is not a valid number")
var errStringIsNotAValidBool = errors.New("string is not a valid bool")
var errStringIsNotAValidDatetime = errors.New("string is not a valid datetime")
var errStringIsNotAValidTimestamp = errors.New("string is not a valid timestamp")

type Parser interface {
	Parse(cfg externaldata.ColumnConfig, initialValue string) (any, error)
}

type parser struct {
	thousandsDelimiters map[string]string
	decimalDelimiters   map[string]string
	validTrueValues     map[string]bool
	validFalseValues    map[string]bool
	validDatetimeFormat []string
}

func NewParser() Parser {
	return &parser{
		thousandsDelimiters: map[string]string{
			"dot":   DelimiterDot,
			"comma": DelimiterComma,
			"space": DelimiterSpace,
			"":      DelimiterNone,
		},
		decimalDelimiters: map[string]string{
			"dot":   DelimiterDot,
			"comma": DelimiterComma,
			"":      DelimiterNone,
		},
		validTrueValues: map[string]bool{
			"yes":  true,
			"y":    true,
			"oui":  true,
			"true": true,
			"1":    true,
		},
		validFalseValues: map[string]bool{
			"no":    true,
			"n":     true,
			"non":   true,
			"false": true,
			"0":     true,
		},
		validDatetimeFormat: []string{
			"2006-01-02T15:04:05.00Z",
			"2006-01-02T15:04:05-07:00",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.00-07:00",
		},
	}
}

func (p *parser) Parse(cfg externaldata.ColumnConfig, initialValue string) (any, error) {
	var transformedValue any
	var err error

	switch cfg.Type {
	case externaldata.ColumnTypeString:
		if len(initialValue) > externaldata.MaxStringLen {
			return nil, fmt.Errorf("string length must be less than %d", externaldata.MaxStringLen)
		}

		transformedValue = initialValue
	case externaldata.ColumnTypeNumber:
		transformedValue, err = p.parseNumber(initialValue, cfg.ThousandsDelimiter, cfg.DecimalDelimiter)
	case externaldata.ColumnTypeBoolean:
		transformedValue, err = p.parseBool(initialValue, p.validTrueValues, p.validFalseValues)
	case externaldata.ColumnTypeStringArray:
		if cfg.StringArrayType == delimiterArray && cfg.StringArrayDelimiter == "" {
			return nil, errors.New("string array delimiter is required")
		}

		transformedValue, err = p.parseStringArray(initialValue, cfg.StringArrayType, cfg.StringArrayDelimiter)
	case externaldata.ColumnTypeDateTime:
		transformedValue, err = p.parseDatetime(initialValue, p.validDatetimeFormat)
	case externaldata.ColumnTypeTimestamp:
		transformedValue, err = p.parseTimestamp(initialValue)
	default:
		return nil, fmt.Errorf("unexpected column type: %d", cfg.Type)
	}

	return transformedValue, err
}

func (p *parser) parseNumber(stringNumber, thousandsDelimiter, decimalDelimiter string) (float64, error) {
	if stringNumber == "" {
		return 0, errStringIsNotAValidNumber
	}

	isNegative := false
	if strings.HasPrefix(stringNumber, "-") {
		isNegative = true
		stringNumber = stringNumber[1:]

		// edge case if input stringNumber = "-"
		if len(stringNumber) == 0 {
			return 0, errStringIsNotAValidNumber
		}
	}

	// edge case if input stringNumber = "-0"
	if isNegative && stringNumber == "0" {
		return 0, errStringIsNotAValidNumber
	}

	thousandsDelimiter, ok := p.thousandsDelimiters[thousandsDelimiter]
	if !ok {
		return 0, fmt.Errorf("thousands delimiter must be one of %q, %q, %q or empty", DelimiterComma, DelimiterDot, DelimiterSpace)
	}

	decimalDelimiter, ok = p.decimalDelimiters[decimalDelimiter]
	if !ok {
		return 0, fmt.Errorf("decimal delimiter must be one of %q, %q or empty", DelimiterDot, DelimiterComma)
	}

	if len(stringNumber) > 1 && stringNumber[0] == '0' {
		if decimalDelimiter != DelimiterNone {
			if !strings.HasPrefix(stringNumber, "0"+decimalDelimiter) {
				return 0, errStringIsNotAValidNumber
			}
		} else {
			if !strings.HasPrefix(stringNumber, "0.") {
				return 0, errStringIsNotAValidNumber
			}
		}
	}

	var escapedThousandsDelimiter, escapedDecimalDelimiter string

	switch thousandsDelimiter {
	case DelimiterDot, DelimiterComma, DelimiterSpace:
		escapedThousandsDelimiter = regexp.QuoteMeta(thousandsDelimiter)
	}

	if decimalDelimiter == DelimiterNone {
		escapedDecimalDelimiter = regexp.QuoteMeta(DelimiterDot)
	} else {
		escapedDecimalDelimiter = regexp.QuoteMeta(decimalDelimiter)
	}

	if thousandsDelimiter == "" {
		numberRegexp, err := regexp.Compile(`^(0|[1-9]\d*)(` + escapedDecimalDelimiter + `\d+)?$`)
		if err != nil {
			return 0, fmt.Errorf("failed to compile regexp: %w", err)
		}

		if !numberRegexp.MatchString(stringNumber) {
			return 0, errStringIsNotAValidNumber
		}

		if decimalDelimiter != DelimiterNone && decimalDelimiter != DelimiterDot {
			stringNumber = strings.ReplaceAll(stringNumber, decimalDelimiter, DelimiterDot)
		}
	} else {
		numberRegexp, err := regexp.Compile(`^(0|[1-9]\d*|\d{1,3}(` + escapedThousandsDelimiter + `\d{3})+)(` + escapedDecimalDelimiter + `\d+)?$`)
		if err != nil {
			return 0, fmt.Errorf("failed to compile regexp: %w", err)
		}

		if !numberRegexp.MatchString(stringNumber) {
			return 0, errStringIsNotAValidNumber
		}

		stringNumber = strings.ReplaceAll(stringNumber, thousandsDelimiter, "")
		if decimalDelimiter != DelimiterNone && decimalDelimiter != DelimiterDot {
			stringNumber = strings.ReplaceAll(stringNumber, decimalDelimiter, DelimiterDot)
		}
	}

	number, err := strconv.ParseFloat(stringNumber, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number: %w", err)
	}

	if isNegative {
		number = -number
	}

	return number, nil
}

func (p *parser) parseBool(stringBool string, trueValues, falseValues map[string]bool) (bool, error) {
	if stringBool == "" {
		return false, errStringIsNotAValidBool
	}

	stringBool = strings.ToLower(stringBool)

	if trueValues[stringBool] {
		return true, nil
	}

	if falseValues[stringBool] {
		return false, nil
	}

	return false, errStringIsNotAValidBool
}

func (p *parser) parseStringArray(stringStringArray string, arrayType int, delimiter string) ([]string, error) {
	switch arrayType {
	case jsonArray:
		if stringStringArray == "" {
			return nil, errors.New("empty string is not a valid JSON array")
		}

		var p fastjson.Parser
		v, err := p.Parse(stringStringArray)
		if err != nil {
			return nil, fmt.Errorf("string is not a valid JSON: %w", err)
		}

		if v.Type() != fastjson.TypeArray {
			return nil, errors.New("string is not a JSON array")
		}

		arr, _ := v.Array()

		result := make([]string, len(arr))
		for i, item := range arr {
			if item.Type() != fastjson.TypeString {
				return nil, fmt.Errorf("string array element at index %d is not a string", i)
			}

			result[i] = string(item.GetStringBytes())
		}

		return result, nil
	case delimiterArray:
		if stringStringArray == "" {
			return []string{}, nil
		}

		if delimiter == "" {
			return nil, errors.New("delimiter cannot be empty for delimiter array type")
		}

		return strings.Split(stringStringArray, delimiter), nil
	default:
		return nil, fmt.Errorf("invalid array type: %d", arrayType)
	}
}

func (p *parser) parseDatetime(stringDatetime string, validFormats []string) (int64, error) {
	if stringDatetime == "" {
		return 0, errStringIsNotAValidDatetime
	}

	for _, format := range validFormats {
		if t, err := time.Parse(format, stringDatetime); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, errStringIsNotAValidDatetime
}

func (p *parser) parseTimestamp(stringTimestamp string) (int64, error) {
	if stringTimestamp == "" {
		return 0, errStringIsNotAValidTimestamp
	}

	if strings.HasPrefix(stringTimestamp, "+") {
		return 0, errStringIsNotAValidTimestamp
	}

	if len(stringTimestamp) > 1 && stringTimestamp[0] == '0' {
		return 0, errStringIsNotAValidTimestamp
	}

	timestamp, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return 0, errStringIsNotAValidTimestamp
	}

	if timestamp < 0 {
		return 0, errStringIsNotAValidTimestamp
	}

	return timestamp, nil
}
