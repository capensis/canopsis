package metaalarm

// RuleType is a type representing the type of an meta-alarm's rule.
type RuleType string

const (
	// Component-resource matching
	RuleTypeRelation  RuleType = "relation"
	RuleTypeTimeBased RuleType = "timebased"
	// RuleTypeAttribute for Attribute matching
	RuleTypeAttribute RuleType = "attribute"
	// RuleTypeComplex for complex rules
	RuleTypeComplex RuleType = "complex"

	RuleValueGroup RuleType = "valuegroup"

	RuleManualGroup RuleType = "manualgroup"

	RuleCorel RuleType = "corel"
)

type Rule struct {
	// ID is a unique id for the rule.
	ID string `bson:"_id"`

	// RuleType is the rule's type.
	Type RuleType `bson:"type"`

	Patterns RulePatterns `bson:"patterns,omitempty"`

	Config RuleConfig `bson:"config"`

	// Name was added to identify manual grouping
	Name string `bson:"name"`

	AutoResolve bool `bson:"auto_resolve"`

	OutputTemplate string `bson:"output_template"`
}
