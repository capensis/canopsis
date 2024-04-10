package request

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

type Parameters struct {
	URL        string                     `bson:"url" json:"url" binding:"required"`
	Method     string                     `bson:"method" json:"method" binding:"required,oneof=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	Auth       *BasicAuth                 `bson:"auth,omitempty" json:"auth"`
	Headers    map[string]string          `bson:"headers,omitempty" json:"headers"`
	Payload    string                     `bson:"payload,omitempty" json:"payload"`
	SkipVerify bool                       `bson:"skip_verify" json:"skip_verify"`
	Timeout    *datetime.DurationWithUnit `bson:"timeout,omitempty" json:"timeout"`
	RetryCount int64                      `bson:"retry_count,omitempty" json:"retry_count"`
	RetryDelay *datetime.DurationWithUnit `bson:"retry_delay,omitempty" json:"retry_delay"`
}

type BasicAuth struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type ParsedParameters struct {
	URL        template.ParsedTemplate
	Method     string
	Auth       *BasicAuth
	Headers    map[string]string
	Payload    template.ParsedTemplate
	SkipVerify bool
	Timeout    *datetime.DurationWithUnit
	RetryCount int64
	RetryDelay *datetime.DurationWithUnit
}

type TemplateErr struct {
	field string
	err   error
}

func (e TemplateErr) Error() string {
	return e.err.Error()
}

func (e TemplateErr) Field() string {
	return e.field
}

func ExecuteTemplates(
	params Parameters,
	templateExecutor template.Executor,
	tplData any,
) (Parameters, error) {
	url, err := templateExecutor.Execute(params.URL, tplData)
	if err != nil {
		return params, TemplateErr{field: "URL", err: err}
	}

	payload, err := templateExecutor.Execute(params.Payload, tplData)
	if err != nil {
		return params, TemplateErr{field: "Payload", err: err}
	}

	headers := make(map[string]string, len(params.Headers))
	for k, v := range params.Headers {
		headers[k], err = templateExecutor.Execute(v, tplData)
		if err != nil {
			return params, TemplateErr{field: "Headers." + k, err: err}
		}
	}

	params.URL = url
	params.Payload = payload
	params.Headers = headers

	return params, nil
}
