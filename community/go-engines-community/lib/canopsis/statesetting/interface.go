package statesetting

//go:generate mockgen -destination=../../../mocks/lib/canopsis/statesetting/statesetting.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting Adapter

type Adapter interface {
	Get(settingType string) (StateSetting, error)
}
