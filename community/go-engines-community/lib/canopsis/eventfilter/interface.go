package eventfilter

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/eventfilter/eventfilter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter RuleApplicator,RuleAdapter,RuleApplicatorContainer,Service,ActionProcessor,FailureService,EventCounter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// outcome constant values
const (
	OutcomePass  = "pass"
	OutcomeDrop  = "drop"
	OutcomeBreak = "break"
)

type ActionProcessor interface {
	Process(
		ctx context.Context,
		rule ParsedRule,
		action ParsedAction,
		event *types.Event,
		updatedEntityInfos map[string]UpdatedValue,
		regexMatch RegexMatch,
		externalData map[string]any,
	) (map[string]UpdatedValue, error)
}

type UpdatedValue struct {
	RuleID             string
	OldValue, NewValue any
}

type RuleApplicator interface {
	// Apply eventfilter rule, the first return value(string) should be one of the outcome constant values
	Apply(context.Context, ParsedRule, *types.Event, map[string]UpdatedValue, RegexMatch, map[string]any) (RuleResult, error)
}

type RuleResult struct {
	Outcome            string
	UpdatedEntityInfos map[string]UpdatedValue
}

type RuleAdapter interface {
	GetAll(context.Context) ([]Rule, error)
	GetByTypes(context.Context, []string) ([]Rule, error)
}

type Service interface {
	ProcessEvent(context.Context, *types.Event, ServiceResult) (ServiceResult, error)
	LoadRules(context.Context, []string) error
}

type ServiceResult struct {
	UpdatedEntityInfos      map[string]UpdatedValue `json:"uei,omitzero"`
	ExecutedEnrichRuleCount int64                   `json:"eerc,omitzero"`
	ExternalRequestCount    map[string]int64        `json:"erc,omitzero"`
	ExternalData            map[string]any          `json:"ed,omitzero"`

	// ExternalDataRequest is set when rule processing is suspended to
	// fetch external data asynchronously through a webhook RPC.
	ExternalDataRequest *ExternalDataRequest `json:"edreq,omitzero"`
	// ExternalDataResponse carries the awaited webhook RPC result fed back when a parked event is resumed.
	ExternalDataResponse *ExternalDataResponse `json:"edres,omitzero"`
	// RulesHash fingerprints the ruleset the event was processed against.
	// It is compared on resume to detect a ruleset change (see ErrRulesetChanged).
	RulesHash string `json:"rh,omitzero"`
}

// ExternalDataRequest describes an event filter run that was suspended while external data is fetched asynchronously.
// RuleID is the rule to resume from once WebhookEvent has been answered.
type ExternalDataRequest struct {
	RuleID        string
	CorrelationID string
	WebhookEvent  rpc.WebhookEvent
}

// ExternalDataResponse is the awaited webhook RPC result fed back into a resumed event filter run.
// RuleID identifies the rule that was suspended.
type ExternalDataResponse struct {
	RuleID   string
	Result   []any
	Err      error
	ErrIndex int64
}

// ResumeWithData turns a suspended result into one carrying the fetched external data,
// ready to be fed back into ProcessEvent.
// The request/response rule linkage is restored here so callers cannot mismatch it.
func (r ServiceResult) ResumeWithData(result []any) ServiceResult {
	return r.resume(ExternalDataResponse{Result: result})
}

// ResumeWithError turns a suspended result into one reporting that external data could not be fetched.
// errIndex is the index of the external-data reference that failed.
func (r ServiceResult) ResumeWithError(err error, errIndex int64) ServiceResult {
	return r.resume(ExternalDataResponse{Err: err, ErrIndex: errIndex})
}

func (r ServiceResult) resume(response ExternalDataResponse) ServiceResult {
	response.RuleID = r.ExternalDataRequest.RuleID
	r.ExternalDataRequest = nil
	r.ExternalDataResponse = &response

	return r
}

type RuleApplicatorContainer interface {
	Get(string) (RuleApplicator, bool)
	Set(string, RuleApplicator)
}
