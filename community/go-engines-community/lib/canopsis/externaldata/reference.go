package externaldata

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

const (
	RefTypeTable = "table"
	RefTypeAPI   = "api"
)

type RefParameters struct {
	Reference string `bson:"reference" json:"reference" binding:"required"`
	Type      string `bson:"type" json:"type" binding:"required"`

	// are used in db external data
	Table        string            `bson:"table,omitempty" json:"table,omitempty"`
	TableName    string            `bson:"table_name,omitempty" json:"-"`
	TableType    *int              `bson:"table_type,omitempty" json:"-"`
	TableColumns []string          `bson:"table_columns,omitempty" json:"-"`
	Select       map[string]string `bson:"select,omitempty" json:"select,omitempty"`
	Regexp       map[string]string `bson:"regexp,omitempty" json:"regexp,omitempty"`
	SortBy       string            `bson:"sort_by,omitempty" json:"sort_by,omitempty"`
	Sort         string            `bson:"sort,omitempty" json:"sort,omitempty" binding:"oneoforempty=asc desc"`
	Optional     bool              `bson:"optional,omitempty" json:"optional,omitempty"`

	// are used in api external data
	Request *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}

type ParsedRefParameters struct {
	Reference string
	Type      string

	TableName    string
	TableType    *int
	TableColumns []string
	Select       map[string]libtemplate.ParsedTemplate
	Regexp       map[string]libtemplate.ParsedTemplate
	SortBy       string
	Sort         string
	Optional     bool

	Request *request.ParsedParameters
}

func ParseRefParameters(data []RefParameters, tplExecutor libtemplate.Executor) []ParsedRefParameters {
	parsed := make([]ParsedRefParameters, len(data))
	for i, params := range data {
		parsedSelect := make(map[string]libtemplate.ParsedTemplate, len(params.Select))
		for k, v := range params.Select {
			parsedSelect[k] = tplExecutor.Parse(v)
		}

		parsedRegexp := make(map[string]libtemplate.ParsedTemplate, len(params.Regexp))
		for k, v := range params.Regexp {
			parsedRegexp[k] = tplExecutor.Parse(v)
		}

		var parsedRequestParameters *request.ParsedParameters
		if params.Request != nil {
			parsedRequestParameters = &request.ParsedParameters{
				URL:        tplExecutor.Parse(params.Request.URL),
				Method:     params.Request.Method,
				Auth:       params.Request.Auth,
				Headers:    params.Request.Headers,
				Payload:    tplExecutor.Parse(params.Request.Payload),
				SkipVerify: params.Request.SkipVerify,
				Timeout:    params.Request.Timeout,
				RetryCount: params.Request.RetryCount,
				RetryDelay: params.Request.RetryDelay,
			}
		}

		parsed[i] = ParsedRefParameters{
			Reference:    params.Reference,
			Type:         params.Type,
			TableName:    params.TableName,
			TableType:    params.TableType,
			TableColumns: params.TableColumns,
			Select:       parsedSelect,
			Regexp:       parsedRegexp,
			SortBy:       params.SortBy,
			Sort:         params.Sort,
			Optional:     params.Optional,
			Request:      parsedRequestParameters,
		}
	}

	return parsed
}
