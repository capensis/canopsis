package types

// CorrelationRuleTags is moved to types because it's used in Event.
type CorrelationRuleTags struct {
	CopyFromChildren bool     `bson:"copy_from_children,omitempty" json:"copy_from_children,omitempty"`
	FilterByLabel    []string `bson:"filter_by_label,omitempty" json:"filter_by_label,omitempty"`
}

// CorrelationRuleInfo is moved to types because it's used in Event.
type CorrelationRuleInfo struct {
	Name             string `bson:"name" json:"name"`
	Description      string `bson:"description,omitempty" json:"description,omitempty"`
	Value            any    `bson:"value,omitempty" json:"value,omitempty"`
	CopyFromChildren bool   `bson:"copy_from_children,omitempty" json:"copy_from_children,omitempty"`
}
