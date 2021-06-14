package eventfilter

import (
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
)

// EnrichmentRule is a type that contains the rule's parameters that are
// specific to enrichment rules.
type EnrichmentRule struct {
	// ExternalData contains external data sources that can be used by the
	// rule's actions
	ExternalData map[string]*DataSource `bson:"external_data" json:"external_data"`

	// Actions is a list of actions to be executed by the rule.
	Actions []Action `bson:"actions" json:"actions"`

	// OnSuccess is the outcome of the rule in case of success of all the
	// actions.
	OnSuccess Outcome `bson:"on_success" json:"on_success"`

	// OnFailure is the outcome of the rule in case of failure of one of the
	// actions.
	OnFailure Outcome `bson:"on_failure" json:"on_failure"`
}

// Empty returns true if the EnrichmentRule have not been changed.
func (p EnrichmentRule) Empty() bool {
	return (len(p.ExternalData) == 0 && len(p.Actions) == 0 &&
		p.OnSuccess == UnsetOutcome && p.OnFailure == UnsetOutcome)
}

// Valid returns true if the EnrichmentRule are valid.
func (p EnrichmentRule) Valid() bool {
	return len(p.Actions) > 0
}

// getExternalData fetches the external data corresponding to an event using
// the external data sources of the rule.
func (p EnrichmentRule) getExternalData(parameters DataSourceGetterParameters) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, source := range p.ExternalData {
		if source.DataSourceGetter == nil {
			return externalData, fmt.Errorf("dataSourceGetter is not initialized")
		}
		data, err := source.Get(parameters)
		if err != nil {
			return externalData, err
		}
		externalData[name] = data
	}

	return externalData, nil
}

// applyActions applies the actions of an enrichment rule to an event.
func (p EnrichmentRule) applyActions(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	var err error

	for _, action := range p.Actions {
		parameters.Event = event
		event, err = action.Apply(event, parameters, report)
		if err != nil {
			return event, err
		}
	}

	return event, nil
}

// Apply applies the enrichment rule to an event.
func (p EnrichmentRule) Apply(event types.Event, regexMatches pattern.EventRegexMatches, report *Report, logger zerolog.Logger, ruleId string) (types.Event, Outcome) {
	parameters := DataSourceGetterParameters{
		Event:      event,
		RegexMatch: regexMatches,
	}
	externalData, err := p.getExternalData(parameters)
	if err != nil {
		logger.Warn().Err(err).Msg("eventfilter | failed to get external data")
		return event, p.OnFailure
	}

	actionParameters := ActionParameters{
		DataSourceGetterParameters: parameters,
		ExternalData:               externalData,
	}
	event, err = p.applyActions(event, actionParameters, report)
	if err != nil {
		marshalledEvent, mErr := json.Marshal(event)
		if mErr != nil {
			logger.Warn().Msg("Failed to marshal event for log")
		}

		logger.Warn().Msgf("Failed to apply %s on %s because of '%s'", ruleId, marshalledEvent, err.Error())

		return event, p.OnFailure
	}

	return event, p.OnSuccess
}
