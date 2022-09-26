package eventfilter

import "fmt"

type Container struct {
	applicators map[string]RuleApplicator
}

func (c *Container) Get(ruleType string) (RuleApplicator, bool) {
	return c.applicators[ruleType], c.Has(ruleType)
}

func (c *Container) Set(ruleType string, service RuleApplicator) {
	if c.Has(ruleType) {
		panic(fmt.Errorf("applicator %q already exists", ruleType))
	}

	c.applicators[ruleType] = service
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

type ExternalDataContainer struct {
	applicators map[string]ExternalDataGetter
}

func (c *ExternalDataContainer) Get(dataType string) (ExternalDataGetter, bool) {
	return c.applicators[dataType], c.Has(dataType)
}

func (c *ExternalDataContainer) Set(dataType string, service ExternalDataGetter) {
	if c.Has(dataType) {
		panic(fmt.Errorf("data getter %q already exists", dataType))
	}

	c.applicators[dataType] = service
}

func (c *ExternalDataContainer) Has(dataType string) bool {
	_, ok := c.applicators[dataType]

	return ok
}

func NewExternalDataGetterContainer() *ExternalDataContainer {
	return &ExternalDataContainer{
		applicators: make(map[string]ExternalDataGetter),
	}
}
