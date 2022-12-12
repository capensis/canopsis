package bdd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"
	"time"

	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	tpl, err := template.New("tpl").
		Option("missingkey=error").
		Funcs(getTplFuncs()).
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

func getTplFuncs() template.FuncMap {
	return template.FuncMap{
		"now": func() int64 {
			return time.Now().Unix()
		},
		"nowAdd": func(s string) (int64, error) {
			d, err := libtypes.ParseDurationWithUnit(s)
			if err != nil {
				return 0, err
			}

			return d.AddTo(libtypes.NewCpsTime()).Unix(), nil
		},
		"nowDate": func() int64 {
			y, m, d := time.Now().UTC().Date()

			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
		},
		"nowDateAdd": func(s string) (int64, error) {
			d, err := libtypes.ParseDurationWithUnit(s)
			if err != nil {
				return 0, err
			}

			year, month, day := time.Now().UTC().Date()
			now := libtypes.CpsTime{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}

			return d.AddTo(now).Unix(), nil
		},
		"parseTime": func(s string) (int64, error) {
			t, err := time.ParseInLocation("02-01-2006 15:04", s, time.UTC)
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
	}
}
