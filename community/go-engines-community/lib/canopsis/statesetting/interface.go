package statesetting

import "context"

//go:generate mockgen -destination=../../../mocks/lib/canopsis/statesetting/statesetting.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting Adapter

type Adapter interface {
	Get(ctx context.Context, settingType string) (StateSetting, error)
}
