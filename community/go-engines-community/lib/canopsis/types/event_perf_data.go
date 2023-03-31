package types

import (
	"strconv"
	"strings"
)

const (
	equalRune        = '='
	singleQuoteRune  = '\''
	dotRune          = '.'
	spaceRune        = ' '
	zeroRune         = '0'
	nineRune         = '9'
	undefinedValRune = 'U'
	invalidIndex     = -1
)

type PerfData struct {
	Name  string
	Value float64
	Unit  string
}

func (e *Event) GetPerfData() []PerfData {
	if e.PerfData == "" {
		return nil
	}

	parsed := make([]PerfData, 0)
	i := 0
	l := len(e.PerfData)
	for {
		d, lastIndex := parseNextPerfData(e.PerfData[i:])
		if lastIndex < 0 {
			return nil
		}
		if d.Name != "" {
			parsed = append(parsed, d)
		}

		i += lastIndex + 1
		if i == l {
			break
		}
	}

	return parsed
}

func parseNextPerfData(s string) (PerfData, int) {
	data := PerfData{}
	lastIndex := invalidIndex
	l := len(s)

	name, i := parsePerfDataName(s)
	if i < 0 || i == l-1 {
		return data, lastIndex
	}

	val, unit, ok, j := parsePerfDataValue(s[i+1:])
	if j < 0 {
		return data, lastIndex
	}
	lastIndex = i + j + 1
	if !ok {
		return data, lastIndex
	}

	if unit != "" {
		name += "_" + unit
	}

	data = PerfData{
		Name:  name,
		Value: val,
		Unit:  unit,
	}

	return data, lastIndex
}

func parsePerfDataName(s string) (string, int) {
	lastIndex := strings.IndexRune(s, equalRune)
	if lastIndex < 0 {
		return "", lastIndex
	}
	name := s[:lastIndex]
	if name == "" {
		return "", invalidIndex
	}

	l := len(name)
	if name[0] == singleQuoteRune && name[l-1] == singleQuoteRune {
		// Unquote
		name = name[1 : l-1]
		l = len(name)
	} else if i := strings.IndexRune(name, spaceRune); i >= 0 {
		return "", invalidIndex
	}

	// Check escaped quotes
	for i := 0; i < l; {
		if name[i] == singleQuoteRune {
			if i == l-1 || name[i+1] != singleQuoteRune {
				return "", invalidIndex
			}
			i += 2
		} else {
			i++
		}
	}
	// Unescape quotes
	name = strings.ReplaceAll(name, "''", "'")

	return name, lastIndex
}

func parsePerfDataValue(s string) (float64, string, bool, int) {
	lastIndex := strings.IndexRune(s, spaceRune)
	paramsStr := ""
	if lastIndex < 0 {
		paramsStr = s
		lastIndex = len(s) - 1
	} else {
		paramsStr = s[:lastIndex]
	}

	params := strings.Split(paramsStr, ";")
	if len(params) == 0 {
		return 0, "", false, invalidIndex
	}

	valWithUnit := params[0]
	if len(valWithUnit) == 0 {
		return 0, "", false, invalidIndex
	}

	if valWithUnit[0] == undefinedValRune {
		return 0, "", false, lastIndex
	}

	notDigitFilter := func(r rune) bool {
		return r < zeroRune || r > nineRune
	}
	i := strings.IndexFunc(valWithUnit, notDigitFilter)
	valStr := ""
	unit := ""
	if i < 0 {
		valStr = valWithUnit
	} else if valWithUnit[i] == dotRune {
		if i == len(valWithUnit)-1 {
			return 0, "", false, invalidIndex
		}

		j := strings.IndexFunc(valWithUnit[i+1:], notDigitFilter)
		if j < 0 {
			valStr = valWithUnit
		} else {
			valStr = valWithUnit[:i+j+1]
			unit = valWithUnit[i+j+1:]
		}
	} else {
		valStr = valWithUnit[:i]
		unit = valWithUnit[i:]
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, "", false, invalidIndex
	}

	return val, unit, true, lastIndex
}
