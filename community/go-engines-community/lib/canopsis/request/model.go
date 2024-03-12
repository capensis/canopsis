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
