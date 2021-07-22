package neweventfilter

type Container struct {
	applicators map[string]RuleApplicator
}

func (c *Container) Get(ruleType string) (RuleApplicator, bool) {
	return c.applicators[ruleType], c.Has(ruleType)
}

func (c *Container) Set(ruleType string, service RuleApplicator) {
	if !c.Has(ruleType) {
		c.applicators[ruleType] = service
	}
}

func (c *Container) Has(ruleType string) bool {
	_, ok := c.applicators[ruleType]

	return ok
}

func NewRuleApplicatorContainer() *Container {
	return &Container{
		applicators: make(map[string]RuleApplicator),
	}
}
