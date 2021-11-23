package eventfilter

import (
	"context"
	"log"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Rule is a type representing a rule.
type Rule struct {
	// ID is a unique id for the rule.
	ID string `bson:"_id"`

	// Description contains a short description of the rule.
	Description string `bson:"description"`

	// Type is the rule's type.
	Type Type `bson:"type"`

	// Patterns is the pattern that select the events to which the rule should
	// be applied.
	Patterns pattern.EventPatternList `bson:"patterns"`

	// Priority is an integer used to determine the order in which the rules
	// will be executed.
	Priority int `bson:"priority"`

	// Enabled is a boolean indicating whether the rule is enabled or not.
	Enabled types.OptionalBool `bson:"enabled"`

	// EnrichmentRule contains parameters that are specific to enrichment
	// rules.
	EnrichmentRule `bson:",inline"`

	Created *types.CpsTime `bson:"created"`
	Updated *types.CpsTime `bson:"updated"`
	Author  string         `bson:"author"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// IsEnabled return true if the rule is enabled.
func (r Rule) IsEnabled() bool {
	return !r.Enabled.Set || r.Enabled.Value
}

// Apply applies the rule to an event. Note that Apply does not check that the
// rule is enabled or that its pattern match the event.
func (r Rule) Apply(ctx context.Context, event types.Event, regexMatches pattern.EventRegexMatches, report *Report,
	timezoneConfig *config.TimezoneConfig, logger zerolog.Logger) (types.Event, Outcome) {
	switch r.Type {
	case RuleTypeBreak:
		return event, Break
	case RuleTypeDrop:
		return event, Drop
	case RuleTypeEnrichment:
		return r.EnrichmentRule.Apply(ctx, event, regexMatches, report, timezoneConfig, logger, r.ID)
	default:
		logger.Info().Msgf("Skipping rule with unknown rule type: %s", r.Type)
		return event, Pass
	}
}

// RuleUnpacker is a wrapper arround Rule that implements the bson.Setter
// interface.
type RuleUnpacker struct {
	Rule

	// Valid is a boolean indicating whether the rule is valid or not. It is
	// false when the rule failed to be unmarshalled. It is necessary because
	// returning an error in UnmarshalBSONValue would prevent from querying a MongoDB
	// collection containing invalid rules.
	Valid bool
}

func (r *RuleUnpacker) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	// Mark the rule as invalid (this will be set to true if the unmarshalling
	// succeeds)
	r.Valid = false

	// Parse the rule base (we need to parse this part separately to get the
	// rule's type)
	err := bson.Unmarshal(b, &r.Rule)
	if err != nil {
		log.Printf("unable to parse rule %s: %v", r.ID, err)
		return nil
	}

	if !r.Patterns.IsValid() {
		log.Printf("unable to parse rule %s: invalid patterns", r.ID)
		return nil
	}

	// Check that the unexpected fields of the RuleProcessor are the fields of
	// the RuleBase (that have already been unmarshalled).
	if len(r.UnexpectedFields) > 0 {
		unexpectedFieldNames := make([]string, 0, len(r.UnexpectedFields))
		for key := range r.UnexpectedFields {
			unexpectedFieldNames = append(unexpectedFieldNames, key)
		}

		log.Printf("unexpected fields in rule %s: %s", r.ID, strings.Join(unexpectedFieldNames, ", "))
		return nil
	}

	switch r.Type {
	case RuleTypeDrop, RuleTypeBreak:
		// Check that the enrichment parameters are not defined for drop/break
		if !r.EnrichmentRule.Empty() {
			log.Printf("unable to parse rule %s: the fields external_data, actions, on_success and on_failure should only be defined for enrichment rules.", r.ID)
			return nil
		}
	case RuleTypeEnrichment:
		// Check that the enrichment parameters are valid
		if !r.EnrichmentRule.Valid() {
			log.Printf("unable to parse rule %s: enrichment rules should have at least one action.", r.ID)
			return nil
		}
	default:
		log.Printf("unable to parse rule %s: missing type field", r.ID)
		return nil
	}

	// Mark the rule as valid.
	r.Valid = true
	return nil
}
