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
	// Table contains id of external table
	Table string `bson:"table,omitempty" json:"table,omitempty"`
	// TableName, TableType, TableColumns contains data of external table to avoid lookup.
	TableName    string   `bson:"table_name,omitempty" json:"-"`
	TableType    *int     `bson:"table_type,omitempty" json:"-"`
	TableColumns []string `bson:"table_columns,omitempty" json:"-"`
	// Select defines equal match conditions.
	// Key contains field name of table, value contains template.
	Select map[string]string `bson:"select,omitempty" json:"select,omitempty"`
	// Regexp defines regexp match conditions.
	// Key contains field name of table with regexp, value contains template.
	Regexp map[string]string `bson:"regexp,omitempty" json:"regexp,omitempty"`
	// SortBy and Sort define priority of found rows. It's id by default.
	SortBy string `bson:"sort_by,omitempty" json:"sort_by,omitempty"`
	Sort   string `bson:"sort,omitempty" json:"sort,omitempty" binding:"oneoforempty=asc desc"`
	// Optional defines will it cause error or not if external data isn't found.
	Optional bool `bson:"optional,omitempty" json:"optional,omitempty"`

	// are used in api external data
	Request *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}

// TemplateRefParameters is a model with all required RefParameters fields for template validation requests.
type TemplateRefParameters struct {
	Reference string                      `json:"reference" binding:"required"`
	Type      string                      `json:"type" binding:"required"`
	Table     string                      `json:"table"`
	Select    map[string]string           `json:"select"`
	Regexp    map[string]string           `json:"regexp"`
	SortBy    string                      `json:"sort_by"`
	Sort      string                      `json:"sort" binding:"oneoforempty=asc desc"`
	Optional  bool                        `json:"optional"`
	Request   *request.TemplateParameters `json:"request"`
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
