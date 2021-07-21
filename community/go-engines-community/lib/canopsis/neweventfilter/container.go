package neweventfilter

type Container struct {
	applicators map[string]RuleApplicator
}

func (c *Container) Get(ruleType string) (RuleApplicator, bool) {
	if c.Has(ruleType) {
		return c.applicators[ruleType], true
	}

	return nil, false
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
	var container Container

	container.applicators = make(map[string]RuleApplicator)

	return &container
}
