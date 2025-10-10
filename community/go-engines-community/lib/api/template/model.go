package template

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	TypeTestDataEvent = iota
	TypeTestDataResponse
)

const (
	TypeTestEventFilterRule = iota
	TypeTestLinkRule
	TypeTestActionScenario
	TypeTestWidget
	TypeTestDeclareTicketRule
	TypeTestDynamicInfosRule
	TypeTestInstruction
	TypeTestJob
	TypeTestMetaAlarmRule
)

type ValidateResponse struct {
	IsValid bool                 `json:"is_valid"`
	Err     *validator.ErrReport `json:"err"`
	Result  string               `json:"result"`
}

type EditDataRequest struct {
	ID string `json:"-"`
	// Possible types:
	//   * `0` - Event test data
	//   * `1` - Response test data
	Type        *int              `json:"type" binding:"required,oneof=0 1"`
	Name        string            `json:"name" binding:"required,max=255"`
	Description string            `json:"description" binding:"max=500"`
	Body        map[string]any    `json:"body" binding:"required"`
	Headers     map[string]string `json:"headers" binding:"dive,max=500"`
	Author      string            `json:"author" swaggerignore:"true"`
}

type ListDataRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=_id name description type"`
	// Possible types:
	//   * `0` - Event test data
	//   * `1` - Response test data
	Type         *int   `form:"type"`
	EventPattern string `form:"event_pattern"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type VarResponse struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type DataResponse struct {
	ID          string            `json:"_id" bson:"_id"`
	Type        int               `json:"type" bson:"type"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Body        map[string]any    `json:"body" bson:"body"`
	Headers     map[string]string `json:"headers,omitempty" bson:"headers"`
	Created     datetime.CpsTime  `json:"created" bson:"created" swaggertype:"integer"`
	Updated     datetime.CpsTime  `json:"updated" bson:"updated" swaggertype:"integer"`
}

type AggregationDataResult struct {
	Data       []DataResponse `json:"data" bson:"data"`
	TotalCount int64          `json:"total_count" bson:"total_count"`
}

func (r *AggregationDataResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationDataResult) GetTotal() int64 {
	return r.TotalCount
}

type DataModel struct {
	ID          string            `bson:"_id,omitempty"`
	Type        int               `bson:"type"`
	Name        string            `bson:"name"`
	Description string            `bson:"description"`
	Body        map[string]any    `bson:"body"`
	Headers     map[string]string `bson:"headers"`
	Author      string            `bson:"author"`
	Created     *datetime.CpsTime `bson:"created,omitempty"`
	Updated     *datetime.CpsTime `bson:"updated"`
}

type EditTestRequest struct {
	ID          string `json:"-"`
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=500"`
	// Possible types:
	//   * `0` - Event filter rule
	//   * `1` - Link rule
	//   * `2` - Action scenario
	//   * `3` - Widget
	//   * `4` - Declare ticket rule
	//   * `5` - Dynamic-infos rule
	//   * `6` - Instruction
	//   * `7` - Job
	//   * `8` - Meta-alarm rule
	Type *int   `json:"type" binding:"required"`
	Rule string `json:"rule" binding:"required"`
	Data struct {
		Event     string         `json:"event"`
		Responses map[int]string `json:"responses"`
		Alarm     string         `json:"alarm"`
		Entity    string         `json:"entity"`
		User      string         `json:"user"`
	} `json:"data" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}

type ListTestRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=_id name description type"`
	// Possible types:
	//   * `0` - Event filter rule
	//   * `1` - Link rule
	//   * `2` - Action scenario
	//   * `3` - Widget
	//   * `4` - Declare ticket rule
	//   * `5` - Dynamic-infos rule
	//   * `6` - Instruction
	//   * `7` - Job
	//   * `8` - Meta-alarm rule
	Type *int     `form:"type"`
	Rule string   `form:"rule"`
	IDs  []string `form:"ids[]"`
}

type TestResponse struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Type        int    `json:"type" bson:"type"`
	Rule        struct {
		ID   string `json:"_id" bson:"_id"`
		Name string `json:"name" bson:"name"`
	} `json:"rule" bson:"rule"`
	Data struct {
		Event *struct {
			ID   string `json:"_id" bson:"_id"`
			Name string `json:"name" bson:"name"`
		} `json:"event,omitempty" bson:"event,omitempty"`
		Responses map[int]struct {
			ID   string `json:"_id" bson:"_id"`
			Name string `json:"name" bson:"name"`
		} `json:"responses,omitempty" bson:"responses,omitempty"`
		Alarm  *alarm.RefResponse  `json:"alarm,omitempty" bson:"alarm,omitempty"`
		Entity *entity.RefResponse `json:"entity,omitempty" bson:"entity,omitempty"`
		User   *struct {
			ID          string `json:"_id" bson:"_id"`
			DisplayName string `bson:"display_name" json:"display_name"`
		} `json:"user,omitempty" bson:"user,omitempty"`
	} `json:"data" bson:"data"`
	Created datetime.CpsTime `json:"created" bson:"created" swaggertype:"integer"`
	Updated datetime.CpsTime `json:"updated" bson:"updated" swaggertype:"integer"`
}

type TestModel struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Type        int    `bson:"type"`
	Rule        struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	} `bson:"rule"`
	Data struct {
		Event     string                      `bson:"event,omitempty"`
		Responses map[int]string              `bson:"responses,omitempty"`
		Alarm     *types.AlarmWithEntityField `bson:"alarm,omitempty"`
		Entity    *types.Entity               `bson:"entity,omitempty"`
		User      string                      `bson:"user,omitempty"`
	} `bson:"data"`
	Author  string            `bson:"author"`
	Created *datetime.CpsTime `bson:"created,omitempty"`
	Updated *datetime.CpsTime `bson:"updated"`
}

type AggregationTestResult struct {
	Data       []TestResponse `json:"data" bson:"data"`
	TotalCount int64          `json:"total_count" bson:"total_count"`
}

func (r *AggregationTestResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationTestResult) GetTotal() int64 {
	return r.TotalCount
}

// TemplateRefParameters is a model with all required RefParameters fields for template validation requests.
type TemplateRefParameters struct {
	Reference string              `json:"reference" binding:"required"`
	Type      string              `json:"type" binding:"required"`
	Table     string              `json:"table"`
	Select    map[string]string   `json:"select"`
	Regexp    map[string]string   `json:"regexp"`
	SortBy    string              `json:"sort_by"`
	Sort      string              `json:"sort" binding:"oneoforempty=asc desc"`
	Optional  bool                `json:"optional"`
	Request   *TemplateParameters `json:"request"`
}

// TemplateParameters is a model with all required Parameters fields for template validation requests.
type TemplateParameters struct {
	URL     string            `json:"url" binding:"required"`
	Payload string            `json:"payload"`
	Headers map[string]string `json:"headers"`
}

// TemplateWebhookDeclareTicket is a model with all required WebhookDeclareTicket fields for template validation requests.
type TemplateWebhookDeclareTicket struct {
	TicketIDTpl  string `json:"ticket_id_tpl"`
	TicketURLTpl string `json:"ticket_url_tpl"`
}
