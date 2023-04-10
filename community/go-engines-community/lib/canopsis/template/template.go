package template

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libreflect "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/reflect"
)

const EnvVar = "Env"

type Executor interface {
	Execute(tplStr string, data any) (string, error)
	Parse(tplStr string) (*template.Template, error)
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
	tpl, err := e.Parse(tplStr)
	if err != nil {
		return "", err
	}

	return e.ExecuteByTpl(tpl, data)
}

func (e *executor) Parse(tplStr string) (*template.Template, error) {
	location := e.timezoneConfigProvider.Get().Location
	tpl, err := template.New("tpl").Funcs(GetFunctions(location)).Parse(tplStr)
	if err != nil {
		return nil, err
	}

	tpl.Option("missingkey=error")
	return tpl, err
}

func (e *executor) ExecuteByTpl(tpl *template.Template, data any) (string, error) {
	buf := e.bufPool.Get().(*bytes.Buffer)
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
		"replace": func(oldRegex string, new string, v interface{}) string {
			if s, ok := v.(string); ok {
				re, err := regexp.Compile(oldRegex)
				if err != nil {
					log.Printf("replace : %+v cannot be parsed by regexp, %v", oldRegex, err)
					return ""
				}
				return re.ReplaceAllString(s, new)
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
				timezone = v[1].(string)
				format = v[0].(string)
			} else if len(v) == 2 {
				if value, ok = castTime(v[1]); !ok {
					log.Printf("localtime : %+v is not a CpsTime", v)
					return ""
				}
				format = v[0].(string)
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
		"tag_has_key": func(tags []string, key string) bool {
			for _, tag := range tags {
				split := strings.Split(tag, ":")
				if key == split[0] {
					return true
				}
			}

			return false
		},
		"get_tag": func(tags []string, key string) string {
			for _, tag := range tags {
				split := strings.Split(tag, ": ")
				if key == split[0] {
					if len(split) > 1 {
						return split[1]
					}

					return ""
				}
			}

			return ""
		},
	}
}

func castTime(v interface{}) (time.Time, bool) {
	switch t := v.(type) {
	case types.CpsTime:
		return t.Time, true
	case *types.CpsTime:
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
