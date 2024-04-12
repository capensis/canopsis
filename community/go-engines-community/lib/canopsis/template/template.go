package template

//go:generate mockgen -destination=../../../mocks/lib/canopsis/template/template.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template Executor

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libreflect "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/reflect"
)

const EnvVar = "Env"

var ErrFailedConvertToInt64 = errors.New("failed convert to int64")
var ErrDivisionByZero = errors.New("division by zero")

type ParsedTemplate struct {
	Text string
	Tpl  *template.Template
	Err  error
}

type Executor interface {
	Execute(tplStr string, data any) (string, error)
	Parse(text string) ParsedTemplate
	ExecuteByTpl(tpl *template.Template, data any) (string, error)
}

type executor struct {
	templateConfigProvider config.TemplateConfigProvider
	timezoneConfigProvider config.TimezoneConfigProvider
	bufPool                sync.Pool
}

func NewExecutor(
	templateConfigProvider config.TemplateConfigProvider,
	timezoneConfigProvider config.TimezoneConfigProvider,
) Executor {
	return &executor{
		templateConfigProvider: templateConfigProvider,
		timezoneConfigProvider: timezoneConfigProvider,
		bufPool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (e *executor) Execute(tplStr string, data any) (string, error) {
	if tplStr == "" {
		return "", nil
	}

	tpl := e.Parse(tplStr)
	if tpl.Err != nil {
		return "", tpl.Err
	}

	return e.ExecuteByTpl(tpl.Tpl, data)
}

func (e *executor) Parse(text string) ParsedTemplate {
	if text == "" {
		return ParsedTemplate{}
	}

	location := e.timezoneConfigProvider.Get().Location
	tpl, err := template.New("tpl").Funcs(GetFunctions(location)).Parse(text)
	if err != nil {
		return ParsedTemplate{
			Text: text,
			Err:  err,
		}
	}

	tpl.Option("missingkey=error")
	return ParsedTemplate{
		Text: text,
		Tpl:  tpl,
	}
}

func (e *executor) ExecuteByTpl(tpl *template.Template, data any) (string, error) {
	buf, ok := e.bufPool.Get().(*bytes.Buffer)
	if !ok {
		return "", errors.New("unknown buffer type")
	}

	defer e.bufPool.Put(buf)
	buf.Reset()
	err := tpl.Execute(buf, addEnvVarsToData(data, e.templateConfigProvider.Get().Vars))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetFunctions(appLocation *time.Location) template.FuncMap {
	return template.FuncMap{
		// json will convert an item to an JSON-compatible element,
		// ie ints will be returned as integers and strings returned as strings with quotes
		// It is mostly used for the strings to preserve their content
		// and avoid Go behavior to escape special characters strings by default
		"json": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return err.Error()
			}

			return string(b)
		},
		// Same behavior as json but remove the quotes around the string.
		"json_unquote": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return err.Error()
			}

			jsonStr := string(b)

			if string(jsonStr[0]) != "\"" || string(jsonStr[len(jsonStr)-1]) != "\"" {
				log.Printf("json_unquote : trying to unquote %+v. Returning value directly", jsonStr)
				return jsonStr
			}

			return jsonStr[1 : len(jsonStr)-1]
		},
		// split will split a string according a separator and returns the substring
		"split": func(sep string, index int, v interface{}) string {
			if s, ok := v.(string); ok {
				stringSlice := strings.Split(s, sep)
				if 0 <= index && index < len(stringSlice) {
					return stringSlice[index]
				}
				log.Printf("split : index %+v out of bounds", index)
			} else {
				log.Printf("split : %+v is not a string", v)
			}
			return ""
		},
		// trim will return a string with all leading and trailing white space removed
		"trim": func(v interface{}) string {
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s)
			}
			log.Printf("trim : %+v is not a string", v)
			return ""
		},
		// formattedDate will return a formatted string from a time type
		"formattedDate": func(format string, v interface{}) string {
			if t, ok := castTime(v); ok {
				return t.Format(format)
			}
			log.Printf("formattedDate : %+v is not a time type", v)
			return ""
		},
		// replace will replace a string, replacing matches of the regex with the replacement string
		"replace": func(oldRegex string, newV string, v interface{}) string {
			if s, ok := v.(string); ok {
				re, err := regexp.Compile(oldRegex)
				if err != nil {
					log.Printf("replace : %+v cannot be parsed by regexp, %v", oldRegex, err)
					return ""
				}
				return re.ReplaceAllString(s, newV)
			}
			log.Printf("replace : %+v is not a string", v)
			return ""
		},
		// upper string
		"uppercase": func(v interface{}) string {
			if s, ok := v.(string); ok {
				return strings.ToUpper(s)
			}
			log.Printf("trim : %+v is not a string", v)
			return ""
		},

		// upper string
		"lowercase": func(v interface{}) string {
			if s, ok := v.(string); ok {
				return strings.ToLower(s)
			}
			log.Printf("trim : %+v is not a string", v)
			return ""
		},
		"localtime": func(v ...interface{}) string {
			var value time.Time
			var timezone string
			var format string
			var ok bool

			if len(v) == 3 {
				if value, ok = castTime(v[2]); !ok {
					log.Printf("localtime : %+v is not a CpsTime", v)
					return ""
				}

				if timezone, ok = v[1].(string); !ok {
					log.Printf("localtime : %+v is not a string", v[1])
					return ""
				}

				if format, ok = v[0].(string); !ok {
					log.Printf("localtime : %+v is not a string", v[0])
					return ""
				}
			} else if len(v) == 2 {
				if value, ok = castTime(v[1]); !ok {
					log.Printf("localtime : %+v is not a CpsTime", v)
					return ""
				}

				if format, ok = v[0].(string); !ok {
					log.Printf("localtime : %+v is not a string", v[0])
					return ""
				}
			} else {
				log.Print("localtime : must have 1 or 2 arguments")
				return ""
			}

			var loc *time.Location
			if timezone != "" {
				var err error
				if loc, err = time.LoadLocation(timezone); err != nil {
					log.Print("localtime : invalid timezone")
					return ""
				}
			} else if appLocation != nil {
				loc = appLocation
			}

			if loc == nil {
				return value.Format(format)
			}

			return value.In(loc).Format(format)
		},
		// regex_map_key return map value if key match the regexp
		"regex_map_key": func(m map[string]interface{}, regexpString string) interface{} {
			re, err := regexp.Compile(regexpString)
			if err != nil {
				log.Printf("regex_map_key : failed to compile regexp %s, %v", regexpString, err)
				return ""
			}

			for k, v := range m {
				if re.MatchString(k) {
					return v
				}
			}

			return ""
		},
		"map_has_key": func(m any, key any) bool {
			rv := reflect.ValueOf(m)
			if !rv.IsValid() || rv.Kind() != reflect.Map {
				return false
			}

			return rv.MapIndex(reflect.ValueOf(key)).IsValid()
		},
		"tag_has_key": func(tags []string, searchKey string) bool {
			for _, tag := range tags {
				key, _, _ := strings.Cut(tag, ":")
				if searchKey == key {
					return true
				}
			}

			return false
		},
		"get_tag": func(tags []string, searchKey string) string {
			for _, tag := range tags {
				key, val, _ := strings.Cut(tag, ": ")
				if searchKey == key {
					return val
				}
			}

			return ""
		},
		"substrLeft": func(str string, n int64) string {
			runeStr := []rune(str)

			if n < 1 {
				return ""
			}

			if n >= int64(len(runeStr)) {
				return str
			}

			return string(runeStr[:n])
		},
		"substrRight": func(str string, n int64) string {
			runeStr := []rune(str)
			lenRuneStr := int64(len(runeStr))

			if n < 1 {
				return ""
			}

			if n >= lenRuneStr {
				return str
			}

			return string(runeStr[lenRuneStr-n:])
		},
		"substr": func(str string, start, n int64) string {
			runeStr := []rune(str)
			lenRuneStr := int64(len(runeStr))

			if start < 0 || start >= lenRuneStr || n < 1 {
				return ""
			}

			end := start + n

			if end >= lenRuneStr {
				return string(runeStr[start:])
			}

			return string(runeStr[start:end])
		},
		"strlen": func(str string) int64 {
			return int64(len([]rune(str)))
		},
		"strpos": func(str, substr string) int64 {
			idx := strings.Index(str, substr)
			if idx < 0 {
				return -1
			}

			return int64(len([]rune(str[:idx])))
		},
		"add": func(a, b any) (int64, error) {
			return wrapInt64ArithmeticFunc(a, b, func(x, y int64) (int64, error) {
				return x + y, nil
			})
		},
		"sub": func(a, b any) (int64, error) {
			return wrapInt64ArithmeticFunc(a, b, func(x, y int64) (int64, error) {
				return x - y, nil
			})
		},
		"mult": func(a, b any) (int64, error) {
			return wrapInt64ArithmeticFunc(a, b, func(x, y int64) (int64, error) {
				return x * y, nil
			})
		},
		"div": func(a, b any) (int64, error) {
			return wrapInt64ArithmeticFunc(a, b, func(x, y int64) (int64, error) {
				if y == 0 {
					return 0, ErrDivisionByZero
				}

				return x / y, nil
			})
		},
	}
}

