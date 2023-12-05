package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type Templater struct {
	defaultVars map[string]interface{}
}

func NewTemplater(defaultVars map[string]interface{}) *Templater {
	return &Templater{
		defaultVars: defaultVars,
	}
}

func (t *Templater) Execute(ctx context.Context, text string) (*bytes.Buffer, error) {
	loc, _ := GetTimezone(ctx)
	tpl, err := template.New("tpl").
		Option("missingkey=error").
		Funcs(getTplFuncs(loc)).
		Parse(text)
	if err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	responseBody, _ := getResponseBody(ctx)
	vars, _ := getVars(ctx)
	data := make(map[string]interface{}, len(t.defaultVars)+len(vars)+1)
	for k, v := range t.defaultVars {
		data[k] = v
	}
	for k, v := range vars {
		data[k] = v
	}
	data["lastResponse"] = responseBody

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("cannot execute template: %w", err)
	}

	return buf, nil
}

func getTplFuncs(location *time.Location) template.FuncMap {
	return template.FuncMap{
		// json converts an item to an JSON-compatible element.
		// For the strings it escapes newline and quote chars
		"json": func(v any) string {
			sv := struct {
				V any `json:"v"`
			}{V: v}
			b, err := json.Marshal(sv)
			if err != nil {
				return err.Error()
			}

			return string(b[6 : len(b)-2])
		},
		"now": func() int64 {
			return time.Now().Unix()
		},
		"nowTz": func() int64 {
			return time.Now().In(location).Unix()
		},
		"nowAdd": func(s string) (int64, error) {
			d, err := datetime.ParseDurationWithUnit(s)
			if err != nil {
				return 0, err
			}

			return d.AddTo(datetime.NewCpsTime()).Unix(), nil
		},
		"nowDate": func() int64 {
			y, m, d := time.Now().UTC().Date()

			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
		},
		"nowDateTz": func() int64 {
			y, m, d := time.Now().In(location).Date()

			return time.Date(y, m, d, 0, 0, 0, 0, location).Unix()
		},
		"nowDateAdd": func(s string) (int64, error) {
			d, err := datetime.ParseDurationWithUnit(s)
			if err != nil {
				return 0, err
			}

			year, month, day := time.Now().UTC().Date()
			now := datetime.CpsTime{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}

			return d.AddTo(now).Unix(), nil
		},
		"parseTime": func(s string) (int64, error) {
			t, err := time.ParseInLocation("02-01-2006 15:04", s, time.UTC)
			if err != nil {
				return 0, err
			}

			return t.Unix(), nil
		},
		"parseTimeTz": func(s string) (int64, error) {
			t, err := time.ParseInLocation("02-01-2006 15:04", s, location)
			if err != nil {
				return 0, err
			}

			return t.Unix(), nil
		},
		"sumTime": func(args ...interface{}) (int64, error) {
			var sum int64
			for _, arg := range args {
				switch v := arg.(type) {
				case string:
					i, err := strconv.Atoi(v)
					if err != nil {
						return 0, err
					}

					sum += int64(i)
				case int:
					sum += int64(v)
				case int64:
					sum += v
				default:
					return 0, fmt.Errorf("unexpected type %T of argument %v", arg, arg)
				}
			}

			return sum, nil
		},
		"query_escape": func(s string) string {
			return url.QueryEscape(s)
		},
	}
}
