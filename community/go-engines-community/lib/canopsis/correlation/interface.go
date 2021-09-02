package correlation

//go:generate mockgen -destination=../../../mocks/lib/canopsis/correlation/metaalarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation RulesAdapter

type RulesAdapter interface {
	// Get read meta-alarm rules from db
	Get() ([]Rule, error)

	Save(Rule) error
	GetManualRule() (Rule, error)
	GetRule(id string) (Rule, error)
}