func wrapInt64ArithmeticFunc(a, b any, f func(x, y int64) (int64, error)) (int64, error) {
	aInt64, ok := toInt64(a)
	if !ok {
		return 0, ErrFailedConvertToInt64
	}

	bInt64, ok := toInt64(b)
	if !ok {
		return 0, ErrFailedConvertToInt64
	}

	return f(aInt64, bInt64)
}

func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int8:
		return int64(val), true
	case int16:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true
	case uint:
		return int64(val), true
	case uint8:
		return int64(val), true
	case uint16:
		return int64(val), true
	case uint32:
		return int64(val), true
	case uint64:
		return int64(val), true
	default:
		return 0, false
	}
}

func castTime(v interface{}) (time.Time, bool) {
	switch t := v.(type) {
	case datetime.CpsTime:
		return t.Time, true
	case *datetime.CpsTime:
		if t == nil {
			return time.Time{}, false
		}
		return t.Time, true
	case time.Time:
		return t, true
	case *time.Time:
		if t == nil {
			return time.Time{}, false
		}
		return *t, true
	default:
		return time.Time{}, false
	}
}

func addEnvVarsToData(data any, envVars map[string]any) any {
	if len(envVars) == 0 {
		return data
	}

	mapData, ok := libreflect.ToMap(data)
	if !ok {
		return data
	}

	mapData[EnvVar] = envVars
	return mapData
}
