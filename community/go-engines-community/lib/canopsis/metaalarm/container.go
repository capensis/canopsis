package metaalarm

type Container struct {
	applicators map[RuleType]RuleApplicator
}

func (container *Container) Get(ruleType RuleType) (RuleApplicator, bool) {
	if container.Has(ruleType) {
		return container.applicators[ruleType], true
	}

	return nil, false
}

func (container *Container) Set(ruleType RuleType, service RuleApplicator) {
	if !container.Has(ruleType) {
		container.applicators[ruleType] = service
	}
}

func (container *Container) Has(ruleType RuleType) bool {
	_, ok := container.applicators[ruleType]

	return ok
}

func NewRuleApplicatorContainer() *Container {
	var container Container

	container.applicators = make(map[RuleType]RuleApplicator)

	return &container
}
