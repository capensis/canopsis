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
	SeparatorComma = ","
	SeparatorDot   = "."
	SeparatorSpace = " "
	SeparatorNone  = ""
)

const (
	jsonArray = iota + 1
	separatorArray
)

var errStringIsNotAValidNumber = errors.New("string is not a valid number")
var errStringIsNotAValidBool = errors.New("string is not a valid bool")
var errStringIsNotAValidDatetime = errors.New("string is not a valid datetime")
var errStringIsNotAValidTimestamp = errors.New("string is not a valid timestamp")

type Parser interface {
	Parse(cfg ColumnConfig, initialValue string) (any, error)
}

type parser struct {
	thousandsSeparators map[string]string
	decimalSeparators   map[string]string
	validTrueValues     map[string]bool
	validFalseValues    map[string]bool
	validDatetimeFormat []string
}

func NewParser() Parser {
	return &parser{
		thousandsSeparators: map[string]string{
			"dot":   SeparatorDot,
			"comma": SeparatorComma,
			"space": SeparatorSpace,
			"":      SeparatorNone,
		},
		decimalSeparators: map[string]string{
			"dot":   SeparatorDot,
			"comma": SeparatorComma,
			"":      SeparatorNone,
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

func (p *parser) Parse(cfg ColumnConfig, initialValue string) (any, error) {
	var transformedValue any
	var err error

	switch cfg.Type {
	case externaldata.ColumnTypeString:
		if len(initialValue) > MaxStringLen {
			return nil, fmt.Errorf("string length must be less than %d", MaxStringLen)
		}

		transformedValue = initialValue
	case externaldata.ColumnTypeNumber:
		transformedValue, err = p.parseNumber(initialValue, cfg.ThousandsSeparator, cfg.DecimalSeparator)
	case externaldata.ColumnTypeBoolean:
		transformedValue, err = p.parseBool(initialValue, p.validTrueValues, p.validFalseValues)
	case externaldata.ColumnTypeStringArray:
		if cfg.StringArrayType == separatorArray && cfg.StringArraySeparator == "" {
			return nil, errors.New("string array separator is required")
		}

		transformedValue, err = p.parseStringArray(initialValue, cfg.StringArrayType, cfg.StringArraySeparator)
	case externaldata.ColumnTypeDateTime:
		transformedValue, err = p.parseDatetime(initialValue, p.validDatetimeFormat)
	case externaldata.ColumnTypeTimestamp:
		transformedValue, err = p.parseTimestamp(initialValue)
	default:
		return nil, fmt.Errorf("unexpected column type: %d", cfg.Type)
	}

	return transformedValue, err
}

func (p *parser) parseNumber(stringNumber, thousandsSeparator, decimalSeparator string) (float64, error) {
	if stringNumber == "" || stringNumber == "-" || stringNumber == "-0" {
		return 0, errStringIsNotAValidNumber
	}

	thousandsSeparator, ok := p.thousandsSeparators[thousandsSeparator]
	if !ok {
		return 0, fmt.Errorf("thousands separator must be one of %q, %q, %q or empty", SeparatorComma, SeparatorDot, SeparatorSpace)
	}

	decimalSeparator, ok = p.decimalSeparators[decimalSeparator]
	if !ok {
		return 0, fmt.Errorf("decimal separator must be one of %q, %q or empty", SeparatorDot, SeparatorComma)
	}

	var escapedThousandsSeparator, escapedDecimalSeparator string

	switch thousandsSeparator {
	case SeparatorDot, SeparatorComma, SeparatorSpace:
		escapedThousandsSeparator = regexp.QuoteMeta(thousandsSeparator)
	}

	if decimalSeparator == SeparatorNone {
		escapedDecimalSeparator = regexp.QuoteMeta(SeparatorDot)
	} else {
		escapedDecimalSeparator = regexp.QuoteMeta(decimalSeparator)
	}

	if thousandsSeparator == "" {
		numberRegexp, err := regexp.Compile(`^-?(0|[1-9]\d*)(` + escapedDecimalSeparator + `\d+)?$`)
		if err != nil {
			return 0, fmt.Errorf("failed to compile regexp: %w", err)
		}

		if !numberRegexp.MatchString(stringNumber) {
			return 0, errStringIsNotAValidNumber
		}
	} else {
		numberRegexp, err := regexp.Compile(`^-?(0|[1-9]\d*|[1-9]\d{0,2}(` + escapedThousandsSeparator + `\d{3})+)(` + escapedDecimalSeparator + `\d+)?$`)
		if err != nil {
			return 0, fmt.Errorf("failed to compile regexp: %w", err)
		}

		if !numberRegexp.MatchString(stringNumber) {
			return 0, errStringIsNotAValidNumber
		}

		stringNumber = strings.ReplaceAll(stringNumber, thousandsSeparator, "")
	}

	if decimalSeparator != SeparatorNone && decimalSeparator != SeparatorDot {
		stringNumber = strings.ReplaceAll(stringNumber, decimalSeparator, SeparatorDot)
	}

	number, err := strconv.ParseFloat(stringNumber, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number: %w", err)
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

func (p *parser) parseStringArray(stringStringArray string, arrayType int, separator string) ([]string, error) {
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
	case separatorArray:
		if stringStringArray == "" {
			return []string{}, nil
		}

		if separator == "" {
			return nil, errors.New("separator cannot be empty for separator array type")
		}

		return strings.Split(stringStringArray, separator), nil
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
